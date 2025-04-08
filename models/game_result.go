package models

type GameResult struct {
	GameResultID uint   `gorm:"primaryKey;autoIncrement" json:"game_result_id"`
	GameID       uint   `json:"game_id"`
	TeamID       uint   `json:"team_id"`
	Result       string `gorm:"type:enum('win', 'lose');" json:"result"`
}
