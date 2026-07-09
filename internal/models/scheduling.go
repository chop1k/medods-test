package models

type ConnectivityCheck struct {
	IP     string `json:"ip"`
	Scheme string `json:"scheme"`
}
