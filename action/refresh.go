package action

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func Refresh(r *gin.Context) {
	//连接数据库
	db, err := sql.Open("mysql", "root:123#Mysql@tcp(127.0.0.1:8000)/mysite")
	if err != nil {
		//没连上的操作
		r.String(500, "连接数据库失败")
		panic(err)
	}
	defer db.Close()

	// 更新redis缓存
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
	} else {
		// 连接成功redis
		// mysql中找到所有的picid，然后在zan表中sum统计总的点赞数，缓存到redis
		getpic_id := "select pic_id from pic_info;"
		getid_res, getid_err := db.Query(getpic_id)
		if getid_err != nil {
			r.String(500, "获取所有的图片id失败")
			panic(getid_err)
		}
		count_info := map[int]interface{}{}
		for getid_res.Next() {
			var pic_id int
			getid_res.Scan(&pic_id)
			count_sql := "select count(zan_id) from zan where pic_id = ?"
			count_res, count_err := db.Query(count_sql, pic_id)
			if count_err != nil {
				r.String(500, "统计赞的数量出错")
				panic(count_err)
			}
			var sum int
			for count_res.Next() {
				count_res.Scan(&sum)
				count_info[pic_id] = sum
			}
			rdb.Do(ctx, "set", pic_id, sum)
		}
	}

}
