package models

type HeroBan struct {
	BanID            uint `gorm:"primaryKey;autoIncrement"`
	MatchID          uint `gorm:"not null"`
	TeamID           uint `gorm:"not null"`
	HeroID           uint `gorm:"not null"`
	TotalBans        int  `gorm:"not null"`
	FirstPhaseCount  int  `gorm:"not null"`
	SecondPhaseCount int  `gorm:"not null"`
}
