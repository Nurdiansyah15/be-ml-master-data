package models

type Game struct {
	GameID           uint   `gorm:"primaryKey;autoIncrement" json:"game_id"`
	MatchID          uint   `json:"match_id"`
	FirstPickTeamID  uint   `json:"first_pick_team_id"`
	SecondPickTeamID uint   `json:"second_pick_team_id"`
	WinnerTeamID     uint   `json:"winner_team_id"`
	GameNumber       int    `json:"game_number"`
	VideoLink        string `json:"video_link"`
	FullDraftImage   string `json:"full_draft_image"`
}
