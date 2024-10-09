package models

type TournamentTeam struct {
	TournamentTeamID uint `gorm:"primaryKey;autoIncrement"`
	TournamentID     uint `gorm:"not null"`
	TeamID           uint `gorm:"not null"`
}
