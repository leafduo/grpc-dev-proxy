package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

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
		output, err := client.Invoke(target, service, method, convertQueryToMetadata(r.URL.Query()), string(body))
		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, err.Error())
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

func convertQueryToMetadata(query url.Values) []string {
	metadata := make([]string, 0, len(query))
	for key, value := range query {
		m := fmt.Sprintf("%s: %s", key, strings.Join(value, ","))
		metadata = append(metadata, m)
	}
	return metadata
}
