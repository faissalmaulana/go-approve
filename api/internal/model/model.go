package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(100);not null"`
	Handler  string `gorm:"type:varchar(50);unique;not null"`
	Email    string `gorm:"type:varchar(255);unique;not null"`
	Password string `gorm:"type:varchar(60);not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	return nil
}

type BloclistToken struct {
	Token string `gorm:"unique;type:varchar(255);not null"`
}
