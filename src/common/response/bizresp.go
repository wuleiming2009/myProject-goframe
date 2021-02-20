package response

import (
	"encoding/json"
)

type BizErr struct {
	Code int
	Msg  string
}

func (e *BizErr) Error() string {
	if jsonify, err := json.Marshal(e); err == nil {
		return string(jsonify)
	}
	return ""
}

func NewBizErr(code int, msg string) *BizErr {
	return &BizErr{
		Code: code,
		Msg:  msg,
	}
}

func (e *BizErr) WithMessage(msg string) *BizErr {
	return NewBizErr(e.Code, msg)
}

// 正常返回
var (
	SuccEmpty = ""
)

// 业务错误
var (
	ErrInvalidArgs        = NewBizErr(40000, "Invalid Args")
	ErrNotWaitingDelivery = NewBizErr(40001, "This Order Is Not Waiting For Delivery")
	ErrOrderUpdated       = NewBizErr(40002, "This order is expired, please reload the page")
	ErrCanNotCancelOrder  = NewBizErr(40001, "This order can't be cancelled")
	ErrCanNotSignOrder    = NewBizErr(40001, "This order can't be signed")
)

// 系统错误
var (
	ErrInternal = NewBizErr(50000, "Something Wrong Here, Try Again Later")
	ErrEmptyDB  = NewBizErr(50010, "Connection Failed")
	ErrRelogin  = NewBizErr(50020, "Please Login Again")
)
