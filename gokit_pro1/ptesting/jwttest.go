package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	Uname string `json:"username"`
	jwt.StandardClaims
}

func main() {
	sec := []byte("123abc") // 对称秘钥设置
	token_obj := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{Uname: "alvin"})
	token, _ := token_obj.SignedString(sec)
	fmt.Println(token, "---01")
	// 验证方式2  得到map[username:alvin]
	getToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return sec, nil
	})
	if getToken.Valid { // 获取解析token成功
		fmt.Println(getToken.Claims, "---03")
	}
	//// 验证方式2 得到Uname值
	//unc:=UserClaim{}
	//getToken,_:=jwt.ParseWithClaims(token,&unc, func(token *jwt.Token) (interface{}, error) {
	//	return sec,nil
	//})
	//if getToken.Valid{		// 获取解析token成功
	//	fmt.Println(getToken.Claims.(*UserClaim).Uname,"---03")
	//}
}
