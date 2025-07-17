package proxy

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/xiaozongyang/vm-access/internal/types"
)

var (
	pingClient *client.Client
)

func init() {
	c, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	pingClient = c
}

func Ping(registryAddr, cluster, proxyAddr string, stop chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	pingReq := &types.PingReq{
		Cluster: cluster,
		Addr:    proxyAddr,
	}

	rawReq, err := json.Marshal(pingReq)
	if err != nil {
		hlog.Errorf("marshal ping req failed: %v", err)
		return
	}

	ping(pingClient, registryAddr, rawReq)

	for {
		select {
		case <-ticker.C:
			ping(pingClient, registryAddr, rawReq)
		case <-stop:
			return
		}
	}
}

func ping(c *client.Client, registryAddr string, rawReq []byte) {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)
	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	req.SetMethod("POST")
	req.SetRequestURI(registryAddr + "/api/v1/clusters/ping")
	req.SetBody(rawReq)

	if err := c.Do(context.Background(), req, resp); err != nil {
		hlog.Errorf("ping proxy failed: %v", err)
		return
	}

	if resp.StatusCode() != 200 {
		hlog.Errorf("ping proxy failed: status code %d, body %s", resp.StatusCode(), string(resp.Body()))
	}
}
