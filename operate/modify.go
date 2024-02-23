package operate

import (
	"database/sql"
	"myweb/tools"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Modify(e *gin.Engine) {
	modify_group := e.Group("/modify")
	{
		modify_group.POST("/icon", modifyicon)
		modify_group.POST("/info", modifyinfo)
		modify_group.POST("/pic_info", modify_picinfo)
		modify_group.POST("/store", modify_max)
	}
}

func modifyicon(r *gin.Context) {
	// 验证token
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

		newicon := r.DefaultPostForm("newicon", "")
		if newicon == "" {
			r.String(500, "请传递用户头像")
			panic("请传递用户头像")
		}
		modify_sql := "update user set icon = ? where id = ?"
		mod_res, mod_err := db.Exec(modify_sql, newicon, id)

		if mod_err != nil {
			r.String(500, "修改失败")
			panic(mod_err)
		} else {
			_ = mod_res
			user_info["icon"] = newicon
			r.JSON(200, map[string]interface{}{
				"message": "修改成功",
				"newicon": newicon,
			})
		}

	}
}

func modifyinfo(r *gin.Context) {
	// 验证token
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

		newname := r.DefaultPostForm("newname", "")
		newphone := r.DefaultPostForm("newphone", "")
		if newname == "" || newphone == "" {
			r.String(500, "需要完整用户信息")
			panic("需要完整用户信息")
		}
		modify_sql := "update user set username = ? , phone = ? where id = ?"
		mod_res, mod_err := db.Exec(modify_sql, newname, newphone, id)

		if mod_err != nil {
			r.String(500, "修改用户信息失败")
			panic(mod_err)
		} else {
			_ = mod_res
			r.JSON(200, map[string]interface{}{
				"message":  "修改用户信息成功",
				"newname":  newname,
				"newphone": newphone,
			})
		}
	}
}

func modify_picinfo(r *gin.Context) {
	// 验证token是否有效
	token := r.DefaultPostForm("token", "")
	_, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "用户认证失败")
		panic("用户认证失败")
	}

	label := r.DefaultPostForm("label", "")
	pic_name := r.DefaultPostForm("name", "")
	public := r.DefaultPostForm("public", "")
	pic_id := r.DefaultPostForm("id", "")
	bpublic := 0
	if public == "true" {
		bpublic = 1
	}
	if label == "" || pic_name == "" || public == "" {
		r.String(403, "参数不完整")
		panic("修改图片属性参数不完整")
	}

	// 连接数据库
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	//修改相关字段
	modify_sql := "update pic_info set label = ?,pic_name = ?,public = " + strconv.Itoa(bpublic) + " where pic_id = ?"
	_, modify_err := db.Exec(modify_sql, label, pic_name, pic_id)
	if modify_err != nil {
		r.String(500, "修改图片属性失败")
		panic(modify_err)
	}
	r.String(200, "修改图片属性成功")
}

func modify_max(r *gin.Context) {
	// 验证超级管理员
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
	// 修改max
	newmax := r.DefaultPostForm("newmax", "")
	user_id := r.DefaultPostForm("user_id", "")
	modify_sql := "update store set max = ? where user_id = ?"
	_, modify_err := db.Exec(modify_sql, newmax, user_id)
	if modify_err != nil {
		r.String(500, "修改用户可存储最多数失败")
		panic(modify_err)
	}
	r.String(200, "修改用户最多存储成功")

}
