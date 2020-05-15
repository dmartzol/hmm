package models

import "time"

type ConfirmationCodeType int

const (
	ConfirmationCodeTypeEmail ConfirmationCodeType = iota
	ConfirmationCodeTypePhone
	ConfirmationCodeTypePasswordReset
)

type ConfirmationCode struct {
	Row
	Type        ConfirmationCodeType
	AccountID   int64 `db:"account_id"`
	Code        string
	ConfirmTime *time.Time `db:"confirm_time"`
	ExpireTime  *time.Time `db:"expire_time"`
}
