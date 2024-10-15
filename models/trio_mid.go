package models

type TrioMid struct {
	TrioMidID   uint   `gorm:"primaryKey;autoIncrement" json:"trio_mid_id"`
	GameID      uint   `json:"game_id"`
	TeamID      uint   `json:"team_id"`
	HeroID      uint   `json:"hero_id"`
	Role        string `gorm:"type:enum('jungler', 'midlaner', 'roamer')" json:"role"`
	EarlyResult string `gorm:"type:enum('win', 'draw', 'lose')" json:"early_result"`
}
