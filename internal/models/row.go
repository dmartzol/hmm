package models

import "time"

type Row struct {
	ID         int64
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}
