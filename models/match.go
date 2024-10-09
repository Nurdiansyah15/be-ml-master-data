package models

import "database/sql"

type Match struct {
	MatchID       uint `gorm:"primaryKey;autoIncrement"`
	TournamentID  uint `gorm:"not null"`
	HomeTeamID    uint `gorm:"not null"`
	AwayTeamID    uint
	Week          int
	Day           int
	Date          sql.NullTime `gorm:"type:date"`
	Time          sql.NullTime `gorm:"type:time"`
	HomeTeamScore int
	AwayTeamScore int
	WinnerTeamID  uint
	TotalGames    int
}
