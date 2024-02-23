package tools

import (
	"context"
	"math/rand"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gophish/gomail"
)

func Sendmail(e *gin.Engine) {
	e.POST("/sendcode", sendmail)
}

func sendmail(r *gin.Context) {
	// 生成验证码
	code := randomCode()

	m := gomail.NewMessage()
	//发送人
	m.SetHeader("From", "kunzige666@qq.com")
	//接收人
	to_email := r.DefaultPostForm("mail", "")
	re, _ := regexp.Compile("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$")
	match := re.FindStringIndex(to_email)
	if match == nil {
		r.String(500, "发送失败")
		panic("不能获取有效邮箱！")
	}

	m.SetHeader("To", to_email)
	//主题
	m.SetHeader("Subject", "注册验证码")
	//内容
	m.SetBody("text/html", "<h1>"+"您好，您此次注册的验证码是:"+"<br>"+code+"</h1>"+"<h2>15分钟内有效</h2>")
	//附件
	//m.Attach("./myIpPic.png")

	//拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer("smtp.qq.com", 587, "kunzige666@qq.com", "ftzukinegfohegbd")

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		r.String(500, "发送失败")
		panic(err)
	} else {

		// 发送成功后存放到redis数据库
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
			rdb.Do(ctx, "setex", to_email, 60*15, code).Result()
		}

		r.JSON(200, gin.H{
			"message": "发送成功",
		})
	}

}

// 生成验证码
func randomCode() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	random_code := ""
	rand_int := 0
	rand.Seed(time.Now().Unix())
	for i := 0; i < 4; i++ {
		rand_int = rand.Intn(62)
		random_code += str[rand_int : rand_int+1]
	}
	return random_code
}
