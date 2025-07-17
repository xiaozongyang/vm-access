package route

import (
	"context"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/xiaozongyang/vm-access/internal/api/handler"
	"github.com/xiaozongyang/vm-access/pkg/utils"
)

func Register(h *server.Hertz) {
	// load templates should be before static and api
	h.LoadHTMLGlob(getTemplateRoot())
	h.Static("/static", getStaticRoot())

	h.GET("/", handler.HTMLListClusters)
	clusters := h.Group("/clusters")
	clusters.GET("/", handler.HTMLListClusters)
	clusters.GET("/:cluster", handler.GetScrapesByCluster)

	clusters.GET("/:cluster/vm-static-scrapes/new", handler.HTMLNewVMStaticScrape)
	clusters.GET("/:cluster/vm-static-scrapes/:vm-static-scrape", handler.GetVMStaticScrape)

	clusters.GET("/:cluster/vm-service-scrapes/new", handler.HTMLNewVMServiceScrape)
	clusters.GET("/:cluster/vm-service-scrapes/:vm-service-scrape", handler.GetVMServiceScrape)

	clusters.GET("/:cluster/vm-pod-scrapes/new", handler.HTMLNewVMPodScrape)
	clusters.GET("/:cluster/vm-pod-scrapes/:vm-pod-scrape", handler.GetVMPodScrape)

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	apiV1 := h.Group("/api/v1")
	apiV1.GET("/clusters", handler.ListClusters)
	apiV1.POST("/clusters/ping", handler.Ping)
	apiV1.GET("/proxies", handler.DumpProxies)

	cluster := apiV1.Group("/:cluster")

	vms := cluster.Group("/vm-static-scrapes")
	vms.POST("/", handler.CreateVMStaticScrape)
	vms.GET("/", handler.ListVMStaticScrapes)
	vms.GET("/:vm-static-scrape", handler.GetVMStaticScrape)
	vms.PUT("/:vm-static-scrape", handler.UpdateVMStaticScrape)
	vms.DELETE("/:vm-static-scrape", handler.DeleteVMStaticScrape)

	vmss := cluster.Group("/vm-service-scrapes")
	vmss.POST("/", handler.NewVMServiceScrape)
	vmss.GET("/", handler.ListVMServiceScrapes)
	vmss.GET("/:vm-service-scrape", handler.GetVMServiceScrape)
	vmss.PUT("/:vm-service-scrape", handler.UpdateVMServiceScrape)
	vmss.DELETE("/:vm-service-scrape", handler.DeleteVMServiceScrape)

	vmpods := cluster.Group("/vm-pod-scrapes")
	vmpods.POST("/", handler.NewVMPodScrape)
	vmpods.GET("/", handler.ListVMPodScrapes)
	vmpods.GET("/:vm-pod-scrape", handler.GetVMPodScrape)
	vmpods.PUT("/:vm-pod-scrape", handler.UpdateVMPodScrape)
	vmpods.DELETE("/:vm-pod-scrape", handler.DeleteVMPodScrape)
}

func getStaticRoot() string {
	if os.Getenv("VM_ACCESS_STATIC_ROOT") != "" {
		return os.Getenv("VM_ACCESS_STATIC_ROOT")
	}
	return "/"
}

func getTemplateRoot() string {
	dir := os.Getenv("VM_ACCESS_TEMPLATE_ROOT")
	return dir + "/*.html"
}
