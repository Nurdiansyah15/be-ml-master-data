package models

type PriorityBan struct {
	PriorityBanID     uint    `gorm:"primaryKey;autoIncrement" json:"priority_ban_id"`
	MatchTeamDetailID uint    `json:"match_team_detail_id"`
	HeroID            uint    `json:"hero_id"`
	Total             int     `json:"total"`
	Role              string  `gorm:"type:enum('gold', 'exp', 'roam', 'mid', 'jungler');" json:"role"`
	BanRate           float64 `json:"ban_rate"`
}
