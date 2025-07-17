package types

type VMStaticScrapeList struct {
	Items []VMStaticScrape `json:"items"`
}

type VMStaticScrape struct {
	Name      string `json:"name"`
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`

	JobName  string         `json:"jobName"`
	Endpoint TargetEndpoint `json:"endpoint"`

	Meta Meta `json:"meta"`
}

type TargetEndpoint struct {
	Path    string            `json:"path"`
	Labels  map[string]string `json:"labels"`
	Targets []string          `json:"targets"`
}

func (e *TargetEndpoint) GetPath() string {
	if e.Path != "" {
		return e.Path
	}
	return "/metrics"
}

type VMStaticScrapeListResponse struct {
	Items []VMStaticScrape `json:"items"`
}
