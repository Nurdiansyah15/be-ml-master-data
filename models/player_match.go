package models

type PlayerMatch struct {
	PlayerMatchID     uint `gorm:"primaryKey;autoIncrement" json:"player_match_id"`
	MatchTeamDetailID uint `json:"match_team_detail_id"`
	PlayerID          uint `json:"player_id"`
}