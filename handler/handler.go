package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

//HandleRequest handles incoming requests
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(400)
		_, _ = io.WriteString(w, "Only POST method is supported")
	}

	target := r.Header.Get("target")
	service := r.Header.Get("service")
	method := r.Header.Get("method")
	// TODO: more validation
	if len(method) > 0 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(501)
			return
		}
		err = invoke(w, target, service, method, convertQueryToMetadata(r.URL.Query()), string(body))
		if err != nil {
			logrus.WithError(err).Error("failed to invoke method")
			return
		}
	} else if len(service) > 0 {
		// TODO: list methods
	} else if len(target) > 0 {
		err := listServices(w, target)
		if err != nil {
			logrus.WithError(err).Error("failed to list services")
			return
		}
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
