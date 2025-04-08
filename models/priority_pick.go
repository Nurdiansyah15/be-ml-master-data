package models

type PriorityPick struct {
	PriorityPickID    uint    `gorm:"primaryKey;autoIncrement" json:"priority_pick_id"`
	MatchTeamDetailID uint    `json:"match_team_detail_id"`
	HeroID            uint    `json:"hero_id"`
	Total             int     `json:"total"`
	Role              string  `gorm:"type:enum('gold', 'exp', 'roam', 'mid', 'jungler');" json:"role"`
	PickRate          float64 `json:"pick_rate"`
}
