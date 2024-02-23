package operate

import (
	"database/sql"
	"myweb/tools"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Add(e *gin.Engine) {
	e.POST("/add", add)
}

func add(r *gin.Context) {

	// 验证用户是否有效
	token := r.DefaultPostForm("token", "")
	label := r.DefaultPostForm("label", "")
	pic_name := r.DefaultPostForm("pic_name", "")
	pic_url := r.DefaultPostForm("pic_url", "")
	user_info, parse_err := tools.Jwtparse(token)
	public := r.DefaultPostForm("public", "")
	if parse_err != nil {
		r.String(403, "用户认证失败")
		panic(parse_err)
	} else {
		// 取出id查询，到时候当做插入数据的标识
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
		// 查询该用户是否到达临界值
		getsum_sql := "select count(pic_id) from pic_info where user_id = ?"
		getsum_res, _ := db.Query(getsum_sql, user_info["id"])
		sum := 0
		for getsum_res.Next() {
			var sum_pic int
			getsum_res.Scan(&sum_pic)
			sum = sum_pic

		}

		getmax_sql := "select max from store where user_id = ?"
		getmax_res, _ := db.Query(getmax_sql, user_info["id"])
		for getmax_res.Next() {
			var max int
			getmax_res.Scan(&max)
			if sum >= max {
				r.String(403, "您已达上传上限")
				return
			}
		}

		pub, _ := strconv.Atoi(public)
		insert_sql := "insert into pic_info(user_id,pic_id,label,pic_name,pic_url,public) values(?,0,?,?,?,?);"
		res, err := db.Exec(insert_sql, id, label, pic_name, pic_url, pub)
		if err != nil {
			r.String(500, "添加失败")
			panic(err)
		} else {
			_, _ = res.LastInsertId()
			r.JSON(200, gin.H{
				"message": "添加成功",
			})
		}

	}

}
