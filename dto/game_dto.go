package dto

type GameRequestDto struct {
	FirstPickTeamID  uint   `json:"first_pick_team_id" binding:"required"`
	SecondPickTeamID uint   `json:"second_pick_team_id" binding:"required"`
	WinnerTeamID     uint   `json:"winner_team_id" binding:"required"`
	GameNumber       int    `json:"game_number" binding:"required"`
	VideoLink        string `json:"video_link"`
	FullDraftImage   string `json:"full_draft_image"`
}

type GameResponseDto struct {
	GameID          uint `json:"game_id"`
	MatchID         uint `json:"match_id"`
	FirstPickTeamID uint `json:"first_pick_team_id"`
	FirstTeam       struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:first_team_" json:"first_team"`
	SecondPickTeamID uint `json:"second_pick_team_id"`
	SecondTeam       struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:second_team_" json:"second_team"`
	WinnerTeamID uint `json:"winner_team_id"`
	WinnerTeam   struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:winner_team_" json:"winner_team"`
	GameNumber     int    `json:"game_number"`
	VideoLink      string `json:"video_link"`
	FullDraftImage string `json:"full_draft_image"`
}

type LordResultRequestDto struct {
	TeamID   uint   `json:"team_id" binding:"required"`
	Phase    string `json:"phase" binding:"required"`
	Setup    string `json:"setup" binding:"required"`
	Initiate string `json:"initiate" binding:"required"`
	Result   string `json:"result" binding:"required"`
}

type LordResultResponseDto struct {
	LordResultID uint `json:"lord_result_id"`
	GameID       uint `json:"game_id"`
	TeamID       uint `json:"team_id"`
	Team         struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}
	Phase    string `json:"phase"`
	Setup    string `json:"setup"`
	Initiate string `json:"initiate"`
	Result   string `json:"result"`
}

type TurtleResultRequestDto struct {
	TeamID   uint   `json:"team_id" binding:"required"`
	Phase    string `json:"phase" binding:"required"`
	Setup    string `json:"setup" binding:"required"`
	Initiate string `json:"initiate" binding:"required"`
	Result   string `json:"result" binding:"required"`
}

type TurtleResultResponseDto struct {
	TurtleResultID uint `json:"turtle_result_id"`
	GameID         uint `json:"game_id"`
	TeamID         uint `json:"team_id"`
	Team           struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `json:"team"`
	Phase    string `json:"phase"`
	Setup    string `json:"setup"`
	Initiate string `json:"initiate"`
	Result   string `json:"result"`
}

type ExplanerRequestDto struct {
	HeroID      uint   `json:"hero_id" binding:"required"`
	EarlyResult string `gorm:"type:enum('win', 'draw', 'lose')" json:"early_result"`
}

type ExplanerResponseDto struct {
	ExplanerID uint `json:"explaner_id"`
	GameID     uint `json:"game_id"`
	TeamID     uint `json:"team_id"`
	Team       struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:team_" json:"team"`
	HeroID uint `json:"hero_id"`
	Hero   struct {
		HeroID uint   `json:"hero_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:hero_" json:"hero"`
	EarlyResult string `json:"early_result"`
}

type GoldlanerRequestDto struct {
	HeroID      uint   `json:"hero_id" binding:"required"`
	EarlyResult string `gorm:"type:enum('win', 'draw', 'lose')" json:"early_result"`
}

type GoldlanerResponseDto struct {
	GoldlanerID uint `json:"goldlaner_id"`
	GameID      uint `json:"game_id"`
	TeamID      uint `json:"team_id"`
	Team        struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:team_" json:"team"`
	HeroID uint `json:"hero_id"`
	Hero   struct {
		HeroID uint   `json:"hero_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:hero_" json:"hero"`
	EarlyResult string `json:"early_result"`
}

// khusus
type TrioMidRequestDto struct {
	HeroID      uint   `json:"hero_id" binding:"required"`
	Role        string `json:"role" binding:"required"`
	EarlyResult string `gorm:"type:enum('win', 'draw', 'lose')" json:"early_result"`
}

type TrioMidResponseDto struct {
	TrioMidHeroID uint `json:"trio_mid_hero_id"`
	TrioMidID   uint   `json:"trio_mid_id"`
	GameID      uint   `json:"game_id"`
	Role        string `json:"role"`
	EarlyResult string `json:"early_result"`
	Team        struct {
		TeamID uint   `json:"team_id"` // Tetap menggunakan team_id
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:team_" json:"team"`
	Hero struct {
		HeroID uint   `json:"hero_id"` // Tetap menggunakan hero_id
		Name   string `json:"name"`
		Image  string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:hero_" json:"hero"`
}
