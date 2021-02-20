package bcrypt

import (
	"golang.org/x/crypto/bcrypt"

	"myProject/common/response"
)

func GeneratePassword(password string) (string, error) {
	// 生成加密密码 - 这里使用了Golang官方加密包
	// golang.org/x/crypto/bcrypt
	// 注意 - 每次生成的密码都不一样，但是没关系
	// 登陆是验证原密码再加密后是否能和存储的加密是否一致
	// bcrypt.DefaultCost = 10
	encodePW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", response.ErrInternal
	}
	return string(encodePW), nil
}

func ComparePassword(userPwd, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(userPwd), []byte(password))
}
