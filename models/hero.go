package models

type Hero struct {
	HeroID    uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:50;not null"`
	HeroImage string `gorm:"size:255"`
}
