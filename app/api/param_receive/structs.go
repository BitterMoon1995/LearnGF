package param_receive

import (
	"github.com/gogf/gf/net/ghttp"
)

type registerRequest struct {
	Username string
	Password string
	// 给字段取请求别名
	CheckPassword string `p:"password2"`
}

// 不取别名的话可以不写JSON标签
type registerResponse struct {
	Code  int
	Error string
	Data  interface{} `json:"myData"`
}

/*
Request对象支持非常完美的请求校验能力，通过给结构体属性绑定v标签即可。
由于底层校验功能通过g valid模块实现，更详细的校验规则和介绍请参考 数据校验-结构体校验 章节。
*/
type verifiedRegRequest struct {
	Username string ` v:" required|length:6,12 # 请输入用户名|用户名最小:min位，最大max位 "`
	//要做比对（same:password）就得写别名
	Password      string ` p:"password" v:" required|length:6,18 # 请输入密码|密码最小:min位，最大max位 "`
	CheckPassword string ` v:" required|same:password # 请再次确认密码|两次输入密码不一致！ "`
}

/*
从v1.15版本开始，Request请求对象支持通过struct tag的方式为输入对象的属性绑定默认值。
默认值的struct tag名称为d。
*/
type pageQuery struct {
	Type       string // 内容模型
	CategoryId int    // 栏目ID
	PageNum    int    `d:"0"`  // 分页号码
	PageSize   int    `d:"10"` // 分页数量，最大50
	Sort       int    `d:"1"`  // 排序类型(0:最新, 默认。1:活跃, 2:热度)
}

func middleware1(request *ghttp.Request) {
	request.SetParam("name", "TS吴雨薇")
	request.Middleware.Next()
}

func middleware2(request *ghttp.Request) {
	request.SetParam("age", 998)
	request.Middleware.Next()
}

func contextMidWare(request *ghttp.Request) {
	request.SetCtxVar("userId", "tswei1995")
	request.Middleware.Next()
}
