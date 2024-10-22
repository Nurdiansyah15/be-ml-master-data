package models

type Tournament struct {
	TournamentID uint   `gorm:"primaryKey;autoIncrement" json:"tournament_id"`
	Name         string `gorm:"size:100;" json:"name"`
}
