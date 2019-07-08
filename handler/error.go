package handler

import (
	"io"
	"net/http"
)

func reportError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	_, _ = io.WriteString(w, err.Error())
}
