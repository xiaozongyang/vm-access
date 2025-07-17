package types

type VMPodScrape struct {
	Name      string `json:"name"`
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`

	Selector LabelSelector `json:"selector"`

	Port       string `json:"port,omitempty"`
	PortNumber int    `json:"portNumber,omitempty"`
	Path       string `json:"path,omitempty"`
	JobLabel   string `json:"jobLabel,omitempty"`

	// TODO: support podTargetLabels
	// PodTargetLabels []string `json:"podTargetLabels,omitempty"`

	// TODO: support relabel

	Meta Meta `json:"meta"`
}

type LabelSelector struct {
	MatchLabels      map[string]string `json:"matchLabels,omitempty"`
	MatchExpressions []Expression      `json:"matchExpressions,omitempty"`
}

type Expression struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}
