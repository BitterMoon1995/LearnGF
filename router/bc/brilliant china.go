/*
泪目！！！！！！！！！！！！！！！！！！
 */
package bc

import (
	"github.com/gogf/gf/net/ghttp"
)
func Index(request *ghttp.Request)  {
	request.Response.Writeln("最美中国 【首页】")
}
func Registry(request *ghttp.Request)  {
	request.Response.Writeln("最美中国 【用户注册】")
}
func Login(request *ghttp.Request)  {
	request.Response.Writeln("最美中国 【用户登录】")
}
func EditUser(request *ghttp.Request)  {
	request.Response.Writeln("最美中国 【修改用户】")
}
func SceneList(request *ghttp.Request)  {
	request.Response.Writeln("最美中国 【景区列表】")
}
func DelScene(request *ghttp.Request)  {
	request.Response.Writeln("最美中国 【删除景区】")
}