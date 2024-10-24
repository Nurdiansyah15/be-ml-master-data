package dto

type MatchRequestDto struct {
	Stage      *string `json:"stage" binding:"required"`
	Day        *int    `json:"day" binding:"required"`
	Date       *int    `json:"date" binding:"required"`
	TeamAID    *uint   `json:"team_a_id" binding:"required"`
	TeamBID    *uint   `json:"team_b_id" binding:"required"`
	TeamAScore *int    `json:"team_a_score" binding:"required"`
	TeamBScore *int    `json:"team_b_score" binding:"required"`
}

type MatchResponseDto struct {
	MatchID      *uint   `json:"match_id"`
	TournamentID *uint   `json:"tournament_id"`
	Stage        *string `json:"stage"`
	Day          *int    `json:"day"`
	Date         *int    `json:"date"`
	TeamAID      *uint   `json:"team_a_id"`
	TeamA        *struct {
		TeamID *uint   `json:"team_id"`
		Name   *string `json:"name"`
		Image  *string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:team_a_" json:"team_a"`
	TeamBID *uint `json:"team_b_id"`
	TeamB   *struct {
		TeamID *uint   `json:"team_id"`
		Name   *string `json:"name"`
		Image  *string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:team_b_" json:"team_b"`
	TeamAScore *int `json:"team_a_score"`
	TeamBScore *int `json:"team_b_score"`
}

type PlayerMatchRequestDto struct {
	PlayerID *uint   `json:"player_id" binding:"required"`
	Role     *string `json:"role" binding:"required,oneof=goldlaner explaner roamer midlaner jungler"`
}

type UpdatePlayerMatchRequestDto struct {
	Role *string `json:"role" binding:"required,oneof=goldlaner explaner roamer midlaner jungler"`
}

type PlayerMatchResponseDto struct {
	PlayerMatchID     *uint   `json:"player_match_id"`
	MatchTeamDetailID *uint   `json:"match_team_detail_id"`
	Role              *string `json:"role"`
	Player            *struct {
		PlayerID *uint   `json:"player_id"`
		TeamID   *uint   `json:"team_id"`
		Name     *string `json:"name"`
		Image    *string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:player_" json:"player"`
}

type CoachMatchRequestDto struct {
	CoachID *uint   `json:"coach_id" binding:"required"`
	Role    *string `json:"role" binding:"required"`
}

type UpdateCoachMatchRequestDto struct {
	Role *string `json:"role" binding:"required"`
}

type CoachMatchResponseDto struct {
	CoachMatchID      *uint   `json:"coach_match_id"`
	MatchTeamDetailID *uint   `json:"match_team_detail_id"`
	Role              *string `json:"role"`
	Coach             *struct {
		CoachID *uint   `json:"coach_id"`
		TeamID  *uint   `json:"team_id"`
		Name    *string `json:"name"`
		Image   *string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:coach_" json:"coach"`
}

type HeroPickRequestDto struct {
	HeroID       *uint `json:"hero_id" binding:"required"`
	FirstPhase   *int  `json:"first_phase" binding:"required"`
	SecondPhase  *int  `json:"second_phase" binding:"required"`
	Total        *int  `json:"total" binding:"required"`
	HeroPickGame []struct {
		GameID     *uint `json:"game_id" binding:"required"`
		GameNumber *int  `json:"game_number" binding:"required"`
		IsPicked   *bool `json:"is_picked" binding:"required"`
	} `json:"hero_pick_game"`
}

type HeroPickResponseDto struct {
	HeroPickID        *uint `json:"hero_pick_id"`
	MatchTeamDetailID *uint `json:"match_team_detail_id"`
	HeroID            *uint `json:"hero_id"`
	Hero              *struct {
		HeroID *uint   `json:"hero_id"`
		Name   *string `json:"name"`
		Image  *string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:hero_" json:"hero"`
	FirstPhase   *int `json:"first_phase"`
	SecondPhase  *int `json:"second_phase"`
	Total        *int `json:"total"`
	HeroPickGame []struct {
		HeroPickGameID uint `json:"hero_pick_game_id"`
		HeroPickID     uint `json:"hero_pick_id"`
		GameID         uint `json:"game_id"`
		GameNumber     int  `json:"game_number"`
		IsPicked       bool `json:"is_picked"`
	} `json:"hero_pick_game"`
}

type HeroBanRequestDto struct {
	HeroID      *uint `json:"hero_id" binding:"required"`
	FirstPhase  *int  `json:"first_phase" binding:"required"`
	SecondPhase *int  `json:"second_phase" binding:"required"`
	Total       *int  `json:"total" binding:"required"`
	HeroBanGame []struct {
		GameID     *uint `json:"game_id" binding:"required"`
		GameNumber *int  `json:"game_number" binding:"required"`
		IsBanned   *bool `json:"is_banned" binding:"required"`
	} `json:"hero_ban_game"`
}

type HeroBanResponseDto struct {
	HeroBanID         *uint `json:"hero_ban_id"`
	MatchTeamDetailID *uint `json:"match_team_detail_id"`
	HeroID            *uint `json:"hero_id"`
	Hero              *struct {
		HeroID *uint   `json:"hero_id"`
		Name   *string `json:"name"`
		Image  *string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:hero_" json:"hero"`
	FirstPhase  *int `json:"first_phase"`
	SecondPhase *int `json:"second_phase"`
	Total       *int `json:"total"`
	HeroBanGame []struct {
		HeroBanGameID uint `json:"hero_ban_game_id"`
		HeroBanID     uint `json:"hero_ban_id"`
		GameID        uint `json:"game_id"`
		GameNumber    int  `json:"game_number"`
		IsBanned      bool `json:"is_banned"`
	} `json:"hero_ban_game"`
}

type PriorityPickRequestDto struct {
	HeroID   *uint    `json:"hero_id" binding:"required"`
	Total    *int     `json:"total" binding:"required"`
	Role     *string  `json:"role" binding:"required,oneof=gold exp roam mid jungler"`
	PickRate *float64 `json:"pick_rate" binding:"required"`
}

type PriorityPickResponseDto struct {
	PriorityPickID    *uint `json:"priority_pick_id"`
	MatchTeamDetailID *uint `json:"match_team_detail_id"`
	Hero              *struct {
		HeroID *uint   `json:"hero_id"`
		Name   *string `json:"name"`
		Image  *string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:hero_" json:"hero"`
	Total    *int     `json:"total"`
	Role     *string  `json:"role"`
	PickRate *float64 `json:"pick_rate"`
}

type FlexPickRequestDto struct {
	HeroID   *uint    `json:"hero_id" binding:"required"`
	Total    *int     `json:"total" binding:"required"`
	Role     *string  `json:"role" binding:"required,oneof=gold exp roam mid jungler"`
	PickRate *float64 `json:"pick_rate" binding:"required"`
}

type FlexPickResponseDto struct {
	FlexPickID        *uint `json:"flex_pick_id"`
	MatchTeamDetailID *uint `json:"match_team_detail_id"`
	Hero              *struct {
		HeroID *uint   `json:"hero_id"`
		Name   *string `json:"name"`
		Image  *string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:hero_" json:"hero"`
	Total    *int     `json:"total"`
	Role     *string  `json:"role"`
	PickRate *float64 `json:"pick_rate"`
}

type PriorityBanRequestDto struct {
	HeroID  *uint    `json:"hero_id" binding:"required"`
	Total   *int     `json:"total" binding:"required"`
	Role    *string  `json:"role" binding:"required,oneof=gold exp roam mid jungler"`
	BanRate *float64 `json:"ban_rate" binding:"required"`
}

type PriorityBanResponseDto struct {
	PriorityBanID     *uint `json:"priority_ban_id"`
	MatchTeamDetailID *uint `json:"match_team_detail_id"`
	Hero              *struct {
		HeroID *uint   `json:"hero_id"`
		Name   *string `json:"name"`
		Image  *string `json:"image"`
	} `gorm:"embedded;embeddedPrefix:hero_" json:"hero"`
	Total   *int     `json:"total"`
	Role    *string  `json:"role"`
	BanRate *float64 `json:"ban_rate"`
}
