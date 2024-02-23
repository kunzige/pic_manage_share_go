package show

import (
	"database/sql"
	"myweb/tools"
	"strings"

	"github.com/gin-gonic/gin"
)

func Waterfall(e *gin.Engine) {
	e.POST("/show/generate", waterfall)
}

func waterfall(r *gin.Context) {
	// 验证token并获取Id
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
		// 获取标签列表
		labels := r.DefaultPostForm("labels", "")

		if labels == "" {
			r.String(500, "标签不能为空")
			panic("标签不能为空")
		}
		labels_arr := strtoarr(labels)

		// 查询不能重复
		query_sql := "select labels from show_info where user_id=? and labels like" + "'%" + labels_arr[0] + "%'"
		query_res, query_err := db.Query(query_sql, id)
		if query_err != nil {
			r.String(500, "查询标签信息失败！")
			panic(query_err)
		}
		for query_res.Next() {
			var elabels_str string
			query_res.Scan(&elabels_str)
			elabels_arr := strtoarr(elabels_str)

			count := 0
			for _, v1 := range labels_arr {
				for _, v2 := range elabels_arr {
					if v1 == v2 {
						count += 1
					}
				}
			}
			if len(labels_arr) == len(elabels_arr) {
				if count == len(elabels_arr) {
					r.String(500, "您已生成该标签")
					panic("您已申生成该标签")
				}
			}

		}

		insert_sql := "insert into show_info (show_id,user_id,labels) values (0,?,?)"
		_, insert_err := db.Exec(insert_sql, id, labels)
		if insert_err != nil {
			r.String(500, "写入标签信息失败！")
			panic("写入标签信息失败！")
		}
		show_id__sql := "select show_id from show_info where user_id=? and labels = ?"
		show_id_res, show_id_err := db.Query(show_id__sql, id, labels)
		if show_id_err != nil {
			r.String(500, "获取请求链接失败")
			panic(show_id_err)
		}
		var show_id int
		for show_id_res.Next() {
			show_id_res.Scan(&show_id)
		}
		r.JSON(200, map[string]interface{}{
			"id":      id,
			"show_id": show_id,
		})
	}

}

func strtoarr(str string) []string {
	str = strings.Replace(str, "[", "", -1)
	str = strings.Replace(str, "]", "", -1)
	labels := strings.Split(str, ",")
	return labels
}
