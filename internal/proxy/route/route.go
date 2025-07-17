package route

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/xiaozongyang/vm-access/internal/proxy/handler"
	"github.com/xiaozongyang/vm-access/pkg/metrics"
	"github.com/xiaozongyang/vm-access/pkg/utils"
)

func Register(h *server.Hertz) {
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})
	metrics.RegisterToHertz(h)

	apiV1 := h.Group("/api/v1")
	vms := apiV1.Group("/vm-static-scrapes")
	vms.POST("", handler.CreateVMStaticScrape)
	vms.GET("", handler.ListVMStaticScrapes)
	vms.GET("/:vm-static-scrape", handler.GetVMStaticScrape)
	vms.PUT("/:vm-static-scrape", handler.UpdateVMStaticScrape)
	vms.DELETE("/:vm-static-scrape", handler.DeleteVMStaticScrape)

	serviceScrapes := apiV1.Group("/vm-service-scrapes")
	serviceScrapes.POST("", handler.CreateVMServiceScrape)
	serviceScrapes.GET("", handler.ListVMServiceScrapes)
	serviceScrapes.GET("/:vm-service-scrape", handler.GetVMServiceScrape)
	serviceScrapes.PUT("/:vm-service-scrape", handler.UpdateVMServiceScrape)
	serviceScrapes.DELETE("/:vm-service-scrape", handler.DeleteVMServiceScrape)

	podScrapes := apiV1.Group("/vm-pod-scrapes")
	podScrapes.POST("", handler.CreateVMPodScrape)
	podScrapes.GET("", handler.ListVMPodScrapes)
	podScrapes.GET("/:vm-pod-scrape", handler.GetVMPodScrape)
	podScrapes.PUT("/:vm-pod-scrape", handler.UpdateVMPodScrape)
	podScrapes.DELETE("/:vm-pod-scrape", handler.DeleteVMPodScrape)
}
