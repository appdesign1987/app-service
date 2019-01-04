package response

import (
	"encoding/json"
	"net"
)

type SslData struct {
	HostName     string `json:"host"`
	LeftDays     int    `json:"left_days"`
	Certificates []Cert `json:"certs"`
}

type Cert struct {
	IpAddress  net.IP   `json:"ip"`
	LeftDays   int      `json:"left_days"`
	IssuerName string   `json:"issuer_name"`
	CommonName string   `json:"common_name"`
	Serial     string   `json:"serian_num"`
	DNSNames   []string `json:"dns_names"`
}

// SslHostJsonSuccess
type SslHostJsonSuccess struct {
	Code int     `json:"status_code"`
	Data SslData `json:"data"`
}

func (s SslHostJsonSuccess) GetCode() int {
	return s.Code
}

func (s SslHostJsonSuccess) GetResponse() string {
	b, err := json.Marshal(s)

	if err != nil {
		panic(err)
	}

	return string(b)
}
