package vo

import "github.com/yoyo-inc/yoyo/models"

type UserVO struct {
	models.User
	Roles []int `json:"roles"`
}

type QueryUserVO struct {
	// 账户名
	Username string `form:"username"`
	// 手机号
	Phone string `form:"phone"`
}
