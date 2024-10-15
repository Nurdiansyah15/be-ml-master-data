package models

type HeroBanGame struct {
	HeroBanGameID uint `gorm:"primaryKey;autoIncrement" json:"hero_ban_game_id"`
	HeroBanID     uint `json:"hero_ban_id"`
	GameNumber    int  `json:"game_number"`
	IsBanned      bool `json:"is_banned"`
}
