package client

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"sync"
	"time"
)

const (
	KeyPath = "DCC/"
)

var (
	ErrDCCClientInitFailed = errors.New("DCC client init failed")
	ErrLocalCacheNotExists = errors.New("localcache not exists")
	ErrLocalCacheExpire    = errors.New("localcache expire")
)

var _ Client = (*client)(nil)

type Client interface {
	Get(key string) (value []byte, err error)
}

func New(consulKVClient *api.KV, options ...Option) (Client, error) {
	if consulKVClient == nil {
		return nil, ErrDCCClientInitFailed
	}

	c := &client{
		consulKVClient: consulKVClient,
		ttl:            0,
		debug:          false,
		useWatch:       false,
		data:           make(map[string][]byte),
		setTime:        make(map[string]time.Time),
	}

	for _, option := range options {
		c = option(c)
	}

	if c.useWatch {
		if watch, e := NewWatcher(consulKVClient, c.pullInterval); e != nil {
			return nil, e
		} else {
			c.watch = watch
		}

		go func() {
			c.watchAndUpdate()
		}()
	}

	return c, nil
}

type client struct {
	consulKVClient *api.KV

	ttl   time.Duration
	debug bool

	useWatch     bool
	pullInterval time.Duration

	// 本地缓存配置
	data    map[string][]byte    // 配置数据的本地缓存
	setTime map[string]time.Time // 更新的时间
	lock    sync.RWMutex

	// 本地watch consul的key变更
	watch Watcher
}

func (c *client) Get(key string) (value []byte, err error) {
	key = KeyPath + key

	// 增加监听的 key
	if c.useWatch {
		c.watch.AddKey(key)
	}

	// 拉取流程
	if val, e := c.getFromLocalCache(key); e != nil {
		if val, e := c.getFromConsul(key); e != nil {
			return nil, e
		} else {
			c.setToLocalCache(key, val)
			return val, nil
		}
	} else {
		return val, nil
	}
}

func (c *client) getFromLocalCache(key string) (value []byte, err error) {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if c.debug {
		fmt.Printf("getFromLocalCache  key=%s\n", key)
	}

	// 未失效直接获取
	if value, ok := c.data[key]; ok {
		// 检查key是否失效
		if t, _ := c.setTime[key]; t.Add(c.ttl).After(time.Now()) {
			// 未失效
			return value, nil
		} else {
			// 已失效 删除map
			delete(c.setTime, key)
			delete(c.data, key)
			return nil, ErrLocalCacheExpire
		}
	} else {
		return nil, ErrLocalCacheNotExists
	}
}

func (c *client) setToLocalCache(key string, value []byte) {
	defer c.lock.Unlock()
	c.lock.Lock()
	c.data[key] = value
	c.setTime[key] = time.Now()
}

func (c *client) getFromConsul(key string) (value []byte, err error) {
	fmt.Printf("getFromConsul  key=%s\n", key)

	var kvPair *api.KVPair
	if kvPair, _, err = c.consulKVClient.Get(key, nil); err != nil {
		return nil, err
	}
	if kvPair == nil {
		return nil, nil
	}
	return kvPair.Value, nil
}

// watchAndUpdate 消费 watch 的 chan 更新本地cache
func (c *client) watchAndUpdate() {
	for pair := range c.watch.WatchAllKey() {
		// if c.debug {
		// 	fmt.Printf("watchAndUpdate  轮训的所有KV  pair=%+v\n", pair)
		// }

		if val, e := c.getFromLocalCache(pair.Key); e != nil {
			// FIXME: 优化err case
			continue
		} else if EqualSliceByte(val, pair.Value) {
			// 未更新
			continue
		} else {
			// 更新本地缓存
			c.setToLocalCache(pair.Key, pair.Value)

			if c.debug {
				fmt.Printf("watchAndUpdate  更新KV  pair=%+v\n", pair)
			}
		}
	}
}
