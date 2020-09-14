package main

import (
	"flag"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	httpstransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"goMicroservice/gokit_pro1/Services"
	"goMicroservice/gokit_pro1/util"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// go run main.go --name userservices -p 8080 带参数启动
	// t3 post请求http://192.168.142.1:8083/access-token 参数{"username":"alvin","userpass":"123"}
	//  get请求 http://192.168.142.1:8083/user/101?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFsdmluIiwiZXhwIjoxNjAwMDAzNzUwfQ.W9NcSjeMZ4WT6ZpCgAnfF9a1BQUIsNKJXcESncu3FDY
	name := flag.String("name", "", "服务名称")
	port := flag.Int("p", 0, "服务端口")
	flag.Parse()
	if *name == "" {
		log.Fatal("请指定服务名称")
	}
	if *port == 0 {
		log.Fatal("请指定端口")
	}
	util.SetServiceNamePort(*name, *port) // 设置name,*port
	var logger kitlog.Logger              // 中间件log配置
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)                       // kit日志打印到控制台
		logger = kitlog.WithPrefix(logger, "my gokit_pro", "1.0")        // predix：项目名=版本号
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC) // utc时间设置
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)     // 输出文件名和行号
	}
	fmt.Println("server start")
	user := Services.UserService{} // 用户服务
	//endp := Services.GenUserEndpoint(user)
	limit := rate.NewLimiter(1, 5)                                                                                                                    // 限流
	endp := Services.RateLimit(limit)((Services.CheckTokenMiddleware())((Services.UserServiceLogMiddleware(logger))(Services.GenUserEndpoint(user)))) //这里添加中间件 为endpoint添加集成限流和log功能
	//endp:=Services.RateLimit(limit)((CheckTokenMiddleware())((UserServiceLogMiddleware(logger))(GenUserEnpointMiddleware(user))))
	options := []httpstransport.ServerOption{ //生成ServerOtion切片，传入我们自定义的错误处理函数
		httpstransport.ServerErrorEncoder(Services.MyErrorEncoder),
	}
	// 把我们自己创建的services发布为http服务
	serverHandler := httpstransport.NewServer(endp, Services.DocodeUserRequest, Services.EncodeUserResponse, options...)
	// 权限服务
	accessService := &Services.AccessServer{}
	endp_access := Services.AccessEndpoint(accessService)
	accessHandler := httpstransport.NewServer(endp_access, Services.DecodeAccessRequest, Services.EncodeUserResponse, options...)

	router := mymux.NewRouter() // 使用第三方路由mux
	{
		router.Methods("POST").Path("/access-token").Handler(accessHandler)
		router.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
		router.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type", "application/json")
			writer.Write([]byte(`{"status":"ok"}`))
		})
	}
	errChan := make(chan error)
	//开启协程
	go func() {
		util.RegServices() //consul 注册中心注册服务
		err := http.ListenAndServe(":"+strconv.Itoa(*port), router)
		//err := http.ListenAndServe(":8081",router)
		if err != nil {
			log.Panicln(err)
			//如果有异常，清除注册中心的服务
			errChan <- err
		}
	}()
	// 信号监听
	go (func() {
		sig_c := make(chan os.Signal)
		signal.Notify(sig_c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sig_c)
	})()
	//从通道读取信号
	getErr := <-errChan
	util.UnRegService()
	log.Println(getErr)
}
