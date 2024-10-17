package models

type HeroBan struct {
	HeroBanID         uint `gorm:"primaryKey;autoIncrement" json:"hero_ban_id"`
	MatchTeamDetailID uint `json:"match_team_detail_id"`
	HeroID            uint `json:"hero_id"`
	FirstPhase        int  `json:"first_phase"`
	SecondPhase       int  `json:"second_phase"`
	Total             int  `json:"total"`
}
