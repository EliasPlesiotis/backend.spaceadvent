package database

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	UserName string
	Password string
	Email    string
	Messages []Message `gorm:"ForeignKey:UserID"`
}
