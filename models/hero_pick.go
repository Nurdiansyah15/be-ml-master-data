package models

type HeroPick struct {
	HeroPickID        uint `gorm:"primaryKey;autoIncrement" json:"hero_pick_id"`
	MatchTeamDetailID uint `json:"match_team_detail_id"`
	HeroID            uint `json:"hero_id"`
	FirstPhase        int  `json:"first_phase"`
	SecondPhase       int  `json:"second_phase"`
}
