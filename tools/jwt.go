package tools

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var tokenKey = []byte("mysite.kunzige")

// 根据时间和map进行加密
func Generatejwt(timelong int64, args map[string]interface{}) string {
	security := jwt.New(jwt.SigningMethodHS256)
	claims := security.Claims.(jwt.MapClaims)
	for k, v := range args {
		claims[k] = v
	}
	claims["exp"] = timelong
	tokenString, err := security.SignedString(tokenKey)
	if err != nil {
		fmt.Println("加密失败")
	} else {
		return tokenString
	}
	return ""
}

func Jwtparse(refresh_token string) (map[string]interface{}, error) {
	token, err := jwt.Parse(refresh_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tokenKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析token出错")
	}
	if token.Valid {
		claims, _ := token.Claims.(jwt.MapClaims)
		return claims, nil
	}
	return nil, fmt.Errorf("token is invalid")

}

func Regenerate(timelong int64, refresh_token string) string {
	// 解析出token中的数据，重新生成有效的token
	token_map, err := Jwtparse(refresh_token)
	if err != nil {
		return "error_refresh_token"
	}
	new_map := map[string]interface{}{}
	for k, v := range token_map {
		if k != "exp" {
			new_map[k] = v
		}
	}
	return Generatejwt(timelong, new_map)

}
