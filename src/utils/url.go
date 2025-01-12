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

func SplitHostPort(hostPort string) (host, port string) {
	host = hostPort

	colon := strings.LastIndexByte(host, ':')
	if colon != -1 && validOptionalPort(host[colon:]) {
		host, port = host[:colon], host[colon+1:]
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return
}

// validOptionalPort reports whether port is either an empty string
// or matches /^:\d*$/
func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}
