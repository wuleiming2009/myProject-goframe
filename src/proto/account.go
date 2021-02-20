package proto

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResp struct {
	UserId uint64 `json:"user_id"`
	Token   string `json:"token"`
}

type SignUpArgs struct {
	Password    string `json:"password" binding:"required"` //用户密码
	Email       string `json:"email" binding:"required"`    //注册邮箱
	FacebookId  string `json:"facebook_id"`                 //facebook的三方账号
	GooleplusId string `json:"gooleplus_id"`                //google+的三方账号
}

type ResetPwdReq struct {
	Password string `json:"password" binding:"required"` //用户密码
}
