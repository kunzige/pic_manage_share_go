package main

import (
	"log"
	"myweb/action"
	"myweb/feedback"
	"myweb/operate"
	"myweb/show"
	"myweb/statistic"
	"myweb/tools"
	"myweb/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	e.Use(Cors())
	// 获取ip
	tools.GetIP(e)
	// 发送验证码
	tools.Sendmail(e)
	// 用户注册
	user.Register(e)
	// 用户登陆
	user.Login(e)

	// 用户操作
	operate.Operate(e)
	//用户行为
	action.Action(e)

	// 统计
	statistic.Statistic(e)

	// 生成展示的链接
	show.Waterfall(e)
	// 提供api
	feedback.Api(e)
	// 通知
	feedback.Notice(e)
	// 重新授权
	user.Reauth(e)

	e.Run("0.0.0.0:10013")

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
