package middlewares

import (
	"context"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	requestTotal    = metrics.NewCounter("hertz_request_total")
	requestErrors   = metrics.NewCounter("hertz_request_errors_total")
	requestDuration = metrics.NewSummary("hertz_request_duration_seconds")
)

func RequestObservation() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		beforeTime := time.Now()
		defer afterRequest(beforeTime, c)

		c.Next(ctx)
	}
}

func afterRequest(beforeTime time.Time, c *app.RequestContext) {
	duration := time.Since(beforeTime)

	requestTotal.Inc()
	metrics.GetOrCreateSummary(getDurationMetricName(c)).Update(float64(duration.Milliseconds()))
	if c.Response.StatusCode() >= 400 {
		requestErrors.Inc()
	}

	var reqBody string
	var respBody string
	if len(c.Request.Body()) > 0 {
		reqBody = string(c.Request.Body())
	}
	if len(c.Response.Body()) > 0 {
		respBody = string(c.Response.Body())
	}
	hlog.Infof("request logger: %s %s %s status: %d, duration: %s request body: %s %s response body: %s", c.Request.Method(), c.Request.Host(), c.Request.RequestURI(), c.Response.StatusCode(), duration, reqBody, respBody)
}

func getDurationMetricName(c *app.RequestContext) string {
	return fmt.Sprintf(`hertz_request_duration_ms{method="%s",path="%s"}`, c.Request.Method(), c.Request.RequestURI())
}
