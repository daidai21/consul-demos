package client

import "time"

type Option func(c *client) *client

// WithLocalCacheTTLSeconds 本地缓存时间
func WithLocalCacheTTLSeconds(seconds int) Option {
	return func(c *client) *client {
		if seconds <= 0 {
			seconds = 0
		}
		c.ttl = time.Duration(seconds) * time.Second
		return c
	}
}

// WithDebug debug模式
func WithDebug() Option {
	return func(c *client) *client {
		c.debug = true
		return c
	}
}
