package core

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/server/context"
	"net/http"
)

func (c *CoreServer) CoreServeHTTP(writer http.ResponseWriter, r *http.Request) {
	func() {
	RuleCycle:
		for _, rule := range c.GetRulesList() {
			if !c.matchURL(rule, r) {
				continue RuleCycle
			}

			ctx := context.NewContext(rule, writer, r)

			if !c.checkProxyTrust(ctx) {
				return
			}

			c.writeHuanProxyHeader(ctx)

			if rule.Type == rulescompile.ProxyTypeFile {
				c.fileServer(ctx)
			} else if rule.Type == rulescompile.ProxyTypeDir {
				c.dirServer(ctx)
			} else if rule.Type == rulescompile.ProxyTypeAPI {
				c.apiServer(ctx)
			} else if rule.Type == rulescompile.ProxyTypeRedirect {
				c.redirectServer(ctx)
			} else {
				c.abortServerError(ctx)
			}

			fmt.Printf("ctx.Abort: %v\n", ctx.Abort)
			fmt.Printf("config.GetConfig().NotAbort.IsEnable(false) = %v\n", config.GetConfig().NotAbort.IsEnable(false))
			if ctx.Abort && config.GetConfig().NotAbort.IsEnable(false) {
				_ = ctx.Reset()
				continue RuleCycle
			}

			fmt.Println("TAG 6")
			ctx.MustWriteToResponse()
			return

		}

		c.defaultResponse(writer)
		// 此处虽然w为Writer，但应该交由LoggerServer来处理写入
	}()
}
