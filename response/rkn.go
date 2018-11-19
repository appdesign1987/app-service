package response

import "encoding/json"

// RknIpSuccess
type RknIpSuccess struct {
	Code      int    `json:"status_code"`
	IpAddress string `json:"ip_address"`
	IsBlocked bool   `json:"is_blocked"`
}

func (s RknIpSuccess) GetCode() int {
	return s.Code
}

func (s RknIpSuccess) GetResponse() string {
	b, err := json.Marshal(s)

	if err != nil {
		panic(err)
	}

	return string(b)
}

// RknIpsJsonSuccess
type RknIpsJsonSuccess struct {
	Code int             `json:"status_code"`
	Data map[string]bool `json:"data"`
}

func (s RknIpsJsonSuccess) GetCode() int {
	return s.Code
}

func (s RknIpsJsonSuccess) GetResponse() string {
	b, err := json.Marshal(s)

	if err != nil {
		panic(err)
	}

	return string(b)
}
