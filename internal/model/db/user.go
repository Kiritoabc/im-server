package db

import (
	"time"
)

// User 用户表结构体
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`                // 主键
	PhoneNumber  string    `gorm:"unique;not null" json:"phone_number"` // 用户电话号码，唯一，不能为空
	Email        string    `gorm:"unique;not null" json:"email"`        // 用户邮箱，唯一，不能为空
	Username     string    `gorm:"not null" json:"username"`            // 用户名，不能为空
	PasswordHash string    `gorm:"not null" json:"password_hash"`       // 密码哈希值，不能为空
	AvatarURL    string    `gorm:"default:''" json:"avatar_url"`        // 用户头像URL，允许为空
	Bio          string    `gorm:"default:''" json:"bio"`               // 用户个人简介，允许为空
	Gender       string    `gorm:"not null" json:"gender"`              // 用户性别，不能为空,male、female、other
	Address      string    `gorm:"default:''" json:"address"`           // 用户住址，允许为空
	City         string    `gorm:"default:''" json:"city"`              // 用户所在城市，允许为空
	State        string    `gorm:"default:''" json:"state"`             // 用户所在州/省，允许为空
	Country      string    `gorm:"default:''" json:"country"`           // 用户所在国家，允许为空
	PostalCode   string    `gorm:"default:''" json:"postal_code"`       // 用户邮政编码，允许为空
	DateOfBirth  string    `gorm:"default:''" json:"date_of_birth"`     // 用户出生日期，允许为空
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`    // 记录创建时间，默认为当前时间
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`    // 记录更新时间，在更新时自动设置为当前时间
	//DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"` // 删除时间
}

func Register(user User) error {
	return DB.Create(&user).Error
}

func Login(username, passwordHash string) (*User, error) {
	var user User
	err := DB.Where("username = ? AND password_hash = ?", username, passwordHash).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
