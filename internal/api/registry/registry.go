package registry

import (
	"fmt"
	"sort"
	"sync"

	"github.com/xiaozongyang/vm-access/internal/proxy/client"
	"github.com/xiaozongyang/vm-access/internal/types"
)

var (
	cluster2proxyAddr = make(map[string]string)
	c2proxyAddrMutex  = sync.RWMutex{}

	cluster2proxyClient = make(map[string]*client.Client)
	c2proxyClientMutex  = sync.RWMutex{}
)

func ListClusters() []string {
	c2proxyAddrMutex.RLock()
	defer c2proxyAddrMutex.RUnlock()

	clusters := make([]string, 0, len(cluster2proxyAddr))
	for cluster := range cluster2proxyAddr {
		clusters = append(clusters, cluster)
	}

	sort.Strings(clusters)

	return clusters
}

func Register(cluster, addr string) {
	c2proxyAddrMutex.Lock()
	defer c2proxyAddrMutex.Unlock()

	cluster2proxyAddr[cluster] = addr
}

func Dump() []types.PingReq {
	c2proxyAddrMutex.RLock()
	defer c2proxyAddrMutex.RUnlock()

	proxies := make([]types.PingReq, 0, len(cluster2proxyAddr))
	for cluster, addr := range cluster2proxyAddr {
		proxies = append(proxies, types.PingReq{Cluster: cluster, Addr: addr})
	}

	return proxies
}

func GetOrCreateProxyClientLocked(cluster string) (*client.Client, error) {
	proxyAddr, found := getProxyAddr(cluster)
	if !found {
		return nil, fmt.Errorf("proxy addr for cluster %s not found", cluster)
	}

	c2proxyClientMutex.RLock()
	proxyClient, found := cluster2proxyClient[cluster]
	c2proxyClientMutex.RUnlock()

	if found && proxyClient.Addr() == proxyAddr {
		return proxyClient, nil
	}

	proxyClient, err := client.New(proxyAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy client, cluster: %s, addr: %s, err: %w", cluster, proxyAddr, err)
	}

	c2proxyClientMutex.Lock()
	cluster2proxyClient[cluster] = proxyClient
	c2proxyClientMutex.Unlock()

	return proxyClient, nil
}

func getProxyAddr(cluster string) (string, bool) {
	c2proxyAddrMutex.RLock()
	defer c2proxyAddrMutex.RUnlock()

	proxyClient, ok := cluster2proxyAddr[cluster]
	return proxyClient, ok
}
