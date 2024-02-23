package operate

import "github.com/gin-gonic/gin"

func Operate(e *gin.Engine) {
	// 添加图片的操作
	Add(e)
	Get(e)
	Delete(e)
	Modify(e)
	Search(e)
}
