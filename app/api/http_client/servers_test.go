package http_client

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"testing"
)

func TestNigger(t *testing.T) {
	server := g.Server()

	server.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/get", func(request *ghttp.Request) {
			name := request.GetString("name")
			age := request.GetInt("age")
			faceScore := request.GetFloat32("face_score")
			isBeauty := request.GetBool("is_beauty")

			girl := &girl{
				Name:      name,
				Age:       age,
				FaceScore: faceScore,
				IsBeauty:  isBeauty,
			}
			_ = request.Response.WriteJsonExit(girl)
		})

		group.POST("/post_map", func(request *ghttp.Request) {
			name := request.Get("name")
			age := request.Get("age")
			fmt.Println(name, age)
		})

		group.POST("/post_json", func(request *ghttp.Request) {
			fmt.Println(request.Header["Content-Type"])
			name := request.Get("name")
			age := request.Get("age")
			fmt.Println(name, age)
		})

		group.GET("/getIntList", func(request *ghttp.Request) {
			ints := []int{1, 2, 3, 4, 5, 6}
			request.Response.Write(ints)
		})
	})

	server.SetPort(1937)
	server.Run()
}
