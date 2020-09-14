package Services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	mymux "github.com/gorilla/mux"
	"github.com/tidwall/gjson"
	"goMicroservice/gokit_pro1/util"
	"io/ioutil"
	"net/http"
	"strconv"
)

// 收到外部请求，对外部进行解码，外部请求的格式可能是rpc，http等，参数可能是json或者url等然后封装到endpoint的请求对象中

func DocodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	vars := mymux.Vars(r)
	if uid, ok := vars["uid"]; ok {
		uid, _ := strconv.Atoi(uid)
		return UserRequest{
			Uid:    uid,
			Method: r.Method,
			Token:  r.URL.Query().Get("token"),
		}, nil
	}
	return nil, errors.New("参数错误")
}

// 将我们发出的响应码进行编码，比较通用的是json
func EncodeUserResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Context-type", "application/json")
	if response == nil {
		fmt.Println("aaa")
	}
	return json.NewEncoder(w).Encode(response)
}

func MyErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain;charset=utf-8", []byte(err.Error())
	w.Header().Set("content-type", contentType)
	fmt.Println(body)
	if myerr, ok := err.(*util.MyError); ok {
		w.WriteHeader(myerr.Code)
		w.Write(body)
	} else {
		w.WriteHeader(500)
		w.Write(body)
	}

	//w.WriteHeader(500)
	//w.Write(body)
}

// 第三版
func DecodeAccessRequest(c context.Context, r *http.Request) (interface{}, error) {
	body, _ := ioutil.ReadAll(r.Body)
	result := gjson.Parse(string(body)) //第三方库解析json
	if result.IsObject() {              //如果是json就返回true
		username := result.Get("username")
		userpass := result.Get("userpass")
		return AccessRequest{Username: username.String(), Userpass: userpass.String(), Method: r.Method}, nil
	}
	return nil, errors.New("参数错误")

}
