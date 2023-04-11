package model

import "gorm.io/gorm"

type User struct {
	*gorm.Model `json:"-"`
	UserID      uint32 `json:"user_id" gorm:"primaryKey"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	PasswordMD5 string `json:"password"`
	RealName    string `json:"real_name"`
	Bio         string `json:"bio"`
	Verified    bool   `json:"verified"`
	Deleted     bool   `json:"deleted"`
}
