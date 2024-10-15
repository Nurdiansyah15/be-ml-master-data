package models

type CoachMatch struct {
	CoachMatchID      uint `gorm:"primaryKey;autoIncrement" json:"coach_match_id"`
	MatchTeamDetailID uint `json:"match_team_detail_id"`
	CoachID           uint `json:"coach_id"`
}
