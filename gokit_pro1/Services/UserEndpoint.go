package Services

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"goMicroservice/gokit_pro1/util"
	"golang.org/x/time/rate"
	"strconv"
	"time"
)

// 定义请求和响应的struct
type UserRequest struct {
	Uid    int `json:"uid"`
	Method string
	Token  string
}

type UserResponse struct {
	Result string `json:"result"`
}

// 日志中间件,每一个service都应该有自己的日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			logger.Log("method", r.Method, "event", "get user", "userId", r.Uid)
			return next(ctx, request)
		}
	}
}

// 加入限流功能中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, util.NewMyError(429, "too many requests")
				//return nil,errors.New("too many requests")
				//return nil,errors.NewError(429,"too many requests")
			}
			return next(ctx, request)
		}
	}
}

const secKey = "123abc" // 设置秘钥
type UserClaim struct {
	Uname string `json:"username"`
	jwt.StandardClaims
}
type IAccessService interface {
	GetToken(uname string, upwd string) (string, error)
}
type AccessServer struct {
}

func (a *AccessServer) GetToken(uname string, upwd string) (string, error) {
	if uname == "alvin" && upwd == "123" {
		userinfo := &UserClaim{Uname: uname}
		userinfo.ExpiresAt = time.Now().Add(time.Second * 60).Unix() // 设置60s过期时间
		token_obj := jwt.NewWithClaims(jwt.SigningMethodHS256, userinfo)
		token, err := token_obj.SignedString([]byte(secKey))
		return token, err
	}
	return "", fmt.Errorf("error username and pwd")
}

func GenUserEndpoint(userService IUserService) endpoint.Endpoint {
	//var logger log.Logger		// 项目中用中间件代替
	//{
	//	logger = log.NewLogfmtLogger(os.Stdout)		// kit日志打印到控制台
	//	logger = log.WithPrefix(logger,"my gokit_pro","1.0")	// predix：项目名=版本号
	//	logger = log.With(logger,"time",log.DefaultTimestampUTC)	// utc时间设置
	//	logger = log.With(logger,"caller",log.DefaultCaller)	// 输出文件名和行号
	//}
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest) // 进行断言
		fmt.Println("登录username", ctx.Value("LoginUser"))
		result := "nothing"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid) + strconv.Itoa(util.ServicePort)
			//logger.Log("method",r.Method,"event","get user","userId",r.Uid)
		} else if r.Method == "DELETE" { // 执行删除
			err := userService.DelUser(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprint("%d的用户删除成功", r.Uid)
			}
		}
		return UserResponse{Result: result}, nil
	}
}

type AccessRequest struct {
	Username string
	Userpass string
	Method   string
}
type AccessResponse struct {
	Status string
	Token  string
}

func AccessEndpoint(accessService IAccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(AccessRequest)
		result := AccessResponse{Status: "OK"}
		if r.Method == "POST" {
			token, err := accessService.GetToken(r.Username, r.Userpass)
			if err != nil {
				result.Status = "error:" + err.Error()
			} else {
				result.Token = token
			}
		}
		return result, nil
	}
}

// token 验证中间件
func CheckTokenMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest) // 通过类型断言获取请求结构体
			uc := UserClaim{}
			//下面r.Token是在DecodeUserRequest封装的
			getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (interface{}, error) {
				return []byte(secKey), err
			})
			fmt.Println(err, 123)
			if getToken != nil && getToken.Valid { //验证通过
				newCtx := context.WithValue(ctx, "LoginUser", getToken.Claims.(*UserClaim).Uname)
				return next(newCtx, request)
			} else {
				return nil, util.NewMyError(403, "error token")
			}
		}
	}
}
