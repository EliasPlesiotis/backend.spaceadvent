package database

import (
	"github.com/jinzhu/gorm"
)

// Message for Messages Page
type Message struct {
	gorm.Model
	UserID string
	Mine   bool
	Text   string
}
