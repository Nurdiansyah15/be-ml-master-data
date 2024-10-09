package models

type Match struct {
	MatchID          uint `gorm:"primaryKey;autoIncrement"`
	TournamentTeamID uint
	OpponentTeamID   uint
	Week             string
	Day              string
	Date             int
	HomeTeamScore    int
	AwayTeamScore    int
	WinnerTeamID     uint
	TotalGames       int
}
