package response

import "encoding/json"

type DomainData struct {
	DomainName     string   `json:"domain_name"`
	CreatedDate    string   `json:"created_date"`
	ExpirationDate string   `json:"expiration_date"`
	ExpireLeftDays float64  `json:"expire_left_days"`
	NameServers    []string `json:"name_servers"`
}

// WhoisSuccess
type WhoisSuccess struct {
	Code int        `json:"status_code"`
	Data DomainData `json:"domain_data"`
}

func (s WhoisSuccess) GetCode() int {
	return s.Code
}

func (s WhoisSuccess) GetResponse() string {

	b, err := json.Marshal(s)

	if err != nil {
		panic(err)
	}

	return string(b)
}
