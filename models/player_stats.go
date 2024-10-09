package models

type PlayerStats struct {
	StatID    uint    `gorm:"primaryKey;autoIncrement"`
	PlayerID  uint    `gorm:"not null"`
	GameID    uint    `gorm:"not null"`
	GameRate  float64 `gorm:"type:decimal(5,2)"`
	MatchRate float64 `gorm:"type:decimal(5,2)"`
}
