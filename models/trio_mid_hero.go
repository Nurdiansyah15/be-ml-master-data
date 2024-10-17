package models

type TrioMidHero struct {
	TrioMidHeroID uint   `gorm:"primaryKey;autoIncrement" json:"trio_mid_hero_id"`
	TrioMidID     uint   `json:"trio_mid_id"`
	HeroID        uint   `json:"hero_id"`
	Role          string `gorm:"type:enum('jungler', 'midlaner', 'roamer')" json:"role"`
	EarlyResult   string `gorm:"type:enum('win', 'draw', 'lose')" json:"early_result"`
}
