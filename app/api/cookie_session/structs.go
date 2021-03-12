package cookie_session

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func errorHandler(request *ghttp.Request) {
	request.Middleware.Next()
	if err := request.GetError(); err != nil {
		g.Log("exception").Error(err)
		request.Response.ClearBuffer()
		request.Response.Writeln("服务器开小差了ε=(´ο｀*)))")
	}
}
