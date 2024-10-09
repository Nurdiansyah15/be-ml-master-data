package models

type MatchVideo struct {
	VideoID    uint   `gorm:"primaryKey;autoIncrement"`
	MatchID    uint   `gorm:"not null"`
	GameNumber int
	VideoLink  string `gorm:"size:255"`
}
