package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

type UserClaims struct {
	Uname string `json:"username"`
	jwt.StandardClaims
}

func main() {
	// 非对称秘钥设置
	priKeyBytes, err := ioutil.ReadFile("../pem/private.pem")
	if err != nil {
		fmt.Println("私钥文件读取失败")
	}
	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priKeyBytes)
	if err != nil {
		fmt.Println("私钥不正确")
	}
	pubKeyBytes, err := ioutil.ReadFile("../pem/public.pem")
	if err != nil {
		fmt.Println("公钥文件读取失败")
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		fmt.Println("公钥不正确")
	}
	//// 验证方式2  得到map[username:alvin]
	//token_obj := jwt.NewWithClaims(jwt.SigningMethodHS256,UserClaims{Uname: "alvin"})
	//token,_:=token_obj.SignedString(priKey)
	//fmt.Println(token,"---01")
	//getToken,_:=jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	//	return pubKey,nil
	//})
	//if getToken.Valid{		// 获取解析token成功
	//	fmt.Println(getToken.Claims,"---03")
	//}
	//// 验证方式2 得到Uname值
	users := UserClaims{Uname: "alvin"}
	users.ExpiresAt = time.Now().Add(time.Second * 5).Unix() // 设置过期时间
	token_obj := jwt.NewWithClaims(jwt.SigningMethodRS256, users)
	token, err := token_obj.SignedString(priKey)
	if err != nil {
		fmt.Println(err, "---01")
	}
	fmt.Println(token, "---03")

	i := 0
	for {
		unc := UserClaims{}
		getToken, err := jwt.ParseWithClaims(token, &unc, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})
		if getToken != nil && getToken.Valid { // 获取解析token成功
			fmt.Println(getToken.Claims.(*UserClaims).Uname, "---03")
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("错误的token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				fmt.Println("token过期或者未启用")
			} else {
				fmt.Println("Couldn't handle this token:", err)
			}
		} else {
			fmt.Println("无法解析token", err)
		}
		i++
		fmt.Println(i, "---i")
		time.Sleep(time.Second)
	}
}
