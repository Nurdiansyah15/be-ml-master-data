package models

type Coach struct {
	CoachID uint   `gorm:"primaryKey;autoIncrement" json:"coach_id"`
	TeamID  uint   `json:"team_id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	Image   string `json:"image"`
}
