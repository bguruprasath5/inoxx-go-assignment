package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName    string     `json:"user_name" gorm:"uniqueIndex"`
	AuthToken   string     `json:"auth_token"`
	Password    string     `json:"-"`
	FullName    string     `json:"full_name"`
	Dob         *time.Time `json:"dob"`
	LinkedinURL string     `json:"linkedin_url"`
}
