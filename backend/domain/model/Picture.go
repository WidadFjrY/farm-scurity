package model

import "time"

type Picture struct {
	ID                    string `gorm:"primaryKey;type:char(15);not null"`
	URL                   string `gorm:"type:char(100);not null"`
	IsFromMotionDetection bool   `gorm:"type:boolean;default:false"`
	CreatedAt             time.Time
}
