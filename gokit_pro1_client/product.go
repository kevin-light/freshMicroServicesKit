package main

import (
	"fmt"
	httpstransport "github.com/go-kit/kit/transport/http"
	"net/url"
)

func main() {
	tgt, _ := url.Parse("http:://localhost::8081")
	client := httpstransport.NewClient("GET", tgt, enc, dec)
	fmt.Println(client)
}
