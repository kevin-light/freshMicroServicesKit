package Services

// 定义请求和响应的struct
type UserRequest struct {
	Uid    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

//func GenUserEndpoint(userService IUserService)endpoint.Endpoint  {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		r := request.(UserRequest)	// 进行断言
//		result := "nothing"
//		if r.Method == "GET"{
//			result = userService.GetName(r.Uid)
//		}else if r.Method == "DELETE"{	// 执行删除
//			err := userService.DelUser(r.Uid)
//			if err != nil {
//				result = err.Error()
//			}else {
//				result = fmt.Sprint("%d的用户删除成功",r.Uid)
//			}
//		}
//		return UserResponse{Result: result},nil
//	}
//}

//// 将我们发出的响应进行编码
//func EncodeUserResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
//	w.Header().Set("Content-type","application/json")	// 设置请求头的数据返回格式
//	return json.NewEncoder(w).Encode(response)
//}
