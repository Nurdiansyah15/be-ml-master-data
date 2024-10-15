package models

type Player struct {
	PlayerID uint   `gorm:"primaryKey;autoIncrement" json:"player_id"`
	TeamID   uint   `json:"team_id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Image    string `json:"image"`
}
