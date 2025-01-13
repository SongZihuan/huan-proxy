package api

const XHuanProxyHeaer = "X-Huan-Proxy"
const ViaHeader = "Via"

var WarningHeader = []string{
	"Host",
	"Referer",
	"User-Agent",
	"Forwarded",
	"Content-Length",
	"Transfer-Encoding",
	"Upgrade",
	"Connection",
	"X-Forwarded-For",
	"X-Forwarded-Host",
	"X-Forwarded-Proto",
	"X-Real-Ip",
	"X-Real-Port",
}

func isNotGoodHeader(header string) bool {
	for _, h := range WarningHeader {
		if h == header {
			return true
		}
	}

	return false
}
