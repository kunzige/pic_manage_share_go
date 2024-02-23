package user

import (
	"myweb/tools"
	"time"

	"github.com/gin-gonic/gin"
)

func Reauth(e *gin.Engine) {
	e.POST("/user/reauth", reauth)
}

func reauth(r *gin.Context) {
	refre_token := r.DefaultPostForm("refresh_token", "")
	if refre_token == "" {
		r.String(403, "无效")
		panic("无法重新授权")
	}
	r.JSON(200, map[string]interface{}{
		"token": tools.Regenerate(int64(time.Now().Unix()+60*60*3), refre_token),
	})
}
