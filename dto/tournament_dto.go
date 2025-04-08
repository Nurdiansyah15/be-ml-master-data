package dto

type TournamentRequestDto struct {
	Name string `json:"name" binding:"required"`
}
