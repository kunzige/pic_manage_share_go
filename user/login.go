package user

import (
	"database/sql"
	"myweb/tools"
	"time"

	"github.com/gin-gonic/gin"
)

func Login(e *gin.Engine) {
	e.POST("/login", login)
}

func login(r *gin.Context) {
	// 获取邮箱和密码
	email := r.DefaultPostForm("email", "")
	password := r.DefaultPostForm("password", "")

	// 验证是否输入用户名密码
	if email == "" || password == "" {
		r.String(403, "用户信息不完整")
		panic("用户信息不完整")
	}

	// 连接数据库查看邮箱密码是否正确
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "数据库没开启")

		panic(err)
	}
	defer db.Close()

	// 查询邮箱是否已经存在
	email_res, err2 := db.Query("SELECT password,id,username,icon,register_time FROM user where email=" + "'" + email + "'")
	if err2 != nil {
		r.String(500, "连接数据库失败")
		panic(err2)
	} else {
		//在循环中进行变量的声明，每次都是新的，块级作用域的概念，在外部声明数组循环中append
		var u_password string
		var id int
		var username string
		var icon string
		var register_time string
		var user_info = map[string]interface{}{}

		for email_res.Next() { //遍历查询到的每一行结果

			err = email_res.Scan(&u_password, &id, &username, &icon, &register_time) //使变量指向查询到的结果，理解为将查询到的结果进行赋值

			if err != nil {
				panic(err)
			} else {
				//将查询到的数据添加到数组中

				user_info["username"] = username
				user_info["icon"] = icon
				user_info["register_time"] = register_time
				user_info["id"] = id
				user_info["email"] = email
			}
		}

		if u_password == password {
			r.JSON(200, gin.H{
				"message":       "登陆成功",
				"token":         tools.Generatejwt(int64(time.Now().Unix()+60*60*3), user_info),
				"refresh_token": tools.Generatejwt(int64(time.Now().Unix()+60*60*24*3), user_info),
			})
		} else {
			r.String(403, "密码错误")
		}

	}
}
