package models

type Player struct {
	PlayerID uint   `gorm:"primaryKey;autoIncrement" json:"player_id"`
	TeamID   uint   `json:"team_id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
}
