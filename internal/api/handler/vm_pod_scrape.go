package handler

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/xiaozongyang/vm-access/internal/api/registry"
	"github.com/xiaozongyang/vm-access/internal/types"
	"github.com/xiaozongyang/vm-access/pkg/utils"
)

type VMPodScrapeTemplateArgs struct {
	Create      bool
	Readonly    bool
	Cluster     string
	VMPodScrape *types.VMPodScrape
}

func NewVMPodScrape(ctx context.Context, c *app.RequestContext) {
	var req types.VMPodScrape
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

	if err := proxyClient.CreateVMPodScrape(ctx, &req); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if shouldReturnHTML(c) {
		c.HTML(consts.StatusOK, "vm-pod-scrape.html", VMPodScrapeTemplateArgs{Create: true, VMPodScrape: &req, Cluster: cluster, Readonly: false})
	} else {
		c.JSON(consts.StatusOK, req)
	}
}

func HTMLNewVMPodScrape(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")
	c.HTML(consts.StatusOK, "vm-pod-scrape.html", utils.H{"Cluster": cluster, "Create": true})
}

func GetVMPodScrape(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	scrape, err := proxyClient.GetVMPodScrape(ctx, c.Param("vm-pod-scrape"))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if shouldReturnHTML(c) {
		c.HTML(consts.StatusOK, "vm-pod-scrape.html", VMPodScrapeTemplateArgs{Create: false, VMPodScrape: scrape, Cluster: cluster, Readonly: isReservedVMPodScrape(scrape.Name)})
	} else {
		c.JSON(consts.StatusOK, scrape)
	}
}

func ListVMPodScrapes(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	scrapes, err := proxyClient.ListVMPodScrapes(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, scrapes)
}

func UpdateVMPodScrape(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	scrapeName := c.Param("vm-pod-scrape")

	var req types.VMPodScrape
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if scrapeName != req.Name {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": "name cannot be changed"})
		return
	}

	if err := proxyClient.UpdateVMPodScrape(ctx, scrapeName, &req); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if shouldReturnHTML(c) {
		c.HTML(consts.StatusOK, "vm-pod-scrape.html", VMPodScrapeTemplateArgs{Create: false, VMPodScrape: &req, Cluster: cluster, Readonly: isReservedVMPodScrape(scrapeName)})
	} else {
		c.JSON(consts.StatusOK, req)
	}
}

func DeleteVMPodScrape(ctx context.Context, c *app.RequestContext) {
	cluster := c.Param("cluster")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	scrapeName := c.Param("vm-pod-scrape")

	if isReservedVMPodScrape(scrapeName) {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": fmt.Sprintf("VMPodScrape %s is reserved, cannot be deleted", c.Param("vm-pod-scrape"))})
		return
	}

	if err := proxyClient.DeleteVMPodScrape(ctx, scrapeName); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, map[string]string{"message": "ok"})
}
