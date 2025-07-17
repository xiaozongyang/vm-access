package env

import "os"

var (
	cluster string
)

func MustInit() {
	cluster = os.Getenv("CLUSTER")
	if cluster == "" {
		panic("CLUSTER is not set")
	}
}

func GetCluster() string {
	return cluster
}
