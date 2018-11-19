package response

import (
	"fmt"
	"net/http"
)

func SendResponse(w http.ResponseWriter, r Responsable) {
	w.WriteHeader(r.GetCode())
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, r.GetResponse())
}

func SendRawResponse(w http.ResponseWriter, Body string, Code int) {
	w.WriteHeader(Code)
	fmt.Fprint(w, Body)
}
