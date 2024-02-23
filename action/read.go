package action

import (
	"database/sql"
	"myweb/tools"

	"github.com/gin-gonic/gin"
)

func Read(r *gin.Context) {
	// 获取token id notice_id
	// 验证token是否有效
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	notice_id := r.DefaultPostForm("notice_id", "")
	if parse_err != nil {
		r.String(403, "用户认证失败")
		panic(parse_err)
	}
	//连接数据库
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	// 更新状态
	update_sql := "update notice set is_readed = 1 where user_id = ? and notice_id= ?"
	_, update_err := db.Exec(update_sql, user_info["id"], notice_id)
	if update_err != nil {
		r.String(500, "更新通知已读状态失败")
		panic(update_err)
	}
	r.String(200, "更新通知已读状态成功")

}
