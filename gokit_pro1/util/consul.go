package util

import (
	"fmt"
	consul_api "github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"log"
)

var consulClient *consul_api.Client
var ServiceID string
var ServicesName string
var ServicePort int

func init() {
	config := consul_api.DefaultConfig()
	config.Address = "192.168.31.129:8500"
	// 服务注册
	client, err := consul_api.NewClient(config)
	if err != nil {
		log.Fatal(err, "---01")
	}
	consulClient = client

	ServiceID = "userService" + uuid.NewV1().String()
}
func SetServiceNamePort(name string, port int) {
	ServicesName = name
	ServicePort = port
}

// 服务注册
func RegServices() {
	// 创建服务注册信息
	reg := consul_api.AgentServiceRegistration{}
	reg.ID = ServiceID // 不能重复
	reg.Name = ServicesName
	reg.Address = "192.168.142.1" // 	// API接口地址
	//reg.Address = "192.168.31.129"		// consul服务地址error
	reg.Port = ServicePort
	//reg.Port = 8081
	reg.Tags = []string{"consul_test"}

	check := consul_api.AgentServiceCheck{}
	check.Interval = "5s"                                                     // 检查间隔时间
	check.HTTP = fmt.Sprintf("http://%s:%d/health", reg.Address, ServicePort) // API接口地址
	//check.HTTP = "http://192.168.31.129:8080/health"
	reg.Check = &check

	err := consulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err, "---02")
	}
}
func UnRegService() {
	// 反注册
	consulClient.Agent().ServiceDeregister("ServiceID")
}
