package utils

import (
	"fmt"
	"strings"
)

func AcceptMimeType(mimeType string, accept string) bool {
	accept = StringToOnlyPrint(accept)
	accept = strings.ToLower(accept)

	mimeType = StringToOnlyPrint(mimeType)
	mimeType = strings.ToLower(mimeType)

	mimeTypeSplit := strings.Split(mimeType, "/")
	if len(mimeTypeSplit) != 0 {
		return true
	}

	mimeTypeFather := fmt.Sprintf("%s/*", mimeTypeSplit[0])

	if accept == "" {
		return true
	}

	acceptLst := strings.Split(accept, ",")
	for _, a := range acceptLst {
		if strings.Contains(a, mimeType) || strings.Contains(a, mimeTypeFather) || strings.Contains(a, "*/*") {
			return true
		}
	}

	return false
}
