package main

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httpstransports "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"goMicroservice/gokit_pro2/Services"
	"goMicroservice/gokit_pro2/utils"
	"io"
	"net/url"
	"os"
	"time"
)

func main() { // 封装整合main3
	configA := hystrix.CommandConfig{
		Timeout:                2000,
		MaxConcurrentRequests:  5,
		RequestVolumeThreshold: 3,
		ErrorPercentThreshold:  20,
		SleepWindow:            int(time.Second * 100),
	}
	hystrix.ConfigureCommand("getuser", configA)
	err := hystrix.Do("getuser", func() error {
		res, err := utils.GetUser()
		fmt.Println(res)
		return err
	}, func(e error) error {
		fmt.Println("降级用户")
		return e
	})

	if err != nil {

	}

	fmt.Println(err)
}

func main_t02() { // 添加熔断器之前： 后面在UserUtils.go添加熔断器
	//curl --request PUT http://192.168.31.129:8500/v1/agent/service/deregister/userServiceID  删除service
	{
		config := consulapi.DefaultConfig()
		config.Address = "192.168.31.129:8500" // 注册中心地址
		apiClient, _ := consulapi.NewClient(config)
		client := consul.NewClient(apiClient)
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
		}
		{
			tags := []string{"consul_test"}
			// 实时查询服务实例状态
			instances := consul.NewInstancer(client, logger, "userService", tags, true) //
			{
				factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) {

					tgt, _ := url.Parse("http://" + service_url) // 真实服务器地址
					return httpstransports.NewClient("GET", tgt, Services.GetUserInfoRequest, Services.GetUserInfoResponse).Endpoint(), nil, nil
				}
				endpointer := sd.NewEndpointer(instances, factory, logger)
				endpoiints, _ := endpointer.Endpoints()

				fmt.Println("服务有：", len(endpoiints))
				//mylb := lb.NewRoundRobin(endpointer) // 轮询
				mylb := lb.NewRandom(endpointer, time.Now().UnixNano()) // 随机
				for {
					//getUserInfo := endpoiints[0]		 // 写死获取第一个api服务
					getUserInfo, _ := mylb.Endpoint() //
					// 第三，四步，创建context上下文，并执行
					res, err := getUserInfo(context.Background(), Services.UserRequest{Uid: 102})
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					// 第五步 断言，得到响应
					userInfo := res.(Services.UserResponse)
					fmt.Println(userInfo.Result)
					time.Sleep(3 * time.Second)
				}
			}
		}
	}
}

// 直接连接API方式
func main_t01() {
	tgt, _ := url.Parse("http:://localhost::8081")
	// 1、创建一个直连client，必须写两个func，一个是如何请求，一个是响应我们怎么处理
	client := httpstransports.NewClient("GET", tgt, Services.GetUserInfoRequest, Services.GetUserInfoResponse)
	// 2、暴露出endpoint的func，方便执行
	getUserInfo := client.Endpoint()
	ctx := context.Background() // 3、创建context上下文对象
	// 4、 执行
	res, err := getUserInfo(ctx, Services.UserRequest{Uid: 101})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 5、断言得到响应
	userinfo := res.(Services.UserResponse)
	fmt.Println(userinfo.Result)
}
