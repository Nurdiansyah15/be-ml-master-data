package models

type Tournament struct {
	TournamentID uint   `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"size:100;not null"`
	Season       string `gorm:"size:50"`
}
