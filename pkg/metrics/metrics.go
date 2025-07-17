package metrics

import (
	"context"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func StartMetricServer() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WriteProcessMetrics(w)
		metrics.WritePrometheus(w, true)
	})

	go func() {
		hlog.Infof("metrics server started on port 19191")
		if err := http.ListenAndServe(":19191", nil); err != nil {
			hlog.Errorf("failed to start metrics server: %v", err)
		}
	}()
}

func RegisterToHertz(h *server.Hertz) {
	h.GET("/metrics", func(c context.Context, ctx *app.RequestContext) {
		metrics.WriteProcessMetrics(ctx.Response.BodyWriter())
		metrics.WritePrometheus(ctx.Response.BodyWriter(), true)
	})
}
