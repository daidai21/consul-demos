package client

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"testing"
	"time"
)

func Test_Get(t *testing.T) {
	consulClient, _ := api.NewClient(api.DefaultConfig())
	DCCClient, _ := New(consulClient.KV(), WithDebug(), WithLocalCacheTTLSeconds(3))

	fmt.Println(DCCClient.Get("test")) // from consul
	fmt.Println(DCCClient.Get("test")) // from localcache
	time.Sleep(4 * time.Second)
	fmt.Println(DCCClient.Get("test")) // from consul, localcache expired
}

func Test_Get_UpdateDelay(t *testing.T) {
	consulClient, _ := api.NewClient(api.DefaultConfig())
	DCCClient, _ := New(consulClient.KV(),
		WithDebug(), WithLocalCacheTTLSeconds(3),
		WithWatch(), WithWatchPullInterval(1),
	)

	loop := 10
	for {
		if loop <= 0 {
			break
		}
		val, _ := DCCClient.Get("test")
		fmt.Println(string(val))
		time.Sleep(time.Second)
		loop--
	}
}
