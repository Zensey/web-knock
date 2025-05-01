package model

import "time"

type Blacklist struct {
	IP      string    `json:"-" gorm:"primaryKey"`
	Request time.Time `json:"request,omitempty"`
}
