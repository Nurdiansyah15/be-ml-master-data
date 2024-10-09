package models

type Coach struct {
	CoachID uint   `gorm:"primaryKey;autoIncrement"`
	TeamID  uint   `gorm:"not null"`
	Name    string `gorm:"size:100;not null"`
	Role    string `gorm:"size:50"`
	Image   string `gorm:"size:255"`
}
