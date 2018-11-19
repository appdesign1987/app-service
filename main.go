package main

import (
	"./rkn"
	"./whois"
	"fmt"
	Mux "github.com/gorilla/mux"
	Cache "github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var ListenPort = "8080";

	if os.Getenv("LISTEN_PORT") != "" {
		ListenPort = os.Getenv("LISTEN_PORT");
	}

	Store := Cache.New(1*time.Hour, 24*time.Hour)
	Store.Set("started", time.Now(), time.Hour)

	whois.Store = Store

	Router := Mux.NewRouter()

	Router.HandleFunc("/whois/{domain}", whois.DomainRouterHandler).Methods("GET")

	Router.HandleFunc("/rkn/ip/{ip}/short", rkn.IpShortRouterHandler).Methods("GET")
	Router.HandleFunc("/rkn/ip/{ip}", rkn.IpRouterHandler).Methods("GET")
	Router.HandleFunc("/rkn/ips", rkn.IpsJsonRouterHandler).Methods("POST")

	Router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	log.Printf("Listen on port: %s", ListenPort)

	http.ListenAndServe(fmt.Sprintf(":%s", ListenPort), Router)
}
