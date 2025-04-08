package models

type MatchTeamDetail struct {
	MatchTeamDetailID uint `gorm:"primaryKey;autoIncrement" json:"match_team_detail_id"`
	MatchID           uint `json:"match_id"`
	TeamID            uint `json:"team_id"`
}
