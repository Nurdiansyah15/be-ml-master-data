package models

type Coach struct {	
	CoachID uint   `gorm:"primaryKey;autoIncrement" json:"coach_id"`
	TeamID  uint   `json:"team_id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
}
