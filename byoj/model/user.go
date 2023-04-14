package model

import (
	"byoj/utils/logs"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	ID          uint32         `json:"user_id"    form:"user_id"    query:"user_id"   gorm:"primaryKey;unique;not null"`
	CreatedAt   time.Time      `json:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" form:"updated_at" query:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" form:"deleted_at" query:"deleted_at"`
	UserName    string         `json:"user_name"  form:"user_name"  query:"user_name" gorm:"unique;not null"`
	Email       string         `json:"email"      form:"email"      query:"email"     gorm:"unique;not null"`
	PasswordMD5 string         `json:"password"   form:"password"   query:"password"  gorm:"not null"`
	RealName    string         `json:"real_name"  form:"real_name"  query:"real_name" `
	Bio         string         `json:"bio"        form:"bio"        query:"bio" `
	Verified    bool           `json:"verified"   form:"verified"   query:"verified"  gorm:"not null"`
	Deleted     bool           `json:"deleted"    form:"deleted"    query:"deleted"   gorm:"not null"`
}

func UserRegister(userName string, email string, passwordMD5 string, realName string, bio string) error {
	m := GetModel()
	defer m.Close()

	result := m.tx.Create(&User{
		UserName:    userName,
		Email:       email,
		PasswordMD5: passwordMD5,
		RealName:    realName,
		Bio:         bio,
		Verified:    false,
		Deleted:     false,
	})
	if result.Error != nil {
		logs.Warn("Create user failed.", zap.Error(result.Error))
		m.Abort()
		return result.Error
	}

	m.tx.Commit()
	return nil
}

func UserVerify(userID uint32) error {
	m := GetModel()
	defer m.Close()

	result := m.tx.First(&User{}, userID).Updates(User{Verified: true})
	if result.Error != nil {
		logs.Warn("Update user's verification info failed.", zap.Error(result.Error))
		m.Abort()
		return result.Error
	}

	m.tx.Commit()
	return nil
}

func FindUserByID(userID uint32) (User, error) {
	m := GetModel()
	defer m.Close()

	var user User
	result := m.tx.First(&user, userID)
	if result.Error != nil {
		logs.Info("Find user by id failed.", zap.Error(result.Error))
		m.Abort()
		return user, result.Error
	}

	m.tx.Commit()
	return user, nil
}

func FindUserByName(userName string) (User, error) {
	m := GetModel()
	defer m.Close()

	var user User
	result := m.tx.Model(&User{}).Where("user_name = ?", userName).First(&user)
	if result.Error != nil {
		logs.Info("Find user by name failed.", zap.Error(result.Error))
		m.Abort()
		return user, result.Error
	}

	m.tx.Commit()
	return user, nil
}

func FindUserByEmail(email string) (User, error) {
	m := GetModel()
	defer m.Close()

	var user User
	result := m.tx.Model(&User{}).Where("email = ?", email).First(&user)
	if result.Error != nil {
		logs.Info("Find user by email failed.", zap.Error(result.Error))
		m.Abort()
		return user, result.Error
	}

	m.tx.Commit()
	return user, nil
}
