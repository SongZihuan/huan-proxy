package api

const XHuanProxyHeaer = "X-Huan-Proxy"
const ViaHeader = "Via"
const TransferEncoding = "Transfer-Encoding"

var WarningReqHeader = []string{
	"Host",
	"Referer",
	"User-Agent",
	"Forwarded",
	"Content-Length",
	"Upgrade",
	"Connection",
	"X-Forwarded-For",
	"X-Forwarded-Host",
	"X-Forwarded-Proto",
	"X-Real-Ip",
	"X-Real-Port",
}

func isNotGoodHeader(header string) bool {
	for _, h := range WarningReqHeader {
		if h == header {
			return true
		}
	}

	return false
}
