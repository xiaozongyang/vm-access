package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/xiaozongyang/vm-access/internal/api/registry"
	"github.com/xiaozongyang/vm-access/internal/types"
	"github.com/xiaozongyang/vm-access/pkg/utils"
)

func NewVMServiceScrape(ctx context.Context, c *app.RequestContext) {
	var req types.VMServiceScrape
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	cluster := c.Param("cluster")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := proxyClient.CreateVMServiceScrape(ctx, &req); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if shouldReturnHTML(c) {
		c.HTML(consts.StatusOK, "vm-service-scrape.html", utils.H{"Create": true, "VMServiceScrape": req, "Cluster": cluster})
	} else {
		c.JSON(consts.StatusOK, req)
	}
}

func HTMLNewVMServiceScrape(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")
	c.HTML(consts.StatusOK, "vm-service-scrape.html", utils.H{"Cluster": cluster, "Create": true})
}

func GetVMServiceScrape(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	scrape, err := proxyClient.GetVMServiceScrape(ctx, c.Param("vm-service-scrape"))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if shouldReturnHTML(c) {
		c.HTML(consts.StatusOK, "vm-service-scrape.html", utils.H{"Create": false, "VMServiceScrape": scrape, "Cluster": cluster})
	} else {
		c.JSON(consts.StatusOK, scrape)
	}
}

func ListVMServiceScrapes(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	scrapes, err := proxyClient.ListVMServiceScrapes(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, scrapes)
}

func UpdateVMServiceScrape(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var req types.VMServiceScrape
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := proxyClient.UpdateVMServiceScrape(ctx, c.Param("vm-service-scrape"), &req); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if shouldReturnHTML(c) {
		c.HTML(consts.StatusOK, "vm-service-scrape.html", utils.H{"Create": false, "VMServiceScrape": req, "Cluster": cluster})
	} else {
		c.JSON(consts.StatusOK, req)
	}
}

func DeleteVMServiceScrape(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := proxyClient.DeleteVMServiceScrape(ctx, c.Param("vm-service-scrape")); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, map[string]string{"message": "ok"})
}
