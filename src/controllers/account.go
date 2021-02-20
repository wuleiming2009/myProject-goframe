package controllers

import (
	"github.com/gin-gonic/gin"

	"myProject/common/response"
	"myProject/lib/log"
	"myProject/proto"
)

// @summary 注册账号
// @description 注册账号
// @tags account
// @param signInArgs body proto.SignUpArgs true "注册账号"
// @success 200 {object} proto.LoginResp "登录tk"
// @router /account/sign_up [post]
func AddAccount(c *gin.Context) {
	ctx := c.Request.Context()
	zl := log.FromContext(ctx)
	args := &proto.SignUpArgs{}
	err := c.BindJSON(args)
	if err != nil {
		response.Response(c, response.ErrInvalidArgs)
		return
	}

	resp, err := serviceAccount.SignUp(ctx, args)
	if err != nil {
		zl.Errorf("SignIn failure, args:%+v err:%v", args, err)
		response.Response(c, response.ErrInternal)
		return
	}

	response.Response(c, &proto.LoginResp{Token: resp.Token})

}

// @summary 账号登录
// @description 账号登录
// @tags account
// @param loginReq body proto.LoginReq true "账号登录"
// @success 200 {object} proto.LoginResp "登录tk"
// @router /account/login [post]
func UserLogin(c *gin.Context) {
	ctx := c.Request.Context()
	zl := log.FromContext(ctx)

	args := &proto.LoginReq{}
	err := c.BindJSON(args)
	if err != nil {
		response.Response(c, response.ErrInvalidArgs)
		return
	}

	if args.Email == "" || args.Password == "" {
		zl.Errorf("Email and password have empty, args:%+v", args)
		response.Response(c, response.ErrInvalidArgs)
		return
	}

	resp, err := serviceAccount.Login(ctx, args)
	if err != nil {
		response.Response(c, response.ErrInternal)
		return
	}

	response.Response(c, &proto.LoginResp{Token: resp.Token})
	return
}
