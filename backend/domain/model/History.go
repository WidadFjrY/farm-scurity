package model

import "time"

type History struct {
	ID          string  `gorm:"primaryKey;type:char(15);not null"`
	Operation   string  `gorm:"type:char(255);not null"`
	Description string  `gorm:"type:char(255);not null"`
	PictureID   *string `gorm:"type:char(15);unique"`
	CreatedAt   time.Time
}
