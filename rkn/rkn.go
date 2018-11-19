package rkn

import (
	Output "../response"
	"encoding/json"
	Mux "github.com/gorilla/mux"
	IpTree "github.com/zmap/go-iptree/iptree"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

const (
	DumpUrl        = "https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv"
	Timeout        = time.Second * 30
	RetryInterval  = time.Second * 30
	UpdateInterval = time.Minute * 15
)

var Db *IpTree.IPTree
var Lock = sync.Mutex{}

func init() {
	var err error

	CurrentDump := path.Join("/tmp", "rkn-current-dump.csv")
	FreshDump := path.Join("/tmp", "rkn-fresh-dump.csv")

	for {
		err := DownloadDump(DumpUrl, CurrentDump)

		if err == nil {
			break
		}

		log.Println(err, "Retry after", RetryInterval)

		time.Sleep(RetryInterval)
	}

	Db, err = DbLoadDump(CurrentDump)

	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		ticker := time.NewTicker(UpdateInterval).C

		for range ticker {
			log.Println("Updating RKN database...")

			err := DownloadDump(DumpUrl, FreshDump)

			if err != nil {
				log.Println(err)
				continue
			}

			Nb, err := DbLoadDump(FreshDump)

			if err != nil {
				log.Println(err)
				continue
			}

			Lock.Lock()

			Db = Nb

			if err = os.Rename(FreshDump, CurrentDump); err != nil {
				log.Println(err)
			}

			log.Println("RKN Database successful updated")

			Lock.Unlock()
		}
	}()
}

func IpRouterHandler(w http.ResponseWriter, r *http.Request) {
	RequestParams := Mux.Vars(r)

	IpAddress := net.ParseIP(RequestParams["ip"])
	v, ok, err := Db.GetByString(IpAddress.String())

	if err != nil || !ok {
		Output.SendResponse(w, Output.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	IsBlocked := false

	if i, _ := v.(int); i == 1 {
		IsBlocked = true
	}

	Output.SendResponse(w, Output.RknIpSuccess{Code: http.StatusOK, IpAddress: IpAddress.String(), IsBlocked: IsBlocked})
}

func IpShortRouterHandler(w http.ResponseWriter, r *http.Request) {
	RequestParams := Mux.Vars(r)

	IpAddress := net.ParseIP(RequestParams["ip"])
	v, ok, err := Db.GetByString(IpAddress.String())

	if err != nil || !ok {
		Output.SendRawResponse(w, "error", http.StatusInternalServerError)
		return
	}

	if i, _ := v.(int); i == 1 {
		Output.SendRawResponse(w, "block", 200)
	} else {
		Output.SendRawResponse(w, "clean", 200)
	}
}

func IpsJsonRouterHandler(w http.ResponseWriter, r *http.Request) {
	var RequestIps []string

	if err := json.NewDecoder(r.Body).Decode(&RequestIps); err != nil {
		Output.SendResponse(w, Output.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if len(RequestIps) < 1 {
		Output.SendResponse(w, Output.Error{Code: http.StatusBadRequest, Message: "Empty IP list json"})
		return
	}

	Response := make(map[string]bool, len(RequestIps))

	for _, ip := range RequestIps {
		v, ok, err := Db.GetByString(ip)

		if err != nil {
			Output.SendResponse(w, Output.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		Response[ip] = false

		if !ok {
			continue
		}

		if flag, _ := v.(int); flag == 1 {
			Response[ip] = true
		}
	}

	Output.SendResponse(w, Output.RknIpsJsonSuccess{Code: http.StatusOK, Data: Response})
}
