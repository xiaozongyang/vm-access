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

func CreateVMServiceScrape(ctx context.Context, c *app.RequestContext) {
	var req types.VMServiceScrape
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	req.Meta.CreatedAt = time.Now()
	req.Meta.UpdatedAt = req.Meta.CreatedAt

	crd := serviceScrapeToCRD(&req)
	_, err := vm_operator_client.Client().VMServiceScrapes(VMNs).Create(ctx, crd, metav1.CreateOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create VMServiceScrape: %v", err)})
		return
	}

	c.JSON(consts.StatusCreated, req)
}

func GetVMServiceScrape(ctx context.Context, c *app.RequestContext) {
	name := c.Param("vm-service-scrape")

	crd, err := vm_operator_client.Client().VMServiceScrapes(VMNs).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to get VMServiceScrape %s: %v", name, err)})
		return
	}
	c.JSON(consts.StatusOK, serviceScrapeFormCRD(crd))
}

func ListVMServiceScrapes(ctx context.Context, c *app.RequestContext) {
	crds, err := vm_operator_client.Client().VMServiceScrapes(VMNs).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to list VMServiceScrapes: %v", err)})
		return
	}

	items := []types.VMServiceScrape{}
	for _, crd := range crds.Items {
		items = append(items, serviceScrapeFormCRD(&crd))
	}

	c.JSON(consts.StatusOK, items)
}

func UpdateVMServiceScrape(ctx context.Context, c *app.RequestContext) {
	name := c.Param("vm-service-scrape")
	var vss types.VMServiceScrape
	if err := c.BindJSON(&vss); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	vss.Meta.UpdatedAt = time.Now()
	if vss.Meta.CreatedAt.IsZero() {
		vss.Meta.CreatedAt = vss.Meta.UpdatedAt
	}

	crd, err := vm_operator_client.Client().VMServiceScrapes(VMNs).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to get VMServiceScrape %s: %v", name, err)})
		return
	}

	updating := serviceScrapeToCRD(&vss)
	updating.ObjectMeta.ResourceVersion = crd.ObjectMeta.ResourceVersion
	updating.ObjectMeta.Annotations = crd.ObjectMeta.Annotations

	_, err = vm_operator_client.Client().VMServiceScrapes(VMNs).Update(ctx, updating, metav1.UpdateOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to update VMServiceScrape %s: %v", name, err)})
		return
	}

	c.JSON(consts.StatusOK, "ok")
}

func DeleteVMServiceScrape(ctx context.Context, c *app.RequestContext) {
	name := c.Param("vm-service-scrape")

	err := vm_operator_client.Client().VMServiceScrapes(VMNs).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete VMServiceScrape %s: %v", name, err)})
		return
	}

	c.JSON(consts.StatusOK, "ok")
}

func serviceScrapeToCRD(vss *types.VMServiceScrape) *vm_operator.VMServiceScrape {
	crd := &vm_operator.VMServiceScrape{
		ObjectMeta: metav1.ObjectMeta{
			Name: vss.Name,
		},
		Spec: vm_operator.VMServiceScrapeSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: vss.Selector,
			},
			NamespaceSelector: vm_operator.NamespaceSelector{
				MatchNames: []string{vss.Namespace},
			},
			Endpoints: []vm_operator.Endpoint{
				{
					Port: vss.Port,
					EndpointScrapeParams: vm_operator.EndpointScrapeParams{
						Path: vss.Path,
					},
				},
			},
			JobLabel: vss.JobLabel,
		},
	}
	vss.Meta.InjectToLabels(&crd.ObjectMeta)
	return crd
}

func serviceScrapeFormCRD(crd *vm_operator.VMServiceScrape) types.VMServiceScrape {
	meta := types.Meta{}
	meta.ExtractFromLabels(crd.ObjectMeta)
	var ns string
	if len(crd.Spec.NamespaceSelector.MatchNames) > 0 {
		ns = crd.Spec.NamespaceSelector.MatchNames[0]
	} else {
		hlog.Warnf("namespaceSelector.MatchNames is empty, name: %s", crd.Name)
	}
	return types.VMServiceScrape{
		Name:      crd.Name,
		Cluster:   env.GetCluster(),
		Namespace: ns,
		Selector:  crd.Spec.Selector.MatchLabels,
		Port:      crd.Spec.Endpoints[0].Port,
		Path:      crd.Spec.Endpoints[0].EndpointScrapeParams.Path,
		JobLabel:  crd.Spec.JobLabel,
		Meta:      meta,
	}
}
