package user

import (
	"encoding/json"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var tokenKey = []byte("ysite.kunzige")

// 当做中间件就行
func Authenticate(r *gin.Context) {
	token_info := r.DefaultPostForm("token", "")
	fmt.Println(token_info)
	token_map := map[string]interface{}{}
	if err := json.Unmarshal([]byte(token_info), &token_map); err != nil {
		r.String(403, "forbidden")
		r.Abort()
		return
	}
	access_token := token_map["access_token"].(string)
	if token_info != "" {
		// access_token存在再继续，不存在直接结束函数
		token, err := jwt.Parse(access_token, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			// 通过keyfunc返回token和错误
			if !ok {
				return nil, fmt.Errorf("加密方法错误")
			}
			return tokenKey, nil
		})
		if err != nil { //假的token
			r.String(401, "err token")
			r.Abort()
			return
		} else { //加密方法没问题
			if token.Valid { //验证token过期没，没过期，放行并且路由中传递token中携带的信息
				// 解析拿到token中携带参数,还得判断参数中的信息是否是正确的，不用了，账号密码错误，根本就不返回token了
				// verify_email, _ := Jwtparse(access_token)
				// 查询数据库有没有这个邮箱，进行后端验证
				r.Next()
				return

			}
		}
	}
	r.JSON(403, gin.H{
		"error": "no token",
	})
	r.Abort()
}
