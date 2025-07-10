package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	Todo       TaskStatus = "todo"
	InProgress TaskStatus = "in_progress"
	Done       TaskStatus = "done"
)

type Task struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey;default:uuid()" json:"id"`
	ProjectID   uuid.UUID  `gorm:"type:char(36);not null" json:"project_id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      TaskStatus `gorm:"type:enum('todo','in_progress','done');default:'todo';not null" json:"status"`
	Deadline    *time.Time `json:"deadline" time_format:"2006-01-02"`
	Project     Project    `gorm:"foreignKey:ProjectID"`
}
