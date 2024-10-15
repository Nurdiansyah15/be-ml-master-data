package models

type Team struct {
	TeamID uint   `gorm:"primaryKey;autoIncrement" json:"team_id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
}
