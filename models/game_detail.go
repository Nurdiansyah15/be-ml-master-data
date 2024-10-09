package models

type GameDetail struct {
	GameDetailID uint   `gorm:"primaryKey;autoIncrement"`
	GameID       uint   `gorm:"not null"`
	TeamID       uint   `gorm:"not null"`
	HeroID       uint   `gorm:"not null"`
	Role         string `gorm:"size:50"`
	EarlyResult  string `gorm:"size:20"`
	Position     string `gorm:"size:50"`
}
