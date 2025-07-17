package vm_operator_client

import (
	"fmt"
	"os"

	vm_op_client "github.com/VictoriaMetrics/operator/api/client/versioned/typed/operator/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var vmOperatorClient *vm_op_client.OperatorV1beta1Client

func MustInitVMOperatorClient() {
	config, err := getRestConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to get rest config: %v", err))
	}

	vmOperatorClient = vm_op_client.NewForConfigOrDie(config)
}

func getRestConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		path := os.Getenv("KUBECONFIG")
		if path == "" {
			return nil, fmt.Errorf("KUBECONFIG is not set")
		}
		config, err = clientcmd.BuildConfigFromFlags("", path)
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes config: %w", err)
		}
	}
	return config, nil
}

func Client() *vm_op_client.OperatorV1beta1Client {
	return vmOperatorClient
}
