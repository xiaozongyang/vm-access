package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/xiaozongyang/vm-access/internal/api/route"
	"github.com/xiaozongyang/vm-access/internal/middlewares"
	"github.com/xiaozongyang/vm-access/pkg/metrics"
)

func main() {
	h := server.New(
		server.WithHostPorts("0.0.0.0:8080"),
	)

	route.Register(h)
	h.Use(middlewares.RequestObservation())

	metrics.StartMetricServer()

	h.Spin()
}
