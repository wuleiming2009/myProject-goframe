package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"myProject/cache"
	"myProject/common/bcrypt"
	"myProject/common/response"
	"myProject/common/snow_flake"
	"myProject/common/token"
	"myProject/common/verify"
	"myProject/lib/log"
	"myProject/models"
	"myProject/proto"
)

type AccountService interface {
	Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginResp, error)
	SignUp(ctx context.Context, req *proto.SignUpArgs) (*proto.LoginResp, error)

	FetchAccountByEmail(ctx context.Context, email string) (*models.Account, error)
}

type accountService struct {
}

func NewAccountService() (AccountService, error) {
	return &accountService{}, nil
}

func (s *accountService) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginResp, error) {

	zl := log.FromContext(ctx)

	account, err := s.FetchAccountByEmail(ctx, req.Email)
	if err != nil {
		zl.Errorf("Login failure, req:%+v err:%v", req, err)
		return nil, err
	}

	err = bcrypt.ComparePassword(account.Password, req.Password)
	if err != nil {
		zl.Errorf("ComparePassword failure, account:%+v err:%v", account, err)
		return nil, response.ErrInternal
	}

	// 登录token
	tk, err := token.AddToken(ctx, account.UserId)
	if err != nil {
		zl.Errorf("AddToken failure, account:%+v err:%v", account, err)
		return nil, err
	}

	return &proto.LoginResp{
		UserId: account.UserId,
		Token:   tk,
	}, nil
}

func (s *accountService) SignUp(ctx context.Context, req *proto.SignUpArgs) (*proto.LoginResp, error) {

	zl := log.FromContext(ctx)
	zl.Infof("SignIn req:%+v", req)

	if !verify.VerifyPassword(req.Password) {
		return nil, response.ErrInvalidArgs
	}

	// email地址 格式校验
	valid := verify.VerifyEmailFormat(req.Email)
	if !valid {
		return nil, response.ErrInvalidArgs
	}

	// email地址 去重校验
	account, err := s.FetchAccountByEmail(ctx, req.Email)
	if account != nil {
		zl.Errorf("AccountByEmail this email:%v had account", req.Email)
		return nil, response.ErrInternal
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		zl.Errorf("AccountByEmail failure, email:%v err:%v", req.Email, err)
		return nil, response.ErrInternal
	}

	userId := snow_flake.GenSnowFlake()

	pwd, err := bcrypt.GeneratePassword(req.Password)
	if err != nil {
		zl.Errorf("GeneratePassword failure, pwd:%v err:%v", req.Password, err)
		return nil, response.ErrInternal
	}

	// 插入用户数据
	acc := &models.Account{
		UserId:  userId,
		Email:    req.Email,
		Password: pwd,
	}
	acc.Id, err = models.AddAccount(ctx, acc)
	if err != nil {
		zl.Errorf("AddAccount failure, acc:%+v err:%v", acc, err)
		return nil, response.ErrInternal
	}

	// 登录token
	tk, err := token.AddToken(ctx, userId)
	if err != nil {
		zl.Errorf("AddToken failure, userId:%v err:%v", userId, err)
		return nil, response.ErrInternal
	}

	return &proto.LoginResp{
		UserId: userId,
		Token:   tk,
	}, nil
}

///////////// FetchWithJson ////////////////

func KeyAccountByEmail(email string) string {
	return fmt.Sprintf("KeyAccountByEmail:%s", email)
}

func (s *accountService) FetchAccountByEmail(ctx context.Context, email string) (*models.Account, error) {

	zl := log.FromContext(ctx)

	rdsCli, err := cache.RedisClient()
	if err != nil {
		zl.Errorf("RedisClient failure, err:%v", err)
		return nil, err
	}

	key := KeyAccountByEmail(email)
	ret := &models.Account{}
	err = cache.FetchWithJson(ctx, rdsCli, key, 1*time.Minute, ret, func() (interface{}, error) {

		list, err := models.SearchAccount(ctx, nil, 0, 1, "", "email=?", email)
		if err != nil {
			zl.Errorf("SearchAccount failure, email:%v, err:%v", email, err)
			return nil, err
		}

		if len(list) == 0 {
			zl.Errorf("SearchAccount record not found, email:%v, err:%v", email, err)
			err = gorm.ErrRecordNotFound
			return nil, err
		}

		if len(list) != 1 {
			zl.Errorf("SearchAccount failure, len:%v, email:%v, err:%v", len(list), email, err)
			return nil, errors.New("Search account by email length wrong")
		}
		return list[0], nil

	})

	if err != nil {
		zl.Errorf("AccountByEmail FetchWithJson failure, err:%v", err)
		return nil, err
	}

	return ret, nil
}
