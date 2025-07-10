package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey;default:uuid()" json:"id"`
	Name     string    `gorm:"type:varchar(255);not null" json:"name"`
	Email    string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password string    `gorm:"type:varchar(255);not null" json:"-"`
}
