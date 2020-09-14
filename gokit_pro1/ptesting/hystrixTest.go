package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"math/rand"
	"sync"
	"time"
)

type Product struct {
	ID    int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 {
		//time.Sleep(time.Second*3)
	}
	return Product{
		ID:    101,
		Title: "Golang精通",
		Price: 100,
	}, nil
}
func RecProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 {
		time.Sleep(time.Second * 3)
	}
	return Product{
		ID:    10,
		Title: "降级处理:静态数据",
		Price: 333,
	}, nil
}
func main() {
	rand.Seed(time.Now().UnixNano())
	configA := hystrix.CommandConfig{
		Timeout: 4000, // h配置请求延迟时间 ms
	}
	configB := hystrix.CommandConfig{
		Timeout:                2000,                   // h配置请求延迟时间 ms
		MaxConcurrentRequests:  5,                      // 最大并发数
		RequestVolumeThreshold: 5,                      // 默认20，熔断器请求阀值，有20个请求才进行错误百分比计算
		ErrorPercentThreshold:  10,                     // 错误百分比默认50%，
		SleepWindow:            int(time.Second * 100), //
	}
	hystrix.ConfigureCommand("get_prod", configA)  // get_prod是command的名字
	hystrix.ConfigureCommand("get_prodB", configB) // get_prod是command的名字
	c, _, _ := hystrix.GetCircuit("get_prodB")     //返回熔断器是否打开
	restluChan := make(chan Product, 1)            // 熔断器打开
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ { // 携程执行，模块
		go (func() {
			wg.Add(1)
			defer wg.Done()
			//err := hystrix.Do("get_prodB", func() error {
			errs := hystrix.Go("get_prodB", func() error { // 开启异步执行 携程
				p, _ := getProduct() // 随机延迟3s
				restluChan <- p
				fmt.Println(p)
				return nil
			}, func(err error) error { // 延迟超时进行降级处理 ,返回error
				//Rec := Product{
				//	ID:001,
				//	Title: "降级处理:静态数据",
				//	Price: 700,
				//}

				Rec, err := RecProduct()
				restluChan <- Rec
				fmt.Println(Rec)
				//return errors.New("rpc time out--03")
				return err
			})
			select {
			case getProd := <-restluChan:
				fmt.Println(getProd, "---03")
			case err := <-errs:
				fmt.Println(err, "---04")

			}
			if errs != nil {
				fmt.Println(errs, "---001")
			}
			fmt.Println(c.IsOpen())
			time.Sleep(time.Second * 1)
		})()
	}
	wg.Wait()
	//for {		循环执行
	//	//err := hystrix.Do("get_prodB", func() error {
	//	errs := hystrix.Go("get_prodB", func() error {		// 开启异步执行 携程
	//		p,_:=getProduct()	// 随机延迟3s
	//		restluChan<-p
	//		fmt.Println(p)
	//		return nil
	//	}, func(err error) error {	// 延迟超时进行降级处理 ,返回error
	//		Rec := Product{
	//			ID:001,
	//			Title: "降级处理:静态数据",
	//			Price: 700,
	//		}
	//		restluChan<-Rec
	//		fmt.Println(Rec)
	//		return errors.New("rpc time out--03")
	//	})
	//	select {
	//	case getProd := <-restluChan:
	//		fmt.Println(getProd ,"---03")
	//	case err:=<-errs:
	//		fmt.Println(err,"---04")
	//
	//	}
	//	if errs != nil {
	//		fmt.Println(errs,"---001")
	//	}
	//}
}
