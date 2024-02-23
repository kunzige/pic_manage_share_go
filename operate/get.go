package operate

import (
	"context"
	"database/sql"
	"fmt"
	"myweb/tools"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func Get(e *gin.Engine) {
	getGroup := e.Group("/get")
	{
		getGroup.POST("/pic", getpic)
		getGroup.POST("/label", getlabel)
		getGroup.POST("/user", userdata)
		getGroup.POST("/showdata", showdata)
		getGroup.POST("/link", getlink)
		getGroup.POST("/icon", geticon)
		getGroup.POST("/like", getlike)
		getGroup.POST("/info", getuserinfo)
		getGroup.POST("/follow", getfollow)
		getGroup.POST("/qfollow", queryfollow)
		getGroup.POST("/apilist", getapilist)
		getGroup.POST("/notice", getnotice)
		getGroup.POST("/alluser", getalluser)
		getGroup.POST("/store", getstore)
	}
}

// 获取存储得到图片数据
func getpic(r *gin.Context) {
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

		rlabel := r.DefaultPostForm("label", "")
		if rlabel == "" || rlabel == "全部" {
			// 获取全部图片信息
			pic_info, err2 := db.Query("select pic_id,label,pic_name,pic_url,public from pic_info where user_id=" + strconv.FormatFloat(user_info["id"].(float64), 'f', -1, 32))
			if err2 != nil {
				r.String(500, "查询图片数据失败")
				panic(err2)
			} else {

				pic_infos := []map[string]interface{}{}
				//在循环中进行变量的声明，每次都是新的，块级作用域的概念，在外部声明数组循环中append
				for pic_info.Next() { //遍历查询到的每一行结果
					info := make(map[string]interface{})
					var label string
					var pic_name string
					var pic_url string
					var pic_id int
					var public bool

					err = pic_info.Scan(&pic_id, &label, &pic_name, &pic_url, &public) //使变量指向查询到的结果，理解为将查询到的结果进行赋值

					if err != nil {
						r.String(500, "查询图片数据失败")
						panic(err)
					} else {
						//将查询到的数据添加到数组中
						info["pic_id"] = pic_id
						info["label"] = label
						info["pic_name"] = pic_name
						info["pic_url"] = pic_url
						info["public"] = public

						pic_infos = append(pic_infos, info)
					}
				}

				r.JSON(200, pic_infos)
			}
		} else {
			// 根据label进行查询并返回相应的结果
			pic_info, err2 := db.Query("select pic_id,pic_name,pic_url,public from pic_info where user_id=" + strconv.FormatFloat(user_info["id"].(float64), 'f', -1, 32) + " and label='" + rlabel + "'")
			if err2 != nil {
				r.String(500, "查询图片数据失败")
				panic(err2)
			} else {

				pic_infos := []map[string]interface{}{}
				//在循环中进行变量的声明，每次都是新的，块级作用域的概念，在外部声明数组循环中append
				for pic_info.Next() { //遍历查询到的每一行结果
					info := make(map[string]interface{})
					var pic_id int
					var pic_name string
					var pic_url string
					var public bool
					err = pic_info.Scan(&pic_id, &pic_name, &pic_url, &public) //使变量指向查询到的结果，理解为将查询到的结果进行赋值

					if err != nil {
						r.String(500, "查询图片数据失败")
						panic(err)
					} else {
						//将查询到的数据添加到数组中
						info["pic_id"] = pic_id
						info["label"] = rlabel
						info["pic_name"] = pic_name
						info["pic_url"] = pic_url
						info["public"] = public

						pic_infos = append(pic_infos, info)
					}
				}

				r.JSON(200, pic_infos)
			}

		}

	}

}

// 获取所有的标签
func getlabel(r *gin.Context) {
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

			r.JSON(200, label_names)
		}
	}
}

func userdata(r *gin.Context) {
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "认证失败")
		panic("getdata认证失败")
	} else {
		// 连接数据库查看邮箱密码是否正确
		db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
		if err != nil {
			//没连上的操作
			r.String(500, "连接数据库失败")
			panic(err)
		}
		defer db.Close()

		user_data := map[string]interface{}{}
		// 查询所有标签
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
			user_data["labels"] = label_names
		}
		fmt.Printf("user_info: %v\n", user_info)

		// 查询用户数据
		info_res, err2 := db.Query("SELECT email,username,icon,register_time,phone FROM user where id=" + strconv.FormatFloat(user_info["id"].(float64), 'f', -1, 32))
		if err2 != nil {
			panic(err2)
		} else {
			//在循环中进行变量的声明，每次都是新的，块级作用域的概念，在外部声明数组循环中append
			var email string
			var username string
			var icon string
			var register_time string
			var phone string

			for info_res.Next() { //遍历查询到的每一行结果

				err = info_res.Scan(&email, &username, &icon, &register_time, &phone) //使变量指向查询到的结果，理解为将查询到的结果进行赋值

				if err != nil {
					panic(err)
				} else {
					//将查询到的数据添加到数组中

					user_data["username"] = username
					user_data["icon"] = icon
					user_data["register_time"] = register_time
					user_data["email"] = email
					user_data["phone"] = phone
				}
			}

		}

		// 查询图片总数
		sum_pic, sum_err := db.Query("select count(pic_id) from pic_info where user_id = " + strconv.FormatFloat(user_info["id"].(float64), 'f', -1, 32))
		if sum_err != nil {
			r.String(500, "查询图片总数出错")
			panic(sum_err)
		} else {
			var sum int
			for sum_pic.Next() {
				sum_pic.Scan(&sum)
			}
			user_data["sum"] = sum
			r.JSON(200, user_data)
		}

	}

}

func showdata(r *gin.Context) {
	id := r.DefaultPostForm("id", "")
	show_id := r.DefaultPostForm("show_id", "")
	token := r.DefaultPostForm("token", "")
	// 连接数据库，获取数据进行响应
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	opts := redis.Options{
		Addr:     "127.0.0.1:8001", //Addr包含服务器的IP地址与端口 "ip:port"	port在redis.conf文件中进行更改
		Password: "123#Redis",      //有的话就填，没的话就是空，与requirepass参数相对应
		DB:       1,
	}

	rdb := redis.NewClient(&opts)
	ctx := context.Background()
	_, rerr := rdb.Ping(ctx).Result()
	if rerr != nil {
		//没有连接上redis
		panic(rerr)
	}

	// 先获取标签的信息
	query_sql := "select labels from show_info where user_id=? and show_id = ?"
	label_res, query_err := db.Query(query_sql, id, show_id)
	if query_err != nil {
		r.String(500, "查询标签信息失败")
		panic(query_err)
	}
	var labels string
	for label_res.Next() {
		label_res.Scan(&labels)
	}
	label_arr := strtoarr(labels)
	var show_info []map[string]interface{}
	for _, v := range label_arr {
		pic_sql := "select pic_url,pic_name,pic_id from pic_info where label = ? and user_id = ? and public=1"
		pic_res, pic_err := db.Query(pic_sql, strings.Replace(v, `"`, "", -1), id)
		if pic_err != nil {
			r.String(500, "获取指定标签的图片失败")
			panic(pic_err)
		}
		label_info := map[string]interface{}{}
		label_data := []map[string]interface{}{}
		for pic_res.Next() {
			var pic_name string
			var pic_url string
			var pic_id int
			info_item := map[string]interface{}{}
			pic_res.Scan(&pic_url, &pic_name, &pic_id)
			sum_res, _ := rdb.Do(ctx, "get", pic_id).Result()
			info_item["pic_url"] = pic_url
			info_item["pic_name"] = pic_name
			info_item["pic_id"] = pic_id
			info_item["sum"] = sum_res

			// 查询是否喜欢过了
			// 需要先登录才知道喜不喜欢
			user_info, _ := tools.Jwtparse(token)
			if user_info != nil {
				is_liked_sql := "select count(zan_id) from zan where user_id=? and pic_id=?"
				is_liked_res, is_liked_err := db.Query(is_liked_sql, user_info["id"], pic_id)
				if is_liked_err != nil {
					r.String(500, "查询是否喜欢失败")
					panic(is_liked_err)
				}
				for is_liked_res.Next() {
					var is_liked int
					is_liked_res.Scan(&is_liked)
					info_item["is_liked"] = is_liked
				}
			}
			label_data = append(label_data, info_item)
		}
		label_info["label"] = v
		label_info["data"] = label_data
		show_info = append(show_info, label_info)
	}
	r.JSON(200, show_info)
}

func getlink(r *gin.Context) {
	//	验证有没有token
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "认证失败")
		panic("getlink认证失败")
	}

	id := user_info["id"]

	// 连接数据库，获取数据进行响应
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	query_sql := "select user_id,show_id,labels from show_info where user_id=?"
	query_res, query_err := db.Query(query_sql, id)
	if query_err != nil {
		r.String(500, "查询生成的链接失败")
		panic(query_err)
	}
	data := []map[string]interface{}{}
	for query_res.Next() {
		var user_id int
		var show_id int
		var labels string
		query_res.Scan(&user_id, &show_id, &labels)
		data = append(data, map[string]interface{}{
			"label":   labels,
			"user_id": user_id,
			"show_id": show_id,
		})
	}
	r.JSON(200, data)
}

func geticon(r *gin.Context) {
	id := r.DefaultPostForm("id", "")
	geticon_sql := "select icon from user where id = ?"
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()
	icon_res, icon_err := db.Query(geticon_sql, id)
	if icon_err != nil {
		r.String(500, "获取用户头像失败")
		panic(icon_err)
	}
	for icon_res.Next() {
		var icon string
		icon_res.Scan(&icon)
		r.String(200, icon)
	}
}

func getlike(r *gin.Context) {
	// 验证token并获取id
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "认证失败")
		panic("getlike认证失败")
	}

	id := user_info["id"]

	// 链接数据库
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	opts := redis.Options{
		Addr:     "127.0.0.1:8001", //Addr包含服务器的IP地址与端口 "ip:port"	port在redis.conf文件中进行更改
		Password: "123#Redis",      //有的话就填，没的话就是空，与requirepass参数相对应
		DB:       1,
	}

	rdb := redis.NewClient(&opts)
	ctx := context.Background()
	_, rerr := rdb.Ping(ctx).Result()
	if rerr != nil {
		//没有连接上redis
		panic(rerr)
	}
	zaninfo := []map[string]interface{}{}

	getzan_sql := "select zan_id from zan where user_id = ?"
	getzan_res, getzan_err := db.Query(getzan_sql, id)
	if getzan_err != nil {
		r.String(500, "获取用户点赞的图片id失败")
		panic(getzan_err)
	}
	for getzan_res.Next() {
		var zan_id int
		getzan_res.Scan(&zan_id)

		getdata1_sql := "select pic_id from zan where zan_id=?"
		getdata1_res, _ := db.Query(getdata1_sql, zan_id)

		item := map[string]interface{}{}
		for getdata1_res.Next() {
			var pic_id int
			var user_id int
			getdata1_res.Scan(&pic_id)
			// 获取被点赞的用户的ID
			getuser_id := "select user_id from pic_info where pic_id=?"
			getuser_id_res, _ := db.Query(getuser_id, pic_id)
			for getuser_id_res.Next() {
				getuser_id_res.Scan(&user_id)
			}
			getdata2_sql := "select pic_url,icon from pic_info,user where pic_id=? and id=?"
			getdata2_res, _ := db.Query(getdata2_sql, pic_id, user_id)
			for getdata2_res.Next() {
				var pic_url string
				var icon string
				getdata2_res.Scan(&pic_url, &icon)
				// 读取redis缓存
				ret, _ := rdb.Do(ctx, "get", pic_id).Result()
				item["pic_url"] = pic_url
				item["icon"] = icon
				item["like_sum"] = ret
				item["pic_id"] = pic_id
				item["user_id"] = user_id

				isliked_sql := "select count(zan_id) from zan where user_id=? and pic_id=?"
				isliked_res, _ := db.Query(isliked_sql, id, pic_id)
				for isliked_res.Next() {
					var is_liked int
					isliked_res.Scan(&is_liked)
					item["is_liked"] = is_liked
				}
			}
			zaninfo = append(zaninfo, item)
		}
	}
	r.JSON(200, zaninfo)

}

func getuserinfo(r *gin.Context) {
	id := r.DefaultPostForm("id", "")
	userinfo_sql := "select username,email,icon,register_time,phone from user where id = ?"
	// 链接数据库
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()
	var username string
	var email string
	var icon string
	var register_time string
	var phone int
	var sum int
	labels := []string{}
	// 用户基本信息
	userinfo_res, userinfo_err := db.Query(userinfo_sql, id)
	if userinfo_err != nil {
		r.String(500, "查询该与用户信息失败")
		panic(userinfo_err)
	}
	for userinfo_res.Next() {
		userinfo_res.Scan(&username, &email, &icon, &register_time, &phone)
	}

	// 获取存储图片总数
	getsum_sql := "select count(pic_id) from pic_info where user_id=?"

	getsum_res, getsum_err := db.Query(getsum_sql, id)
	if getsum_err != nil {
		r.String(500, "获取用户存储总数失败")
		panic(getsum_err)
	}

	for getsum_res.Next() {
		getsum_res.Scan(&sum)
	}

	// 获取存储的所有标签
	getlabels_res, _ := db.Query("select distinct label from pic_info where user_id=?", id)
	for getlabels_res.Next() {
		var label string
		getlabels_res.Scan(&label)
		labels = append(labels, label)
	}

	r.JSON(200, map[string]interface{}{
		"username":      username,
		"email":         email,
		"icon":          icon,
		"register_time": register_time,
		"phone":         phone,
		"sum":           sum,
		"labels":        labels,
	})

}

func getfollow(r *gin.Context) {
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "认证失败")
		panic("getfollow认证失败")
	}

	id := user_info["id"]

	// 连接数据库，获取数据进行响应
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	// 查询关注了哪些用户
	getfollow_sql := "select target_id from follow where user_id = ?"
	getfollow_res, _ := db.Query(getfollow_sql, id)
	target_ids := []int{}
	follow_info := []map[string]interface{}{}
	for getfollow_res.Next() {
		var target_id int
		getfollow_res.Scan(&target_id)
		target_ids = append(target_ids, target_id)
	}
	for _, v := range target_ids {
		userinfo_res, _ := db.Query("select username,email,icon,id from user where id = ?", v)
		userinfo := map[string]interface{}{}
		for userinfo_res.Next() {
			var username string
			var email string
			var icon string
			var id int
			userinfo_res.Scan(&username, &email, &icon, &id)
			userinfo["username"] = username
			userinfo["email"] = email
			userinfo["icon"] = icon
			userinfo["id"] = id
		}
		follow_info = append(follow_info, userinfo)
	}
	r.JSON(200, follow_info)
}

// 查询是否进行了关注
func queryfollow(r *gin.Context) {

	// 连接数据库，获取数据进行响应
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	token := r.DefaultPostForm("token", "")
	user_info, _ := tools.Jwtparse(token)
	target := r.DefaultPostForm("target", "")

	if user_info != nil {
		id := user_info["id"]
		query_sql := "select count(follow_id) from follow where user_id=? and target_id=?"
		query_res, query_err := db.Query(query_sql, id, target)
		if query_err != nil {
			r.String(500, "查询是否关注出错")
			panic(query_err)
		}
		for query_res.Next() {
			var is_follow int
			query_res.Scan(&is_follow)
			r.JSON(200, gin.H{
				"is_follow": is_follow,
			})
		}
		return
	}
	query_sql := "select count(follow_id) from follow where user_id=? and target_id=?"
	query_res, query_err := db.Query(query_sql, -1, target)
	if query_err != nil {
		r.String(500, "查询是否关注出错")
		panic(query_err)
	}
	for query_res.Next() {
		var is_follow int
		query_res.Scan(&is_follow)
		r.JSON(200, gin.H{
			"is_follow": is_follow,
		})
	}

}

func getapilist(r *gin.Context) {

	// 连接数据库，获取数据进行响应
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	query_sql := "select api_id,api,argument,agreed from apis"
	query_res, _ := db.Query(query_sql)
	list := []map[string]interface{}{}

	for query_res.Next() {
		item := map[string]interface{}{}
		var api_id int
		var agreed bool
		var api string
		var argument string
		query_res.Scan(&api_id, &api, &argument, &agreed)
		item["api_id"] = api_id
		item["api"] = api
		item["argument"] = argument
		item["agreed"] = agreed
		list = append(list, item)
	}
	r.JSON(200, list)

}

func getnotice(r *gin.Context) {
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "认证失败")
		panic("getnotice认证失败")
	}

	id := user_info["id"]

	// 连接数据库，获取数据进行响应
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()
	notice_info := map[string]interface{}{}
	getsum_sql := "select count(notice_id) from notice where is_readed = 0 and user_id=?"
	getsum_res, _ := db.Query(getsum_sql, id)
	for getsum_res.Next() {
		var sum int
		getsum_res.Scan(&sum)
		notice_info["sum"] = sum
	}
	query_sql := "select title,notice_id,content,is_readed from notice where user_id=?"
	query_res, query_err := db.Query(query_sql, id)
	if query_err != nil {
		r.String(500, "查询通知失败")
		panic(query_err)
	}
	item_arr := []map[string]interface{}{}
	for query_res.Next() {
		item := map[string]interface{}{}
		var title string
		var notice_id int
		var content string
		var is_readed bool
		query_res.Scan(&title, &notice_id, &content, &is_readed)
		item["notice_id"] = notice_id
		item["content"] = content
		item["is_readed"] = is_readed
		item["title"] = title
		item_arr = append(item_arr, item)
	}
	notice_info["data"] = item_arr
	r.JSON(200, notice_info)
}

func getalluser(r *gin.Context) {
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "认证失败")
		panic("getalluser认证失败")
	}

	id := user_info["id"]
	if int(id.(float64)) != 1 {
		r.String(403, "权限不够")
		panic("权限不能，不能获取用户数据")
	}

	// 连接数据库，获取数据进行响应
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	// 查询关注了哪些用户
	getid_sql := "select id from user where id != 1"
	getfollow_res, _ := db.Query(getid_sql)
	target_ids := []int{}
	users_info := []map[string]interface{}{}
	for getfollow_res.Next() {
		var target_id int
		getfollow_res.Scan(&target_id)
		target_ids = append(target_ids, target_id)
	}
	for _, v := range target_ids {
		userinfo_res, _ := db.Query("select username,email,icon,id from user where id = ?", v)
		userinfo := map[string]interface{}{}
		for userinfo_res.Next() {
			var username string
			var email string
			var icon string
			var id int
			userinfo_res.Scan(&username, &email, &icon, &id)
			userinfo["username"] = username
			userinfo["email"] = email
			userinfo["icon"] = icon
			userinfo["id"] = id
		}
		users_info = append(users_info, userinfo)
	}
	r.JSON(200, users_info)
}

func getstore(r *gin.Context) {
	// 验证管理员身份
	token := r.DefaultPostForm("token", "")
	user_info, parse_err := tools.Jwtparse(token)
	if parse_err != nil {
		r.String(403, "认证失败")
		panic("getstore认证失败")
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

	// 获取id
	id := r.DefaultPostForm("id", "")
	store_info := map[string]interface{}{}
	// 查询存储总数与Max
	getsum_sql := "select count(pic_id) from pic_info  where user_id=?"
	getsum_res, _ := db.Query(getsum_sql, id)
	for getsum_res.Next() {
		var sum int
		getsum_res.Scan(&sum)
		store_info["sum"] = sum
	}
	getmax_sql := "select max from store where user_id=?"
	getmax_res, _ := db.Query(getmax_sql, id)
	for getmax_res.Next() {
		var max int
		getmax_res.Scan(&max)
		store_info["max"] = max
	}

	r.JSON(200, store_info)
}

func strtoarr(str string) []string {
	str = strings.Replace(str, "[", "", -1)
	str = strings.Replace(str, "]", "", -1)
	labels := strings.Split(str, ",")
	return labels
}
