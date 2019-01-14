package whois

import (
	Output "../response"
	"fmt"
	DateParse "github.com/araddon/dateparse"
	Mux "github.com/gorilla/mux"
	WhoisParser "github.com/likexian/whois-parser-go"
	Cache "github.com/patrickmn/go-cache"
	Whois "github.com/undiabler/golang-whois"
	"math"
	"net/http"
	"strings"
	"time"
)

var Store *Cache.Cache

func DomainRouterHandler(w http.ResponseWriter, r *http.Request) {
	RequestParams := Mux.Vars(r)
	DomainName := strings.ToLower(RequestParams["domain"])

	//if _d, found := Store.Get(DomainName); found {
	//	Data := _d.(Output.DomainData)
	//	Output.SendResponse(w, Output.WhoisSuccess{Code: http.StatusOK, Data: Data})
	//	return
	//}

	RawWhois, err := Whois.GetWhois(DomainName)

	if err != nil {
		Output.SendResponse(w, Output.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	ParsedWhois, err := WhoisParser.Parse(RawWhois)

	if err != nil {
		Output.SendResponse(w, Output.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	if ParsedWhois.Registrar.DomainName == "" {
		Output.SendResponse(w, Output.Error{Code: http.StatusInternalServerError, Message: "Unregistered Domain Name"})
		return
	}

	_c, err := DateParse.ParseLocal(ParsedWhois.Registrar.CreatedDate)

	if err != nil {
		Output.SendResponse(w, Output.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	_e, err := DateParse.ParseLocal(ParsedWhois.Registrar.ExpirationDate)

	if err != nil {
		Output.SendResponse(w, Output.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	Data := Output.DomainData{
		DomainName:     DomainName,
		NameServers:    strings.Split(ParsedWhois.Registrar.NameServers, ","),
		CreatedDate:    fmt.Sprintf("%d-%02d-%02d", _c.Year(), _c.Month(), _c.Day()),
		ExpirationDate: fmt.Sprintf("%d-%02d-%02d", _e.Year(), _e.Month(), _e.Day()),
		ExpireLeftDays: math.Round(_e.Sub(time.Now()).Hours() / 24),
	}

	// Store.SetDefault(DomainName, Data)

	Output.SendResponse(w, Output.WhoisSuccess{Code: http.StatusOK, Data: Data})
}
