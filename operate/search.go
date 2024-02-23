package operate

import (
	"database/sql"
	"myweb/tools"

	"strconv"

	"github.com/gin-gonic/gin"
)

func Search(e *gin.Engine) {
	e.POST("/search", search)
}

func search(r *gin.Context) {
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "用户认证失败")
		panic("用户认证失败")
	} else {
		// 连接数据库
		// 连接数据库查看邮箱密码是否正确
		db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")

		if err != nil {
			//没连上的操作
			r.String(500, "连接数据库失败")
			panic(err)
		} else {
			defer db.Close()
			// 参数不能为空
			id := user_info["id"]
			keyword := r.DefaultPostForm("keyword", "")
			if keyword == "" {
				r.String(500, "关键字不能为空")
				panic("关键字不能为空")
			}
			search_sql := "select label,pic_name,pic_url,public from pic_info where pic_name like '%" + keyword + "%'" + "and user_id =" + strconv.FormatFloat(id.(float64), 'f', -1, 32)
			search_res, search_err := db.Query(search_sql)
			if search_err != nil {
				r.String(500, "搜索图片失败")
			} else {
				var search_pics = []map[string]interface{}{}
				for search_res.Next() {
					var pic_info = map[string]interface{}{}
					var label string
					var pic_name string
					var pic_url string
					var public bool
					search_res.Scan(&label, &pic_name, &pic_url, &public)
					pic_info["label"] = label
					pic_info["pic_name"] = pic_name
					pic_info["pic_url"] = pic_url
					pic_info["public"] = public
					search_pics = append(search_pics, pic_info)
				}
				r.JSON(200, search_pics)
			}
		}
	}
}
