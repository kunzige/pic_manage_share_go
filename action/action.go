package action

import "github.com/gin-gonic/gin"

func Action(e *gin.Engine) {
	e.POST("/action/zan", Zan)
	e.DELETE("/action/cancelzan", CancelZan)
	e.POST("/action/refresh", Refresh)
	e.POST("/action/follow", Follow)
	e.DELETE("action/cancelfollow", CancelFollow)
	e.POST("action/agree", Agree)
	e.POST("action/read", Read)
}
