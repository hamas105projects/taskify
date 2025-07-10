package models

import "github.com/google/uuid"

type Project struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey;default:uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedByID uuid.UUID `gorm:"type:char(36);not null" json:"created_by_id"`
	User        User      `gorm:"foreignKey:CreatedByID"`
}
