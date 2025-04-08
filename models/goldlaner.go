package models

type Goldlaner struct {
	GoldlanerID uint   `gorm:"primaryKey;autoIncrement" json:"goldlaner_id"`
	GameID      uint   `json:"game_id"`
	TeamID      uint   `json:"team_id"`
	HeroID      uint   `json:"hero_id"`
	EarlyResult string `gorm:"type:enum('win', 'draw', 'lose')" json:"early_result"`
}
