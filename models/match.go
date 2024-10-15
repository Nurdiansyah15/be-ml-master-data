package models

type Match struct {
	MatchID      uint `gorm:"primaryKey;autoIncrement" json:"match_id"`
	TournamentID uint `json:"tournament_id"`
	Week         int  `json:"week"`
	Day          int  `json:"day"`
	Date         int  `json:"date"`
	TeamAID      uint `json:"team_a_id"`
	TeamBID      uint `json:"team_b_id"`
	TeamAScore   int  `json:"team_a_score"`
	TeamBScore   int  `json:"team_b_score"`
}
