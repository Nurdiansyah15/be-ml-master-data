package models

type HeroPick struct {
	PickID           uint `gorm:"primaryKey;autoIncrement"`
	MatchID          uint `gorm:"not null"`
	TeamID           uint `gorm:"not null"`
	HeroID           uint `gorm:"not null"`
	TotalPicks       int  `gorm:"not null"` // Field yang ditambahkan
	FirstPhaseCount  int  `gorm:"not null"` // Field yang ditambahkan
	SecondPhaseCount int  `gorm:"not null"` // Field yang ditambahkan
}
