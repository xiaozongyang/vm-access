package handler

import (
	"context"
	"fmt"
	"time"

	vm_operator "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/xiaozongyang/vm-access/internal/proxy/env"
	"github.com/xiaozongyang/vm-access/internal/proxy/vm_operator_client"
	"github.com/xiaozongyang/vm-access/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	VMNs = "prom"
)

func CreateVMStaticScrape(ctx context.Context, c *app.RequestContext) {
	var req types.VMStaticScrape
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	req.Meta.CreatedAt = time.Now()
	req.Meta.UpdatedAt = req.Meta.CreatedAt

	crd := staticScrapeToCRD(&req)

	vss, err := vm_operator_client.Client().VMStaticScrapes(VMNs).Create(ctx, crd, metav1.CreateOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create VMStaticScrape: %v", err)})
		return
	}
	hlog.Infof("created VMStaticScrape: %+v", vss)

	c.JSON(consts.StatusCreated, "ok")
}

func ListVMStaticScrapes(ctx context.Context, c *app.RequestContext) {
	vmStaticScrapes, err := vm_operator_client.Client().VMStaticScrapes(VMNs).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to list VMStaticScrapes: %v", err)})
		return
	}

	response := types.VMStaticScrapeListResponse{
		Items: make([]types.VMStaticScrape, 0),
	}

	for _, item := range vmStaticScrapes.Items {
		response.Items = append(response.Items, staticScrapeFormCRD(&item))
	}

	c.JSON(consts.StatusOK, response)
}

func GetVMStaticScrape(ctx context.Context, c *app.RequestContext) {
	name := c.Param("vm-static-scrape")

	crd, err := vm_operator_client.Client().VMStaticScrapes(VMNs).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to get VMStaticScrape %s: %v", name, err)})
		return
	}

	c.JSON(consts.StatusOK, staticScrapeFormCRD(crd))
}

func UpdateVMStaticScrape(ctx context.Context, c *app.RequestContext) {
	name := c.Param("vm-static-scrape")

	var vss types.VMStaticScrape
	if err := c.BindJSON(&vss); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if vss.Name != name {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": "name in request body and path parameter must be the same"})
		return
	}

	vss.Meta.UpdatedAt = time.Now()
	if vss.Meta.CreatedAt.IsZero() {
		vss.Meta.CreatedAt = vss.Meta.UpdatedAt
	}

	crd, err := vm_operator_client.Client().VMStaticScrapes(VMNs).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to get VMStaticScrape %s: %v", name, err)})
		return
	}

	updating := staticScrapeToCRD(&vss)
	updating.ObjectMeta.ResourceVersion = crd.ObjectMeta.ResourceVersion
	updating.ObjectMeta.Annotations = crd.ObjectMeta.Annotations

	_, err = vm_operator_client.Client().VMStaticScrapes(VMNs).Update(ctx, updating, metav1.UpdateOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to update VMStaticScrape: %v", err)})
		return
	}

	c.JSON(consts.StatusOK, map[string]string{"message": "ok"})
}

func DeleteVMStaticScrape(ctx context.Context, c *app.RequestContext) {
	name := c.Param("vm-static-scrape")

	err := vm_operator_client.Client().VMStaticScrapes(VMNs).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete VMStaticScrape %s: %v", name, err)})
		return
	}

	c.JSON(consts.StatusOK, map[string]string{"message": "ok"})
}

func staticScrapeToCRD(vss *types.VMStaticScrape) *vm_operator.VMStaticScrape {
	crd := &vm_operator.VMStaticScrape{
		ObjectMeta: metav1.ObjectMeta{
			Name:   vss.Name,
			Labels: map[string]string{},
		},
		Spec: vm_operator.VMStaticScrapeSpec{
			JobName: vss.JobName,
			TargetEndpoints: []*vm_operator.TargetEndpoint{
				{
					Targets: vss.Endpoint.Targets,
					Labels:  vss.Endpoint.Labels,
					EndpointScrapeParams: vm_operator.EndpointScrapeParams{
						Path: vss.Endpoint.GetPath(),
					},
				},
			},
		},
	}
	vss.Meta.InjectToLabels(&crd.ObjectMeta)
	return crd
}

func staticScrapeFormCRD(crd *vm_operator.VMStaticScrape) types.VMStaticScrape {
	meta := types.Meta{}
	meta.ExtractFromLabels(crd.ObjectMeta)
	return types.VMStaticScrape{
		Name:    crd.Name,
		JobName: crd.Spec.JobName,
		Cluster: env.GetCluster(),
		Meta:    meta,
		Endpoint: types.TargetEndpoint{
			Path:    crd.Spec.TargetEndpoints[0].Path,
			Labels:  crd.Spec.TargetEndpoints[0].Labels,
			Targets: crd.Spec.TargetEndpoints[0].Targets,
		},
	}
}
