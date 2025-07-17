package types

type PingReq struct {
	Cluster string `json:"cluster"`
	Addr    string `json:"addr"`
}
