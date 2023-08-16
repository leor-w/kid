package kid

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ResponseTypeKey = "resp_type"
)

type ResponseType int8

const (
	ResponseTypeJson           = iota + 1 // 1 json 类型
	ResponseTypeXml                       // 2 xml 类型
	ResponseTypeYaml                      // 3 yaml 类型
	ResponseTypeHtml                      // 4 html 类型
	ResponseTypeIndentedJson              // 5 格式化 json
	ResponseTypeSecureJson                // 6 安全的 json
	ResponseTypeJsonP                     // 7 返回 jsonp 格式
	ResponseTypeAsciiJson                 // 8 返回 Ascii 码的 json 会将 Unicode 转换为 ASCII
	ResponseTypePureJson                  // 9 pure json 格式化
	ResponseTypeToml                      // 10 分会 toml 格式
	ResponseTypeProtoBuf                  // 11 以 protobuf 格式输出
	ResponseTypeString                    // 12 以字符串格式输出
	ResponseTypeRedirect                  // 13 重定向
	ResponseTypeFile                      // 14 文件
	ResponseTypeFileFromFS                // 15 远程文件系统
	ResponseTypeFileAttachment            // 16 FileAttachment
	ResponseTypeSSEvent                   // 17 SSEvent
	ResponseTypeStream                    // 18 数据流
)

type RouterGroup struct {
	*gin.RouterGroup
}

type HandleFunc func(*Context) interface{}

type Middleware func(*Context)

func (group *RouterGroup) POST(path string, handler HandleFunc, middlewares ...Middleware) {
	group.RouterGroup.POST(path, convert(handler, middlewares...)...)
}

func (group *RouterGroup) GET(path string, handler HandleFunc, middlewares ...Middleware) {
	group.RouterGroup.GET(path, convert(handler, middlewares...)...)
}

func (group *RouterGroup) PUT(path string, handler HandleFunc, middlewares ...Middleware) {
	group.RouterGroup.PUT(path, convert(handler, middlewares...)...)
}

func (group *RouterGroup) DELETE(path string, handler HandleFunc, middlewares ...Middleware) {
	group.RouterGroup.DELETE(path, convert(handler, middlewares...)...)
}

func (group *RouterGroup) Group(path string, middlewares ...Middleware) *RouterGroup {
	return &RouterGroup{group.RouterGroup.Group(path, convertMiddleware(middlewares...)...)}
}

func (group *RouterGroup) UseMiddle(middlewares ...Middleware) {
	group.RouterGroup.Use(convertMiddleware(middlewares...)...)
}

func convert(handler HandleFunc, middlewares ...Middleware) []gin.HandlerFunc {
	h := convertHandleFunc(handler)
	ms := convertMiddleware(middlewares...)
	return append(ms, h)
}

func convertHandleFunc(handler HandleFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		kidCtx := &Context{Context: ctx}
		resp := handler(kidCtx)
		if resp == nil {
			ctx.Abort()
			return
		}
		responseData(ctx, kidCtx, resp)
	}
}

// responseData 通过获取 ctx 的 resp_type key 确定响应的格式类型
func responseData(ctx *gin.Context, kidCtx *Context, resp interface{}) {
	switch kidCtx.GetResponseType() {
	case ResponseTypeJson:
		ctx.JSON(http.StatusOK, resp)
	case ResponseTypeXml:
		ctx.XML(http.StatusOK, resp)
	case ResponseTypeYaml:
		ctx.YAML(http.StatusOK, resp)
	case ResponseTypeHtml:
		ctx.HTML(http.StatusOK, "", resp)
	case ResponseTypeIndentedJson:
		ctx.IndentedJSON(http.StatusOK, resp)
	case ResponseTypeSecureJson:
		ctx.SecureJSON(http.StatusOK, resp)
	case ResponseTypeJsonP:
		ctx.JSONP(http.StatusOK, resp)
	case ResponseTypeAsciiJson:
		ctx.AsciiJSON(http.StatusOK, resp)
	case ResponseTypePureJson:
		ctx.PureJSON(http.StatusOK, resp)
	case ResponseTypeToml:
		ctx.TOML(http.StatusOK, resp)
	case ResponseTypeProtoBuf:
		ctx.ProtoBuf(http.StatusOK, resp)
	case ResponseTypeString:
		ctx.String(http.StatusOK, "", resp)
	case ResponseTypeRedirect:
		url := resp.(string)
		ctx.Redirect(http.StatusOK, url)
	case ResponseTypeFile:
		filepath := resp.(string)
		ctx.File(filepath)
	default:
		ctx.JSON(http.StatusOK, resp)
	}
}

func convertMiddleware(middlewares ...Middleware) []gin.HandlerFunc {
	var ginMiddlewares = make([]gin.HandlerFunc, len(middlewares))
	for i := range middlewares {
		handler := middlewares[i]
		ginMiddlewares[i] = func(ctx *gin.Context) {
			handler(&Context{Context: ctx})
		}
	}
	return ginMiddlewares
}
