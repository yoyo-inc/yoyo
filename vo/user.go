package vo

type QueryUser struct {
	// 账户名
	Username string `form:"username"`
	// 手机号
	Phone string `form:"phone"`
}
