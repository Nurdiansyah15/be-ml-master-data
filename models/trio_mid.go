package models

type TrioMid struct {
	TrioMidID   uint   `gorm:"primaryKey;autoIncrement" json:"trio_mid_id"`
	GameID      uint   `json:"game_id"`
	TeamID      uint   `json:"team_id"`
	EarlyResult string `gorm:"type:enum('win', 'draw', 'lose')" json:"early_result"`
}
