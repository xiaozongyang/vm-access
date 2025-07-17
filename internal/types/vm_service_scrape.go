package types

type VMServiceScrape struct {
	Name      string            `json:"name"`
	Cluster   string            `json:"cluster"`
	Namespace string            `json:"namespace"`
	Selector  map[string]string `json:"selector"`

	Port     string `json:"port,omitempty"`
	Path     string `json:"path,omitempty"`
	JobLabel string `json:"jobLabel,omitempty"`

	// TODO: support serviceTargetLabels
	// TODO: support podTargetLabels

	Meta Meta `json:"meta"`
}
