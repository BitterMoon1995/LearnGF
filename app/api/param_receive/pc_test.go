package param_receive

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"testing"
)

/*
参数接收，提交格式规则
*/
func TestReceiveSameName(t *testing.T) {
	server := g.Server()

	/* r.Get("key")
	若有同名参数提交格式形如：k=v1&k=v2 ，后续的变量值将会覆盖前面的变量值。*/
	server.BindHandler("/kv", func(r *ghttp.Request) {
		name := r.Get("name")
		r.Response.Writeln(name)
	})

	/* 数组参数 提交格式形如：k[]=v1&k[]=v2
	这里请求路径应为：/arr?girlList[]=薇儿,abc,安吉拉掰逼 */
	server.BindHandler("/arr", func(r *ghttp.Request) {
		girl := r.Get("girlList")
		r.Response.Writeln(girl)
	})

	/* Map参数提交格式形如：k[a]=m&k[b]=n，并且支持多级Map，例如：k[a][a]=m&k[a][b]=n。
	这里路径应为： /map?weier[name]=薇柔芷&weier[age]=18 （麻了）*/
	server.BindHandler("/map", func(r *ghttp.Request) {
		girl := r.Get("weier")
		r.Response.Writeln(girl)
	})

	server.SetPort(9000)
	server.Run()
}

/*
对象处理：对象转换
规则：
struct中需要匹配的属性必须为公开属性(首字母大写)。
参数名称会自动按照 不区分大小写 且 忽略-/_/空格符号 的形式与struct属性进行匹配。
如果匹配成功，那么将键值赋值给属性，如果无法匹配，那么忽略该键值。
模拟场景：
用户注册，校验用户两次输入的密码是否相同，并JSON返回校验结果
*/
func TestObjectTransfer(t *testing.T) {
	server := g.Server()

	server.BindHandler("/register", func(request *ghttp.Request) {
		var regRequest *registerRequest
		/* Parse(pointer interface{})尝试将请求对象解析为pointer所对应的结构体类型，并将结果赋给该指针；
		这个错误判断用于判断请求对象是否符合该结构体的校验条件*/
		if err := request.Parse(&regRequest); err != nil {
			_ = request.Response.WriteJsonExit(registerResponse{
				Code:  998,
				Error: "校验失败",
			})
		}

		//调用土耳其野男人的包验空
		hasZero := structs.HasZero(regRequest)
		fmt.Println(hasZero)
		if hasZero {
			_ = request.Response.WriteJsonExit(registerResponse{
				Code:  400,
				Error: "参数为空",
			})
		}
		if regRequest.Password != regRequest.CheckPassword {
			_ = request.Response.WriteJsonExit(registerResponse{
				Code:  402,
				Error: "两次密码不一致",
			})
		}

		_ = request.Response.WriteJsonExit(registerResponse{
			Code: 200,
			Data: regRequest,
		})
	})

	server.SetPort(9000)
	server.Run()
}

/*
♥对象处理：请求对象校验、校验错误处理
通过给结构体属性绑定v标签，可以限制请求对象字段的存在与否、长度、是否相同等等，
详询 数据校验-结构体校验 章节
*/
func TestObjectVerify(t *testing.T) {
	server := g.Server()

	server.BindHandler("/register", func(request *ghttp.Request) {
		verifiedRR := new(verifiedRegRequest)
		if err := request.Parse(&verifiedRR); err != nil {
			/* 当请求校验错误时，所有校验失败的错误都返回了，这样对于用户体验不是特别友好。
			当错误产生后，我们可以通过err.(*gvalid.Error)断言的方式判断错误是否为校验错误，
			如果是的话则返回第一条校验错误，而不是所有都返回。*/

			//断言成功说明是校验错误
			if errVal, ok := err.(*gvalid.Error); ok {
				_ = request.Response.WriteJsonExit(registerResponse{
					Code:  444,
					Error: errVal.FirstString(),
				})
				//否则是其他的错误
			} else {
				_ = request.Response.WriteJsonExit(registerResponse{
					Code:  500,
					Error: "非校验错误",
				})
			}
		}
		_ = request.Response.WriteJsonExit(registerResponse{
			Code: 200,
			Data: verifiedRR,
		})
	})
	/* 此外，我们这里也可以使用gerror.Current获取第一条报错信息，而不是使用断言判断*/
	server.BindHandler("/register2", func(request *ghttp.Request) {
		verifiedRR := new(verifiedRegRequest)
		if err := request.Parse(&verifiedRR); err != nil {
			_ = request.Response.WriteJsonExit(registerResponse{
				Code:  444,
				Error: gerror.Current(err).Error(),
			})
		}
	})

	server.SetPort(9000)
	server.Run()
}

/*
自定义变量
开发者可以在请求中自定义一些变量设置，自定义变量的获取优先级是最高的，可以覆盖原有的客户端提交参数。
自定义变量往往也可以做请求流程的变量共享，但是需要注意的是该变量会成为请求参数的一部分，是对业务执行流程公开的变量。
*/
func TestCustomParam(t *testing.T) {
	server := g.Server()

	server.Group("/girl", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware1, middleware2)
		group.ALL("/info", func(request *ghttp.Request) {
			name := request.GetParamVar("name")
			age := request.GetParamVar("age")
			request.Response.Writefln("%s  %s",
				name.String(),
				age.String())
		})
	})

	server.SetPort(4396)
	server.Run()
}

/*
上下文
在GF框架中，我们推荐使用Context上下文对象来处理流程共享的上下文变量，
甚至将该对象进一步传递到依赖的各个模块方法中。该Context对象类型实现了标准库的context.Context接口，
该接口往往会作为模块间调用方法的第一个参数，该接口参数也是Golang官方推荐的在模块间传递上下文变量的推荐方式。
*/
func TestContext(t *testing.T) {
	server := g.Server()

	server.Group("/user", func(group *ghttp.RouterGroup) {
		group.Middleware(contextMidWare)
		group.ALL("/id", func(request *ghttp.Request) {
			userId := request.GetCtxVar("userId")
			request.Response.Writeln(userId)
		})
	})

	server.SetPort(4396)
	server.Run()
}
