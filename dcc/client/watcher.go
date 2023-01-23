package client

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"sync"
	"time"
)

var (
	ErrWatcherInitFailed = errors.New("watcher init failed")
)

var _ Watcher = (*watcher)(nil)

type Pair struct {
	Key   string
	Value []byte
}

// Watcher 监听 KV 的变更
// TCC 也是定时pull的方式
// 看了下还是 goroutine for 轮训不断拉数据
type Watcher interface {
	AddKey(key string)
	WatchAllKey() <-chan *Pair
}

func NewWatcher(consulKVClient *api.KV, interval time.Duration) (Watcher, error) {
	if consulKVClient == nil {
		return nil, ErrWatcherInitFailed
	}
	if interval < 0 {
		return nil, ErrWatcherInitFailed
	}

	w := &watcher{
		consulKVClient: consulKVClient,
		interval:       interval,
		out:            make(chan *Pair, 100),
		lock:           sync.RWMutex{},
	}

	go func() {
		for {
			w.lock.RLock()
			for _, key := range w.keys {
				if kvPair, _, e := w.consulKVClient.Get(key, nil); e != nil {
					// FIXME: 可以增加重试和告警机制
					continue
				} else {
					if kvPair != nil {
						w.out <- &Pair{
							Key:   kvPair.Key,
							Value: kvPair.Value,
						}
					}
				}
			}
			w.lock.RUnlock()
			time.Sleep(w.interval)
		}
	}()

	return w, nil
}

type watcher struct {
	consulKVClient *api.KV

	keys []string

	interval time.Duration

	out chan *Pair

	lock sync.RWMutex
}

func (w *watcher) AddKey(key string) {
	defer w.lock.Unlock()
	w.lock.Lock()
	w.keys = append(w.keys, key)
}

// WatchAllKey 只能获取一次，否则类似mq的分区消费了
func (w *watcher) WatchAllKey() <-chan *Pair {
	return w.out
}
