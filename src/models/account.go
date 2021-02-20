package models

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Account struct {
	Id          uint64    `json:"id" gorm:"column:id;primary_key"`
	UserId     uint64    `json:"user_id" gorm:"column:user_id;not null"`         //用户雪花ID
	Password    string    `json:"password" gorm:"column:password;not null"`         //用户密码 md5
	Email       string    `json:"email" gorm:"column:email;not null"`               //注册邮箱
	FacebookId  string    `json:"facebook_id" gorm:"column:facebook_id;not null"`   //facebook的三方账号
	GooleplusId string    `json:"gooleplus_id" gorm:"column:gooleplus_id;not null"` //google+的三方账号
	IsDeleted   int8      `json:"is_deleted" gorm:"column:is_deleted;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

func (*Account) TableName() string {
	return "account"
}

func AddAccount(ctx context.Context, tbAccount *Account, tx ...*gorm.DB) (uint64, error) {
	var db *gorm.DB
	switch len(tx) {
	case 0:
		db = DBInstance()
	case 1:
		db = tx[0]
	}
	if db == nil {
		return 0, errors.New("empty db")
	}

	err := db.Create(tbAccount).Error
	if err != nil {
		return 0, err
	}
	return tbAccount.Id, nil
}

func AccountById(ctx context.Context, id uint64, tx ...*gorm.DB) (*Account, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx[0]
	} else {
		db = DBInstance()
	}
	tbAccount := &Account{}
	err := db.First(tbAccount, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	return tbAccount, nil
}

func SearchAccount(ctx context.Context, tx *gorm.DB, offset, limit uint64, order, query string, queryArgs ...interface{}) ([]*Account, error) {
	db := tx
	if db == nil {
		db = DBInstance()
	}
	switch query {
	case "":
		query = "is_deleted = 0"
	default:
		query += " and is_deleted = 0"
	}
	qs := db.Where(query, queryArgs...)
	if offset != 0 {
		qs.Offset(offset)
	}
	if limit != 0 {
		qs.Limit(limit)
	}
	if order != "" {
		qs.Order(order)
	}
	account := make([]*Account, 0)
	err := qs.Find(&account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func UpdateAccount(ctx context.Context, id uint64, account *Account, tx ...*gorm.DB) error {
	db, err := switchDB(tx...)
	if err != nil {
		return err
	}
	return db.Model(&Account{}).Where("id=?", id).Updates(account).Error
}
