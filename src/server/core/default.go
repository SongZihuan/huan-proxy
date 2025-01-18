package core

import "net/http"

func (c *CoreServer) defaultResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
