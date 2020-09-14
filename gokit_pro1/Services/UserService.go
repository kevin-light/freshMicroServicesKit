package Services

import "errors"

// 需求：传入用户的用户ID，获取用户的用户名
// 业务
type IUserService interface {
	GetName(userId int) string
	DelUser(userId int) error
}

type UserService struct {
}

func (us UserService) GetName(userId int) string {
	if userId == 101 {
		return "101"
	}
	return "guest"
}

func (us UserService) DelUser(userId int) error {
	if userId == 101 {
		return errors.New("无权限")
	}
	return nil
}
