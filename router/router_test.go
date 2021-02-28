package router

import (
	"LearnGF/router/bc"
	"fmt"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"testing"
)

func returnReqUri(request *ghttp.Request)  {
	request.Response.Writeln(request.Router.Uri)
}

/*
路由注册	pattern的写法规范、3种动态路由
[HTTPMethod:]路由规则[@域名]
 */
func TestPatternRegister(t *testing.T) {
	server := g.Server()
	// 精准匹配规则，该层级必须为定死的值
	server.BindHandler("/user", returnReqUri)

	/* 如果给定HTTPMethod，那么路由规则仅会在该请求方式下有效。@域名可以指定生效的域名名称，那么该路由规则仅会在该域名下生效。
	例如此时只有域名127.0.0.1的post请求有效，localhost都不好使*/
	server.BindHandler("POST:/student@127.0.0.1", returnReqUri)

	/*盘点3种动态路由规则*/
	// 命名匹配规则：该层级必须有值，且层级数也必须一致
	server.BindHandler("/name/:path", returnReqUri)
	server.BindHandler("/:path/name", returnReqUri)
	// 模糊匹配规则：该层级可以为空，且层级数也可以不一致
	server.BindHandler("/fuzzy/*path", returnReqUri)
	server.BindHandler("/nigger/*fuzzy/:name", returnReqUri)
	/* 字段匹配规则：使用{field}方式进行匹配，可对URI任意位置的参数进行截取匹配，
	该URI层级必须有值且形式一致，且层级数也必须一致*/
	server.BindHandler("/user/list/{username}", returnReqUri)
	server.BindHandler("/page/list/{pageName}.html", returnReqUri)
	server.BindHandler("/db-{dbNum}/:tableName", returnReqUri)


	server.SetPort(9000)
	server.Run()
}

type myController struct {
	total *gtype.Int
}

func (myController myController) Total(request *ghttp.Request)  {
	request.Response.Writeln("total:",myController.total.Add(1))
}

type myHandlerObject struct {
	count *gtype.Int
}

func (receiver *myHandlerObject) ShowAdd(request *ghttp.Request)  {
	request.Response.Writeln("add 1,current count:",receiver.count.Add(1))
}

func (receiver *myHandlerObject) ShowMinus(request *ghttp.Request){
	request.Response.Writeln("minus 1,current count:",receiver.count.Add(-1))
}

func (receiver *myHandlerObject) Index(request *ghttp.Request)  {
	request.Response.Writeln("myHandlerObjectの默认方法：杀死那尼哥")
}
// 对象中的Init和Shut是两个在HTTP请求流程中被Server自动调用的特殊方法（类似构造函数和析构函数的作用）。
func (receiver myHandlerObject) Init(request *ghttp.Request)  {
	fmt.Println("我是myHandlerObject的Init方法")
	request.Response.Writeln("我是myHandlerObject的Init方法")
}
func (receiver myHandlerObject) Shut(request *ghttp.Request)  {
	fmt.Println("我是myHandlerObject的Shut方法")
	request.Response.Writeln("我是myHandlerObject的Shut方法")
}

/*
路由注册	handler的3种注册方式
 */
func TestHandlerRegister(t *testing.T) {
	/* 1.通过BindHandler方法进行（回调）函数注册①匿名函数、②包方法、③对象方法，
	匿名函数和包方法不解释，看对象方法*/
	server := g.Server()

	//gtype.Int是一个结构体，必须先初始化分配空间
	mc := &myController{total: gtype.NewInt()}
	server.BindHandler("/total",mc.Total)

	/* 2.通过BindObject方法及其变种进行handler对象注册，该对象常驻内存不释放。服
	务端进程在启动时便需要初始化这些对象，并且这些对象需要自行负责自身数据的并发安全。

	对象的方法当然都得是HandlerFunc，并且请求对应的方法名（全小写）才会调用对应的方法，
	除了Index方法，/object/index 与/object等效

	★如果方法名或对象名为驼峰式复杂命名，那么对应的请求路径默认为：
	依然全小写；使用 - 连接各个单词。 比如这里为：http://localhost:9000/my-handler-object/show-add
	此外，我们可以通过.Server.SetNameToUriType方法来设置对象方法名称的4种路由生成方式。*/

	myObject := &myHandlerObject{count: gtype.NewInt()}
	server.BindObject("/object", myObject)
	//server.SetNameToUriType(3)//驼峰式

	/* 当使用BindObject方法进行对象注册时，在路由规则中可以使用两个内置的变量：
	{.struct}和{.method}，前者表示当前对象名称，后者表示当前注册的方法名。*/
	server.BindObject("/{.struct}/{.method}", myObject)

	server.SetPort(9000)
	server.Run()
}

/*
分组路由 基本使用
 */
func TestGroupedRouter(t *testing.T) {
	server := g.Server()
	userGroup := server.Group("/user")
	//仅允许GET方式访问
	userGroup.GET("/getUserList", func(request *ghttp.Request) {
		request.Response.Writeln("所有用户列表.........")
	})
	userGroup.POST("/login", func(request *ghttp.Request) {
		request.Response.Writeln("用户登录......")
	})
	//允许所有的方式访问
	userGroup.ALL("/makeLove", func(request *ghttp.Request) {
		request.Response.Writeln("用户左爱......")
	})
	userGroup.ALL("/sleep", func(request *ghttp.Request) {
		request.Response.Writeln("用户睡觉......")
	})
	server.SetPort(9000)
	server.Run()
}

/*
分组路由 层级路由注册（注意.Group方法参数的变化）
GF框架的分组路由注册支持更加直观优雅层级的注册方式，以便于开发者更方便地管理路由列表。路由层级注册方式也是推荐的路由注册方式。
 */
func TestHierarchyRouter(t *testing.T) {
	server := g.Server()
	server.Group("/brilliantCN.com", func(group *ghttp.RouterGroup) {
		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.POST("/login", bc.Login)
			group.PUT("/edit",bc.EditUser)
		})
		group.Group("/scene", func(group *ghttp.RouterGroup) {
			group.GET("/sceneList",bc.SceneList)
			group.DELETE("/delScene",bc.DelScene)
		})
	})
	server.SetPort(9000)
	server.Run()
}
/*
分组路由 使用ALLMap方法实现批量的路由注册
 */
func TestBatchRegister(t *testing.T) {
	server := g.Server()

	server.Group("/", func(group *ghttp.RouterGroup) {
		group.ALLMap(g.Map{
			"/" : bc.Index,
			"/registry" : bc.Registry,
			"/login" : bc.Login,
			"/editUser" : bc.EditUser,
			"/sceneList" : bc.SceneList,
			"/delScene" : bc.DelScene,
		})
	})

	server.SetPort(9000)
	server.Run()
}

/*
中间件（拦截器） 中间件编写与注册。中间件分为前置与后置，全局与分组
 */
func globalPre(request *ghttp.Request)  {
	fmt.Println("我是全局前置中间件")
	//放行
	request.Middleware.Next()
}
func globalPost(request *ghttp.Request)  {
	request.Middleware.Next()
	fmt.Println("我是全局后置中间件")
}
func girlMidWare(request *ghttp.Request)  {
	fmt.Println("我是女孩子的中间件")
	request.Middleware.Next()
}
func girlMidWare2(request *ghttp.Request)  {
	fmt.Println("我是女孩子的另一个中间件")
	request.Middleware.Next()
}
func TestGlobalMiddleware(t *testing.T) {
	server := g.Server()

	server.BindHandler("/user",returnReqUri)
	server.BindHandler("/nigger",returnReqUri)

	/* 全局中间件可通过Server或Domain对象进行绑定
	执行顺序规则：1.当同一个路由匹配到多个中间件时，会按照路由的深度优先规则执行
	2.同一个路由规则下，会按照中间件的注册先后顺序执行*/

	//仅对‘/user’适用globalPre
	server.BindMiddleware("/user",globalPre)
	//对所有请求适用globalPost
	server.BindMiddlewareDefault(globalPost)

	server.Group("/girl", func(group *ghttp.RouterGroup) {
		/* 分组中间件仅由一个group.Middleware方法进行绑定
		执行顺序规则：只会按照注册的先后顺序执行*/

		group.Middleware(girlMidWare2)
		//对/girl路由组下所有请求适用girlMidWare
		group.Middleware(girlMidWare)
		group.ALLMap(g.Map{
			"/skirt": returnReqUri,
		})
	})

	server.SetPort(9000)
	server.Run()
}
/*
中间件 典型案例
允许跨域请求 请求鉴权处理 鉴权例外处理 统一的错误处理 自定义日志处理，详询文档
 */


