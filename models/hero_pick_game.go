package models

type HeroPickGame struct {
	HeroPickGameID uint `gorm:"primaryKey;autoIncrement" json:"hero_pick_game_id"`
	HeroPickID     uint `json:"hero_pick_id"`
	GameNumber     int  `json:"game_number"`
	IsPicked       bool `json:"is_picked"`
}
