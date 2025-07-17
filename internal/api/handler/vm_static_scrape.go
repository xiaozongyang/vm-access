package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/xiaozongyang/vm-access/internal/api/registry"
	"github.com/xiaozongyang/vm-access/internal/types"
	"github.com/xiaozongyang/vm-access/pkg/utils"
)

func CreateVMStaticScrape(ctx context.Context, c *app.RequestContext) {
	var req types.VMStaticScrape
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	proxyClient, err := registry.GetOrCreateProxyClientLocked(req.Cluster)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = proxyClient.CreateVMStaticScrape(ctx, &req)
	if err != nil {
		hlog.Errorf("failed to create vm static scrape, req:%+v, err: %v", req, err)
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, "ok")
}

func ListVMStaticScrapes(ctx context.Context, c *app.RequestContext) {
	clusterName := c.Param("cluster")
	if clusterName == "" {
		err := fmt.Errorf("cluster name is required")
		c.AbortWithError(consts.StatusBadRequest, err)
		return
	}

	vss, err := listVMStaticScrapes(ctx, c, clusterName)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if shouldReturnHTML(c) {
		c.HTML(consts.StatusOK, "list-vm-static-scrapes.html", utils.H{"Cluster": clusterName, "VMStaticScrapes": vss})
	} else {
		c.JSON(consts.StatusOK, vss)
	}
}

func listVMStaticScrapes(ctx context.Context, c *app.RequestContext, cluster string) ([]types.VMStaticScrape, error) {

	proxyClient, err := registry.GetOrCreateProxyClientLocked(cluster)
	if err != nil {
		c.AbortWithError(consts.StatusInternalServerError, err)
		return nil, err
	}

	list, err := proxyClient.ListVMStaticScrapes(ctx)
	if err != nil {
		c.AbortWithError(consts.StatusInternalServerError, err)
		return nil, err
	}
	return list.Items, nil
}

func HTMLNewVMStaticScrape(ctx context.Context, c *app.RequestContext) {
	clusterName := c.Param("cluster")
	c.HTML(consts.StatusOK, "vm-static-scrape.html", utils.H{"Cluster": clusterName, "Create": true})
}

func GetVMStaticScrape(ctx context.Context, c *app.RequestContext) {
	clusterName := c.Param("cluster")
	vmStaticScrapeName := c.Param("vm-static-scrape")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(clusterName)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	vss, err := proxyClient.GetVMStaticScrape(ctx, vmStaticScrapeName)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if shouldReturnHTML(c) {
		c.HTML(consts.StatusOK, "vm-static-scrape.html", utils.H{"Cluster": clusterName, "VMStaticScrape": vss, "Create": false})
	} else {
		c.JSON(consts.StatusOK, vss)
	}
}

func UpdateVMStaticScrape(ctx context.Context, c *app.RequestContext) {
	clusterName := c.Param("cluster")
	vmStaticScrapeName := c.Param("vm-static-scrape")

	var req types.VMStaticScrape
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	proxyClient, err := registry.GetOrCreateProxyClientLocked(clusterName)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = proxyClient.UpdateVMStaticScrape(ctx, vmStaticScrapeName, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, "ok")
}

func DeleteVMStaticScrape(ctx context.Context, c *app.RequestContext) {
	clusterName := c.Param("cluster")
	vmStaticScrapeName := c.Param("vm-static-scrape")

	proxyClient, err := registry.GetOrCreateProxyClientLocked(clusterName)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = proxyClient.DeleteVMStaticScrape(ctx, vmStaticScrapeName)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, "ok")
}

func shouldReturnHTML(c *app.RequestContext) bool {
	return strings.Contains(string(c.GetHeader("Accept")), "text/html")
}
