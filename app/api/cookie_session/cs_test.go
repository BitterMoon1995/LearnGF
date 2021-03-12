package cookie_session

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
	"github.com/gogf/gf/os/gtime"
	"testing"
	"time"
)

func TestCookie(t *testing.T) {
	server := g.Server()
	g.Redis()

	server.BindHandler("/cookie", func(request *ghttp.Request) {
		cookie := request.Cookie
		// cookie-session登录方案的常规操作：简单设个sessionID
		cookie.SetSessionId("20181103ig")
		cookie.Set("cur_time", gtime.Datetime())
		fmt.Println(cookie)
		_ = request.Response.WriteJson(cookie)
	})

	server.SetPort(1103)
	server.Run()
}

/*
基于文件存储的session（默认）

使用场景：单机；多读少写；持久化（只要没过期就可以从文件中恢复）
*/
func TestFileSession(t *testing.T) {
	server := g.Server()
	_ = server.SetConfigWithMap(g.Map{
		//为方便观察设置2分钟过期
		"SessionMaxAge": time.Minute * 2,
	})

	server.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/set", func(request *ghttp.Request) {
			_ = request.Session.Set("time", gtime.Timestamp())
		})
		group.ALL("/map", func(request *ghttp.Request) {
			sessionMap := request.Session.Map()
			request.Response.Write(sessionMap)
		})
		group.ALL("/clear", func(request *ghttp.Request) {
			_ = request.Session.Clear()
		})
	})

	server.SetPort(1984)
	server.Run()
}

func TestRedisSession(t *testing.T) {
	server := g.Server()
	_ = server.SetConfigWithMap(g.Map{
		"SessionMaxAge":  time.Hour * 2,
		"SessionStorage": gsession.NewStorageRedis(g.Redis()),
		"sessionIdName":  "Jsessionid",
		"sessionPath":    "/temp/MySessionStoragePath",
	})

	server.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/set", func(request *ghttp.Request) {
			_ = request.Session.Set("time", gtime.Timestamp())
		})
		group.ALL("/get", func(request *ghttp.Request) {
			timeD := request.Session.Get("time")
			request.Response.Write(timeD)
		})
		group.ALL("/remove", func(request *ghttp.Request) {
			_ = request.Session.Remove("time")
		})
	})

	server.SetPort(1984)
	server.Run()
}

func TestErrorHandler(t *testing.T) {
	server := g.Server()

	server.Use(errorHandler)
	server.BindHandler("/getUserInfo", func(request *ghttp.Request) {
		panic("数据库down了")
	})

	server.SetPort(1984)
	server.Run()
}
