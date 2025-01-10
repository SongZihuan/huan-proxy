package utils

import (
	"net/http"
)

var statusOK = []int{
	http.StatusOK,
	http.StatusCreated,
	http.StatusAccepted,
	http.StatusNonAuthoritativeInfo,
	http.StatusNoContent,
	http.StatusResetContent,
	http.StatusPartialContent,
	http.StatusMultiStatus,
	http.StatusAlreadyReported,
}

func HttpStatusOK(status int) bool {
	for _, s := range statusOK {
		if status == s {
			return true
		}
	}
	return false
}
