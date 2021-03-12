package data_return

import (
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

func errorHandleMW(request *ghttp.Request) {

	request.Middleware.Next()
	if request.Response.Status >= http.StatusInternalServerError {
		request.Response.ClearBuffer()
		request.Response.Writeln("服务器开小差了ε=(´ο｀*)))")
		//request.Response.Flush()
	}
}
