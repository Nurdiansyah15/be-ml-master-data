package models

type Explaner struct {
	ExplanerID  uint   `gorm:"primaryKey;autoIncrement" json:"explaner_id"`
	GameID      uint   `json:"game_id"`
	TeamID      uint   `json:"team_id"`
	HeroID      uint   `json:"hero_id"`
	EarlyResult string `gorm:"type:enum('win', 'draw', 'lose')" json:"early_result"`
}
