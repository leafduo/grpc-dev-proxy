package handler

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/leafduo/grpc-dev-proxy/client"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	target := r.Header.Get("target")
	service := r.Header.Get("service")
	method := r.Header.Get("method")
	// TODO: more validation
	if len(method) > 0 {
		// TODO: describe method if it's GET request
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		output, err := client.Invoke(target, service, method, string(body))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		io.WriteString(w, output)
	} else if len(service) > 0 {
		// TODO: list methods
	} else if len(target) > 0 {
		// TODO: list services
	} else {
		// TODO: print some help message?
	}
}
