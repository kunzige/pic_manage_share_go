package statistic

import (
	"database/sql"
	"myweb/tools"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Labelcount(e *gin.Engine) {
	e.POST("/statistic/labelcount", labelcount)
}

func labelcount(r *gin.Context) {
	// 解析token拿到用户id然后返回数据
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "用户认证失败")
		panic("token解析失败")
	} else {
		// 连接数据库
		// 连接数据库查看邮箱密码是否正确
		db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
		if err != nil {
			//没连上的操作
			r.String(500, "连接数据库失败")
			panic(err)
		}
		defer db.Close()

		labels, err2 := db.Query("select distinct label from pic_info where user_id=" + strconv.FormatFloat(user_info["id"].(float64), 'f', -1, 32))
		if err2 != nil {
			r.String(500, "查询图片数据失败")
			panic(err2)
		} else {
			label_names := make([]string, 0)
			for labels.Next() {
				var label string
				labels.Scan(&label)
				label_names = append(label_names, label)
			}
			label_pic_count := map[string]interface{}{}
			for _, v := range label_names {
				countsql, count_err := db.Query("select count(pic_id) from pic_info where label= " + "'" + v + "'" + "and user_id = " + strconv.FormatFloat(user_info["id"].(float64), 'f', -1, 32))
				if count_err != nil {
					r.String(500, "统计标签数量失败")
					panic("统计标签数量失败")
				}
				for countsql.Next() {
					var count int
					countsql.Scan(&count)
					label_pic_count[v] = count
				}

			}
			// 查询图片总数
			sum_pic, sum_err := db.Query("select count(pic_id) from pic_info where user_id = " + strconv.FormatFloat(user_info["id"].(float64), 'f', -1, 32))
			if sum_err != nil {
				r.String(500, "查询图片总数出错")
				panic(sum_err)
			} else {
				var pic_sum int
				for sum_pic.Next() {
					sum_pic.Scan(&pic_sum)
				}
				label_pic_count["pic_sum"] = pic_sum

			}

			// 查询图片总数
			sum_label, label_err := db.Query("select count(distinct label) from pic_info where user_id =" + strconv.FormatFloat(user_info["id"].(float64), 'f', -1, 32))
			if label_err != nil {
				r.String(500, "查询图片总数出错")
				panic(sum_err)
			} else {
				var label_sum int
				for sum_label.Next() {
					sum_label.Scan(&label_sum)
				}
				label_pic_count["label_sum"] = label_sum

			}

			r.JSON(200, label_pic_count)

		}
	}
}
