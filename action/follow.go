package action

import (
	"database/sql"
	"myweb/tools"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Follow(r *gin.Context) {
	// 验证token是否有效
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	target := r.DefaultPostForm("target", "")
	idstr := strconv.FormatFloat(user_info["id"].(float64), 'f', -1, 32)
	if idstr == target {
		r.String(401, "自己不能关注自己")
		return
	}

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
	insert_sql := "insert into follow(follow_id,user_id,target_id) values(0,?,?)"
	_, insert_err := db.Exec(insert_sql, user_info["id"], target)
	if insert_err != nil {
		r.String(500, "关注失败")
		panic(insert_err)
	}
}

func CancelFollow(r *gin.Context) {
	// 验证token是否有效
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	target := r.DefaultPostForm("target", "")
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
	insert_sql := "delete from follow where user_id=? and target_id=?"
	_, insert_err := db.Exec(insert_sql, user_info["id"], target)
	if insert_err != nil {
		r.String(500, "取消关注失败")
		panic(insert_err)
	}
}
