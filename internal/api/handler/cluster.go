package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/xiaozongyang/vm-access/internal/api/registry"
	"github.com/xiaozongyang/vm-access/internal/types"
	"github.com/xiaozongyang/vm-access/pkg/utils"
)

func HTMLListClusters(ctx context.Context, c *app.RequestContext) {
	clusters := registry.ListClusters()
	c.HTML(consts.StatusOK, "list-clusters.html", utils.H{"Clusters": clusters})
}

func ListClusters(ctx context.Context, c *app.RequestContext) {
	clusters := registry.ListClusters()

	c.JSON(consts.StatusOK, clusters)
}

func Ping(ctx context.Context, c *app.RequestContext) {
	var req types.PingReq
	if err := c.BindJSON(&req); err != nil {
		c.String(consts.StatusBadRequest, "invalid request")
		return
	}
	if req.Cluster == "" || req.Addr == "" {
		c.String(consts.StatusBadRequest, "invalid request, required non-empty cluster and addr")
		return
	}

	registry.Register(req.Cluster, req.Addr)

	c.Status(consts.StatusOK)
}

func DumpProxies(ctx context.Context, c *app.RequestContext) {
	proxies := registry.Dump()
	c.JSON(consts.StatusOK, proxies)
}

func GetScrapesByCluster(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")
	if cluster == "" {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": "cluster is required"})
		return
	}
	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	staticScrapes, err := proxyClient.ListVMStaticScrapes(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	serviceScrapes, err := proxyClient.ListVMServiceScrapes(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	podScrapes, err := proxyClient.ListVMPodScrapes(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if shouldReturnHTML(c) {
		c.HTML(consts.StatusOK, "cluster.html", utils.H{"Cluster": cluster, "VMStaticScrapes": staticScrapes.Items, "VMServiceScrapes": serviceScrapes, "VMPodScrapes": podScrapes})
	} else {
		c.JSON(consts.StatusOK, utils.H{"Cluster": cluster, "VMStaticScrapes": staticScrapes.Items, "VMServiceScrapes": serviceScrapes, "VMPodScrapes": podScrapes})
	}
}
