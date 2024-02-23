package operate

import (
	"database/sql"
	"myweb/tools"

	"github.com/gin-gonic/gin"
)

func Delete(e *gin.Engine) {
	deleteGroup := e.Group("/delete")
	{
		deleteGroup.DELETE("/pic", pic)
		deleteGroup.DELETE("/share", share)
		deleteGroup.DELETE("/user", deleteuser)
	}

}

func pic(r *gin.Context) {
	// 检验token是否有效
	pic_id := r.DefaultPostForm("pic_id", "")
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "用户认证失败")
		panic("用户认证失败")
	} else {
		// 获取pic_id与id
		id := user_info["id"]
		// 连接数据库
		// 连接数据库查看邮箱密码是否正确
		db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
		if err != nil {
			//没连上的操作
			r.String(500, "连接数据库失败")
			panic(err)
		}
		defer db.Close()

		insert_sql := "delete from pic_info where user_id = ? and pic_id = ?"
		res, err := db.Exec(insert_sql, id, pic_id)
		_ = res //没什么用
		if err != nil {
			r.String(500, "删除数据失败")
			panic("删除数据失败")
		} else {
			r.JSON(200, gin.H{
				"status":  "ok",
				"message": "删除成功",
			})
		}
	}
}

func share(r *gin.Context) {
	// 检验token是否有效
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "用户认证失败")
		panic("用户认证失败")
	} else {
		// 获取pic_id与id
		id := user_info["id"]
		show_id := r.DefaultPostForm("show_id", "")
		if show_id == "" {
			r.String(403, "参数不正确")
			panic("删除分享参数不完整")
		}
		// 连接数据库
		// 连接数据库查看邮箱密码是否正确
		db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
		if err != nil {
			//没连上的操作
			r.String(500, "连接数据库失败")
			panic(err)
		}
		defer db.Close()

		insert_sql := "delete from show_info where user_id = ? and show_id = ?"
		res, err := db.Exec(insert_sql, id, show_id)
		_ = res //没什么用
		if err != nil {
			r.String(500, "删除分享数据失败")
			panic(err)
		} else {
			r.JSON(200, gin.H{
				"status":  "ok",
				"message": "删除成功",
			})
		}
	}
}

func deleteuser(r *gin.Context) {
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

	user_id := r.DefaultPostForm("user_id", "")
	delete_sql := "delete from user where id = ?"
	_, delete_err := db.Exec(delete_sql, user_id)
	if delete_err != nil {
		r.String(500, "删除用户发生了一些错误")
		panic(delete_err)
	}
	r.String(200, "成功删除用户。")

}
