package models

type FlexPick struct {
	FlexPickID        uint    `gorm:"primaryKey;autoIncrement" json:"flex_pick_id"`
	MatchTeamDetailID uint    `json:"match_team_detail_id"`
	HeroID            uint    `json:"hero_id"`
	Total             int     `json:"total"`
	Role              string  `gorm:"type:enum('Roam/Exp', 'Jung/Gold', 'Jung/Mid', 'Jung/Exp');" json:"role"`
	PickRate          float64 `json:"pick_rate"`
}
