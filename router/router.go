package router

import (
    "LearnGF/app/api/hello"
    "github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()
	s.SetPort(9000)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/hello", hello.Hello)
	})
}
