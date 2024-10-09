package models

type Game struct {
	GameID             uint   `gorm:"primaryKey;autoIncrement"`
	MatchID            uint   `gorm:"not null"`
	GameNumber         int    `gorm:"not null"`
	FirstPickTeamID    uint   `gorm:"not null"`
	SecondPickTeamID   uint   `gorm:"not null"`
	WinnerTeamID       uint   `gorm:"not null"`
	TrioMidOverallResult string `gorm:"size:20"`
	EarlyGameResult    string `gorm:"size:20"`
}
