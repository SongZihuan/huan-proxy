package respheader

const XHuanProxyHeaer = "X-Huan-Proxy"
const ViaHeader = "Via"
const ContentLength = "Content-Length"
const TransferEncoding = "Transfer-Encoding"

var WarningRespHeader = []string{
	"Host",
	"Referer",
	"User-Agent",
	"Upgrade",
	"Connection",
	"Cache-Control",
}

func isNotGoodHeader(header string) bool {
	for _, h := range WarningRespHeader {
		if h == header {
			return true
		}
	}

	return false
}
