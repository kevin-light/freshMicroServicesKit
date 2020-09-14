package main

import (
	"goMicroservice/gokit_pro1/util"
	"log"
)

func main() {
	// 生成私钥公钥文件
	err := util.GenRSAPubAndPri(1024, "./pem")
	if err != nil {
		log.Fatal(err)
	}
}
