package model

type Device struct {
	ID       string `gorm:"primaryKey;type:char(15);not null"`
	Location string `gorm:"type:char(100);not null"`
	IsActive bool   `gorm:"type:bool;default:true;"`
}
