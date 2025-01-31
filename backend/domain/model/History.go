package model

import "time"

type History struct {
	ID         string   `gorm:"primaryKey;type:char(15);not null"`
	Operation  string   `gorm:"type:char(255);not null"`
	PictureID  *string  `gorm:"type:char(15);unique"` // Nullable foreign key
	URLPicture *Picture `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PictureID;references:ID"`
	CreatedAt  time.Time
}
