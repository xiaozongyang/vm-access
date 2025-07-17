package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	hzclient "github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/xiaozongyang/vm-access/internal/types"
)

const (
	apiPrefix = "/api/v1"
)

type Client struct {
	addr string

	hertzClient *hzclient.Client
}

func New(addr string) (*Client, error) {
	hertzClient, err := hzclient.NewClient(
		hzclient.WithDialTimeout(time.Second * 10),
	)
	hertzClient.Use(RequestLogger())

	if err != nil {
		return nil, err
	}
	return &Client{
		addr:        addr,
		hertzClient: hertzClient,
	}, nil
}

func (c *Client) Addr() string {
	return c.addr
}

func (c *Client) ListVMStaticScrapes(ctx context.Context) (*types.VMStaticScrapeList, error) {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodGet)
	req.SetRequestURI(c.getUrl("/vm-static-scrapes"))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err := c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	var list types.VMStaticScrapeList
	err = json.Unmarshal(resp.Body(), &list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (c *Client) CreateVMStaticScrape(ctx context.Context, vss *types.VMStaticScrape) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodPost)
	req.SetRequestURI(c.getUrl("/vm-static-scrapes"))
	req.SetHeader("Content-Type", "application/json")

	reqBody, err := json.Marshal(vss)
	if err != nil {
		return err
	}
	req.SetBody(reqBody)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err = c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() >= consts.StatusBadRequest {
		return errors.New("failed to create vm static scrape")
	}

	return nil
}

func (c *Client) GetVMStaticScrape(ctx context.Context, name string) (*types.VMStaticScrape, error) {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodGet)
	req.SetRequestURI(c.getUrl(fmt.Sprintf("/vm-static-scrapes/%s", name)))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err := c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != consts.StatusOK {
		return nil, fmt.Errorf("failed to get vm static scrape: %s", string(resp.Body()))
	}

	var scrape types.VMStaticScrape
	err = json.Unmarshal(resp.Body(), &scrape)
	if err != nil {
		return nil, err
	}

	return &scrape, nil
}

func (c *Client) UpdateVMStaticScrape(ctx context.Context, name string, scrape *types.VMStaticScrape) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodPut)
	req.SetRequestURI(c.getUrl(fmt.Sprintf("/vm-static-scrapes/%s", name)))

	reqBody, err := json.Marshal(scrape)
	if err != nil {
		return err
	}
	req.SetBody(reqBody)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err = c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != consts.StatusOK {
		return fmt.Errorf("failed to update vm static scrape: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) DeleteVMStaticScrape(ctx context.Context, name string) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodDelete)
	req.SetRequestURI(c.getUrl(fmt.Sprintf("/vm-static-scrapes/%s", name)))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err := c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != consts.StatusOK {
		return fmt.Errorf("failed to delete vm static scrape: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) ListVMServiceScrapes(ctx context.Context) ([]types.VMServiceScrape, error) {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodGet)
	req.SetRequestURI(c.getUrl("/vm-service-scrapes"))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err := c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != consts.StatusOK {
		return nil, fmt.Errorf("failed to list vm service scrapes: %s", string(resp.Body()))
	}

	var scrapes []types.VMServiceScrape
	if err := json.Unmarshal(resp.Body(), &scrapes); err != nil {
		return nil, err
	}

	return scrapes, nil
}

func (c *Client) GetVMServiceScrape(ctx context.Context, name string) (*types.VMServiceScrape, error) {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodGet)
	req.SetRequestURI(c.getUrl(fmt.Sprintf("/vm-service-scrapes/%s", name)))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err := c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != consts.StatusOK {
		return nil, fmt.Errorf("failed to get vm service scrape: %s", string(resp.Body()))
	}

	var scrape types.VMServiceScrape
	if err := json.Unmarshal(resp.Body(), &scrape); err != nil {
		return nil, err
	}

	return &scrape, nil
}

func (c *Client) CreateVMServiceScrape(ctx context.Context, scrape *types.VMServiceScrape) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodPost)
	req.SetRequestURI(c.getUrl("/vm-service-scrapes"))

	reqBody, err := json.Marshal(scrape)
	if err != nil {
		return err
	}
	req.SetBody(reqBody)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err = c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != consts.StatusCreated {
		return fmt.Errorf("failed to create vm service scrape: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) UpdateVMServiceScrape(ctx context.Context, name string, scrape *types.VMServiceScrape) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodPut)
	req.SetRequestURI(c.getUrl(fmt.Sprintf("/vm-service-scrapes/%s", name)))

	reqBody, err := json.Marshal(scrape)
	if err != nil {
		return err
	}
	req.SetBody(reqBody)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err = c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != consts.StatusOK {
		return fmt.Errorf("failed to update vm service scrape: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) DeleteVMServiceScrape(ctx context.Context, name string) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodDelete)
	req.SetRequestURI(c.getUrl(fmt.Sprintf("/vm-service-scrapes/%s", name)))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err := c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != consts.StatusOK {
		return fmt.Errorf("failed to delete vm service scrape: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) ListVMPodScrapes(ctx context.Context) ([]types.VMPodScrape, error) {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodGet)
	req.SetRequestURI(c.getUrl("/vm-pod-scrapes"))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err := c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != consts.StatusOK {
		return nil, fmt.Errorf("failed to list vm pod scrapes: %s", string(resp.Body()))
	}

	var scrapes []types.VMPodScrape
	if err := json.Unmarshal(resp.Body(), &scrapes); err != nil {
		return nil, err
	}

	return scrapes, nil
}

func (c *Client) GetVMPodScrape(ctx context.Context, name string) (*types.VMPodScrape, error) {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodGet)
	req.SetRequestURI(c.getUrl(fmt.Sprintf("/vm-pod-scrapes/%s", name)))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err := c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != consts.StatusOK {
		return nil, fmt.Errorf("failed to get vm pod scrape: %s", string(resp.Body()))
	}

	var scrape types.VMPodScrape
	if err := json.Unmarshal(resp.Body(), &scrape); err != nil {
		return nil, err
	}

	return &scrape, nil
}

func (c *Client) CreateVMPodScrape(ctx context.Context, scrape *types.VMPodScrape) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodPost)
	req.SetRequestURI(c.getUrl("/vm-pod-scrapes"))

	reqBody, err := json.Marshal(scrape)
	if err != nil {
		return err
	}
	req.SetBody(reqBody)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err = c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != consts.StatusCreated {
		return fmt.Errorf("failed to create vm pod scrape: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) UpdateVMPodScrape(ctx context.Context, name string, scrape *types.VMPodScrape) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodPut)
	req.SetRequestURI(c.getUrl(fmt.Sprintf("/vm-pod-scrapes/%s", name)))

	reqBody, err := json.Marshal(scrape)
	if err != nil {
		return err
	}
	req.SetBody(reqBody)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err = c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != consts.StatusOK {
		return fmt.Errorf("failed to update vm pod scrape: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) DeleteVMPodScrape(ctx context.Context, name string) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetMethod(consts.MethodDelete)
	req.SetRequestURI(c.getUrl(fmt.Sprintf("/vm-pod-scrapes/%s", name)))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	err := c.hertzClient.Do(ctx, req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != consts.StatusOK {
		return fmt.Errorf("failed to delete vm pod scrape: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) getUrl(path string) string {
	return fmt.Sprintf("http://%s%s%s", c.addr, apiPrefix, path)
}
