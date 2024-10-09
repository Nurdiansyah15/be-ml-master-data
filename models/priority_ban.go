package models

type PriorityBan struct {
	PriorityBanID  uint    `gorm:"primaryKey;autoIncrement"`
	HeroID         uint    `gorm:"not null"`
	MatchID        uint    `gorm:"not null"`
	TeamID         uint    `gorm:"not null"`
	Total          int
	Role           string  `gorm:"size:50"`
	RateBan        float64 `gorm:"type:decimal(5,2)"`
}
