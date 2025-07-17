package client

import (
	"context"
	"time"

	hzclient "github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
)

func RequestLogger() hzclient.Middleware {
	return func(next hzclient.Endpoint) hzclient.Endpoint {
		return func(ctx context.Context, req *protocol.Request, resp *protocol.Response) error {
			beforeTime := time.Now()

			err := next(ctx, req, resp)
			if err != nil {
				return err
			}

			afterRequest(beforeTime, req, resp)

			return nil
		}
	}
}

func afterRequest(beforeTime time.Time, req *protocol.Request, resp *protocol.Response) {
	duration := time.Since(beforeTime)

	var reqBody string
	var respBody string
	if len(req.Body()) > 0 {
		reqBody = string(req.Body())
	}
	if len(resp.Body()) > 0 {
		respBody = string(resp.Body())
	}

	hlog.Infof("request logger: %s %s %s status: %d, duration: %s request body: %s %s response body: %s", req.Method(), req.Host(), req.RequestURI(), resp.StatusCode(), duration, reqBody, respBody)
}
