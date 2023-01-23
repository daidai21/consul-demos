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
		data:           make(map[string][]byte),
		setTime:        make(map[string]time.Time),
	}

	for _, option := range options {
		c = option(c)
	}

	return c, nil
}

type client struct {
	consulKVClient *api.KV

	ttl   time.Duration
	debug bool

	// 本地缓存配置
	data    map[string][]byte
	setTime map[string]time.Time
	lock    sync.RWMutex

	// TODO:
	// 本地watch consul的key变更
	// 新开一个单独的watcher文件
	// https://github.com/pteich/consul-kv-watcher/blob/main/watcher.go
}

func (c *client) Get(key string) (value []byte, err error) {
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
	if kvPair, _, err = c.consulKVClient.Get(KeyPath+key, nil); err != nil {
		return nil, err
	}
	if kvPair == nil {
		return nil, nil
	}
	return kvPair.Value, nil
}
