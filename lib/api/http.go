package api

import "net/http"

type HTTP interface {
	Authorize(w http.ResponseWriter, r *http.Request) error
}
