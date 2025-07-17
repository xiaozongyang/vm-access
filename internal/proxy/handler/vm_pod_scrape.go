package handler

import (
	"context"
	"time"

	vm_operator "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/xiaozongyang/vm-access/internal/proxy/env"
	"github.com/xiaozongyang/vm-access/internal/proxy/vm_operator_client"
	"github.com/xiaozongyang/vm-access/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateVMPodScrape(ctx context.Context, c *app.RequestContext) {
	var req types.VMPodScrape
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	req.Meta.CreatedAt = time.Now()
	req.Meta.UpdatedAt = req.Meta.CreatedAt

	crd := podScrapeToCRD(&req)

	_, err := vm_operator_client.Client().VMPodScrapes(VMNs).Create(ctx, crd, metav1.CreateOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusCreated, map[string]string{"message": "ok"})
}

func UpdateVMPodScrape(ctx context.Context, c *app.RequestContext) {
	var req types.VMPodScrape
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	crd, err := vm_operator_client.Client().VMPodScrapes(VMNs).Get(ctx, c.Param("vm-pod-scrape"), metav1.GetOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	req.Meta.UpdatedAt = time.Now()
	if req.Meta.CreatedAt.IsZero() {
		req.Meta.CreatedAt = crd.ObjectMeta.CreationTimestamp.Time
	}

	updating := podScrapeToCRD(&req)
	updating.ObjectMeta.ResourceVersion = crd.ObjectMeta.ResourceVersion
	updating.ObjectMeta.Annotations = crd.ObjectMeta.Annotations

	_, err = vm_operator_client.Client().VMPodScrapes(VMNs).Update(ctx, updating, metav1.UpdateOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, "ok")
}

func DeleteVMPodScrape(ctx context.Context, c *app.RequestContext) {
	err := vm_operator_client.Client().VMPodScrapes(VMNs).Delete(ctx, c.Param("vm-pod-scrape"), metav1.DeleteOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, "ok")
}

func GetVMPodScrape(ctx context.Context, c *app.RequestContext) {
	crd, err := vm_operator_client.Client().VMPodScrapes(VMNs).Get(ctx, c.Param("vm-pod-scrape"), metav1.GetOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, podScrapeFormCRD(crd))
}

func ListVMPodScrapes(ctx context.Context, c *app.RequestContext) {
	crds, err := vm_operator_client.Client().VMPodScrapes(VMNs).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var podScrapes []types.VMPodScrape
	for _, crd := range crds.Items {
		podScrapes = append(podScrapes, podScrapeFormCRD(&crd))
	}

	c.JSON(consts.StatusOK, podScrapes)
}

func podScrapeToCRD(vps *types.VMPodScrape) *vm_operator.VMPodScrape {
	var portNumber *int32
	if vps.PortNumber > 0 {
		portNumber = new(int32)
		*portNumber = int32(vps.PortNumber)
	}

	crd := &vm_operator.VMPodScrape{
		ObjectMeta: metav1.ObjectMeta{
			Name: vps.Name,
		},
		Spec: vm_operator.VMPodScrapeSpec{
			Selector: labelSelectorToCRD(vps.Selector),
			JobLabel: vps.JobLabel,
			NamespaceSelector: vm_operator.NamespaceSelector{
				MatchNames: []string{vps.Namespace},
			},
			PodMetricsEndpoints: []vm_operator.PodMetricsEndpoint{
				{
					PortNumber: portNumber,
					EndpointScrapeParams: vm_operator.EndpointScrapeParams{
						Path: vps.Path,
					},
					Port: &vps.Port,
				},
			},
		},
	}
	vps.Meta.InjectToLabels(&crd.ObjectMeta)
	return crd
}

func podScrapeFormCRD(crd *vm_operator.VMPodScrape) types.VMPodScrape {
	portNumber := 0
	if crd.Spec.PodMetricsEndpoints[0].PortNumber != nil {
		portNumber = int(*crd.Spec.PodMetricsEndpoints[0].PortNumber)
	}
	var port string
	if crd.Spec.PodMetricsEndpoints[0].Port != nil {
		port = *crd.Spec.PodMetricsEndpoints[0].Port
	}
	var namespace string
	if len(crd.Spec.NamespaceSelector.MatchNames) > 0 {
		namespace = crd.Spec.NamespaceSelector.MatchNames[0]
	}

	var meta types.Meta
	meta.ExtractFromLabels(crd.ObjectMeta)

	return types.VMPodScrape{
		Name:      crd.Name,
		Cluster:   env.GetCluster(),
		Namespace: namespace,

		Selector: labelSelectorFromCRD(crd.Spec.Selector),

		Port:       port,
		PortNumber: portNumber,
		Path:       crd.Spec.PodMetricsEndpoints[0].EndpointScrapeParams.Path,
		JobLabel:   crd.Spec.JobLabel,

		Meta: meta,
	}
}

func labelSelectorToCRD(selector types.LabelSelector) metav1.LabelSelector {
	return metav1.LabelSelector{
		MatchLabels:      selector.MatchLabels,
		MatchExpressions: expressionsToLabelSelectorRequirements(selector.MatchExpressions),
	}
}

func labelSelectorFromCRD(selector metav1.LabelSelector) types.LabelSelector {
	return types.LabelSelector{
		MatchLabels:      selector.MatchLabels,
		MatchExpressions: expressionsFromLabelSelectorRequirements(selector.MatchExpressions),
	}
}

func expressionsFromLabelSelectorRequirements(requirements []metav1.LabelSelectorRequirement) []types.Expression {
	var expressions []types.Expression
	for _, requirement := range requirements {
		expressions = append(expressions, types.Expression{
			Key:      requirement.Key,
			Operator: string(requirement.Operator),
			Values:   requirement.Values,
		})
	}
	return expressions
}

func expressionsToLabelSelectorRequirements(expressions []types.Expression) []metav1.LabelSelectorRequirement {
	var requirements []metav1.LabelSelectorRequirement
	for _, expression := range expressions {
		requirements = append(requirements, metav1.LabelSelectorRequirement{
			Key:      expression.Key,
			Operator: metav1.LabelSelectorOperator(expression.Operator),
			Values:   expression.Values,
		})
	}
	return requirements
}
