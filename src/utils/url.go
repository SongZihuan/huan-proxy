package utils

import "strings"

func ProcessPath(path string, defaultUrl ...string) string {
	if len(path) == 0 && len(defaultUrl) == 1 {
		path = defaultUrl[0]
	}

	path = strings.TrimSpace(path)

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = strings.TrimRight(path, "/")

	if !IsValidURLPath(path) {
		panic("A serious error occurred in 'ProcessPath', and the generated Path does not conform to the 'IsValidURLPath' validation logic.")
	}

	return path
}
