package models

type TurtleResult struct {
	TurtleResultID uint   `gorm:"primaryKey;autoIncrement" json:"turtle_result_id"`
	GameID         uint   `json:"game_id"`
	TeamID         uint   `json:"team_id"`
	Phase          string `json:"phase"`
	Setup          string `gorm:"type:enum('early', 'late', 'no')" json:"setup"`
	Initiate       string `gorm:"type:enum('yes', 'no')" json:"initiate"`
	Result         string `gorm:"type:enum('yes', 'no')" json:"result"`
}
