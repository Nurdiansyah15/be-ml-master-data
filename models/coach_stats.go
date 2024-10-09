package models

type CoachStats struct {
	StatID    uint    `gorm:"primaryKey;autoIncrement"`
	CoachID   uint    `gorm:"not null"`
	MatchID   uint    `gorm:"not null"`
	GameRate  float64 `gorm:"type:decimal(5,2)"`
	MatchRate float64 `gorm:"type:decimal(5,2)"`
}
