package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/xiaozongyang/vm-access/internal/middlewares"
	"github.com/xiaozongyang/vm-access/internal/proxy"
	"github.com/xiaozongyang/vm-access/internal/proxy/env"
	"github.com/xiaozongyang/vm-access/internal/proxy/route"
	"github.com/xiaozongyang/vm-access/internal/proxy/vm_operator_client"
)

func main() {
	env.MustInit()

	registryAddr := os.Getenv("REGISTRY_ADDR")
	if registryAddr == "" {
		panic("env REGISTRY_ADDR is not set")
	}

	proxyIP := os.Getenv("PROXY_IP")
	if proxyIP == "" {
		panic("env PROXY_IP is not set")
	}

	proxyPort := os.Getenv("PROXY_PORT")
	if proxyPort == "" {
		panic("env PROXY_PORT is not set")
	}
	proxyAddr := proxyIP + ":" + proxyPort

	vm_operator_client.MustInitVMOperatorClient()

	stop := make(chan struct{})
	go proxy.Ping(registryAddr, env.GetCluster(), proxyAddr, stop)

	h := server.Default()

	route.Register(h)
	h.Use(middlewares.RequestObservation())

	go func() {
		h.Spin()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	close(stop)
	h.Shutdown(context.Background())
}
