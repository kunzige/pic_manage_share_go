package feedback

import (
	"database/sql"
	"myweb/tools"

	"github.com/gin-gonic/gin"
)

func Notice(e *gin.Engine) {
	e.POST("feedback/notice", notice)
	e.POST("feedback/infonotice", infonotice)
}

func notice(r *gin.Context) {
	// 验证token是否有效
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	content := r.DefaultPostForm("content", "")
	title := r.DefaultPostForm("title", "")
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

	// 写入通知内容
	insert_sql := "insert into notice(notice_id,user_id,content,is_readed,title) values(0,?,?,0,?)"
	_, insert_err := db.Exec(insert_sql, user_info["id"], content, title)
	if insert_err != nil {
		r.String(500, "写入通知信息失败")
		panic(insert_err)
	}
	r.String(200, "写入通知信息成功")
}

func infonotice(r *gin.Context) {

	// 验证管理员身份
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "认证失败")
		panic("getdata认证失败")
	}

	if int(user_info["id"].(float64)) != 1 {
		r.String(403, "权限不够")
		panic("权限不能，不能获取用户存储信息")
	}

	// 连接数据库，获取数据进行响应
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	// 写入通知内容
	content := r.DefaultPostForm("content", "")
	title := r.DefaultPostForm("title", "")
	user_id := r.DefaultPostForm("user_id", "")
	insert_sql := "insert into notice(notice_id,user_id,content,is_readed,title) values(0,?,?,0,?)"
	_, insert_err := db.Exec(insert_sql, user_id, content, title)
	if insert_err != nil {
		r.String(500, "写入通知信息失败")
		panic(insert_err)
	}
	r.String(200, "写入通知信息成功")
}
