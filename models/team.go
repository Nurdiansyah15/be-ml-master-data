package models

type Team struct {
	TeamID uint   `gorm:"primaryKey;autoIncrement"`
	Name   string `gorm:"size:50;not null"`
	Logo   string `gorm:"size:255"`
}
