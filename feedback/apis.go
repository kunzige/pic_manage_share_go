package feedback

import (
	"database/sql"
	"myweb/tools"

	"github.com/gin-gonic/gin"
)

func Api(e *gin.Engine) {
	e.POST("/feedback/api", api)
}

func api(r *gin.Context) {
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "用户认证失败")
		panic(parse_err)
	}

	api := r.DefaultPostForm("api", "")
	argument := r.DefaultPostForm("argument", "")

	// 连接数据库
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	// 向数据库中插入值
	insert_sql := "insert into apis(api_id,api,argument,user_id) values(?,?,?,?)"
	_, insert_err := db.Exec(insert_sql, 0, api, argument, user_info["id"])
	if insert_err != nil {
		r.String(500, "写入api失败")
		panic(insert_err)
	}
	r.String(200, "写入api成功!")
}
