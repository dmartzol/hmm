package hmm

import "time"

type ConfirmationType int

const (
	ConfirmationTypeEmail ConfirmationType = iota
	ConfirmationTypePhone
	ConfirmationTypePasswordReset
)

type Confirmation struct {
	Row
	Type                     ConfirmationType
	AccountID                int64 `db:"account_id"`
	Key                      string
	ConfirmationTarget       *string    `db:"confirmation_target"`
	FailedConfirmationsCount int64      `db:"failed_confirmations_count"`
	ConfirmTime              *time.Time `db:"confirm_time"`
	ExpireTime               *time.Time `db:"expire_time"`
}

type ConfirmationService interface {
	PendingConfirmationByKey(key string) (*Confirmation, error)
	Confirm(id int64) (*Confirmation, error)
	FailedConfirmationIncrease(id int64) (*Confirmation, error)
}
