package model

import (
	"time"
)

type BaseModel struct {
	ID       	uint64 		`gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}