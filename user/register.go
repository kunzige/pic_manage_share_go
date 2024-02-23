package user

import (
	"context"
	"database/sql"
	"fmt"
	"myweb/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

func Register(e *gin.Engine) {
	e.POST("/register", register)
}

func register(r *gin.Context) {
	// 获取明文数据，进行MD5加密并添加到数据库
	password := r.DefaultPostForm("password", "")
	username := r.DefaultPostForm("username", "")
	reg_email := r.DefaultPostForm("email", "")
	register_time := r.DefaultPostForm("register_time", "")
	code := r.DefaultPostForm("code", "")

	// MD5处理后的密码
	e_password := tools.Md5(password)

	// 验证验证码是否有效
	opts := redis.Options{
		Addr:     "127.0.0.1:8001", //Addr包含服务器的IP地址与端口 "ip:port"	port在redis.conf文件中进行更改
		Password: "123#Redis",      //有的话就填，没的话就是空，与requirepass参数相对应
		DB:       0,
	}

	rdb := redis.NewClient(&opts)
	ctx := context.Background()
	_, rerr := rdb.Ping(ctx).Result()
	if rerr != nil {
		//没有连接上redis
		panic(rerr)
	} else {
		//进行增删改查，不用else也可以，因为有panic
		ret, _ := rdb.Do(ctx, "get", reg_email).Result()
		ecode := fmt.Sprintf("%v", ret)
		if code != ecode {
			r.String(500, "验证码错误")
			panic("错误的验证码")
		}
	}

	// 连接数据库
	// 先初始化，再看有没有错
	//先初始化
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		panic(err)
	}
	defer db.Close()

	// 查询邮箱是否已经存在
	email_res, err2 := db.Query("SELECT email FROM user")
	if err2 != nil {
		panic(err2)
	} else {
		//在循环中进行变量的声明，每次都是新的，块级作用域的概念，在外部声明数组循环中append
		for email_res.Next() { //遍历查询到的每一行结果

			var email string

			err = email_res.Scan(&email) //使变量指向查询到的结果，理解为将查询到的结果进行赋值

			if err != nil {
				panic(err)
			} else {
				//将查询到的数据添加到数组中
				if reg_email == email {
					r.String(503, "邮箱已经存在")
					panic("邮箱已存在")
				}
			}
		}
	}

	//先校验验证码是否正确，正确之后在插入用户表
	// 与插入有关的代码
	insert_sql := `insert into user(id,username,email,register_time,password) values (?,?,?,?,?);`
	ret, err := db.Exec(insert_sql, 0, username, reg_email, register_time, e_password)
	if err != nil {
		panic(err)
	} else {
		// store表
		// 获取最新的user_id
		getlatest_id := "select @@identity"
		latestid_res, _ := db.Query(getlatest_id)
		latest_id := 0
		for latestid_res.Next() {
			var user_id int
			latestid_res.Scan(&user_id)
			latest_id = user_id
		}
		_, store_err := db.Exec("insert into store (store_id,user_id,max) values(0,?,5000)", latest_id)
		if store_err != nil {
			r.String(500, "插入存储信息表失败")
			panic(store_err)
		}
		_, _ = ret.LastInsertId()
		r.JSON(200, gin.H{
			"message": "注册成功",
		})
	}

}
