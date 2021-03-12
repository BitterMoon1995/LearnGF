package data_return

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"testing"
)

/*  缓冲控制

Response输出采用了缓冲控制，输出的内容预先写入到一块缓冲区，等待服务方法执行完毕后才真正地输出到客户端。
该特性在提高执行效率同时为输出内容的控制提供了更高的灵活性。

举个例子：
我们通过后置中间件统一对返回的数据做处理，如果服务方法产生了异常，那么不能将敏感错误信息推送到客户端，而统一设置错误提示信息。
*/
func TestBufferControl(t *testing.T) {
	server := g.Server()

	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(errorHandleMW)
		group.ALL("/err_handle", func(request *ghttp.Request) {
			panic("db error: bad sql")
		})
	})

	server.SetPort(4396)
	server.Run()
}

/*
json/xml支持
Response提供了对JSON/XML数据格式输出的原生支持，通过以下方法实现：

WriteJson* 方法用于返回JSON数据格式，参数为任意类型，可以为string、map、struct等等。返回的Content-Type为application/json。
WriteXml* 方法用于返回XML数据格式，参数为任意类型，可以为string、map、struct等等。返回的Content-Type为application/xml。

如果要返回对象，必须保证结构体的每个字段都是公开的，最好有JSON tag
*/
func TestReturnJsonXml(t *testing.T) {
	server := g.Server()

	server.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/json", func(request *ghttp.Request) {
			_ = request.Response.WriteJson(g.Map{
				"name": "尼哥",
				"age":  88,
			})
		})
		group.ALL("/xml", func(request *ghttp.Request) {
			_ = request.Response.WriteXml(g.Map{
				"name": "麻友",
				"age":  26,
			})
		})

		server.SetPort(4396)
		server.Run()
	})

}

/*
Exit控制
Exit: 仅退出当前执行的逻辑方法，不退出后续的请求流程，可用于替代return。
*/
func TestExitControl(t *testing.T) {
	server := g.Server()

	server.BindHandler("/exit", func(request *ghttp.Request) {
		nameId := request.GetInt("nameId")
		if nameId == 1 {
			request.Response.Writeln("black nigger")
			request.Exit()
		}
		request.Response.Writeln("slave")
	})

	server.SetPort(4396)
	server.Run()
}
