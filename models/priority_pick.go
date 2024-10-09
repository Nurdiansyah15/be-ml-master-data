package models

type PriorityPick struct {
	PriorityPickID uint `gorm:"primaryKey;autoIncrement"`
	HeroID         uint `gorm:"not null"`
	MatchID        uint `gorm:"not null"`
	TeamID         uint `gorm:"not null"`
	Total          int
	Role           string  `gorm:"size:50"`
	RatePick       float64 `gorm:"type:decimal(5,2)"`
}
