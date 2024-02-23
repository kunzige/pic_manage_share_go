package action

import (
	"database/sql"
	"myweb/tools"

	"github.com/gin-gonic/gin"
)

func Agree(r *gin.Context) {
	// 验证身份
	// 验证token是否有效
	token := r.DefaultPostForm("token", "")
	_, parse_err := tools.Jwtparse(token)
	api_id := r.DefaultPostForm("api_id", "")
	if parse_err != nil {
		r.String(403, "用户认证失败")
		panic(parse_err)
	}

	// 连接数据库
	//连接数据库
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	// 更新同意状态
	update_sql := "update apis set agreed = 1 where api_id=?"
	_, update_err := db.Exec(update_sql, api_id)
	if update_err != nil {
		r.String(500, "更新同意状态出错")
		panic(update_err)
	}

}
