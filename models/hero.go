package models

type Hero struct {
	HeroID uint   `gorm:"primaryKey;autoIncrement" json:"hero_id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
}
