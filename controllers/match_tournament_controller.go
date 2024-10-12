package controllers

import (
	"ml-master-data/config"
	"ml-master-data/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllTeamMatchesinTournament(c *gin.Context) {

	tournamentID := c.Param("tournamentID")
	teamID := c.Param("teamID")

	tournamentTeam := models.TournamentTeam{}

	if err := config.DB.Where("tournament_id = ? AND team_id = ?", tournamentID, teamID).First(&tournamentTeam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var matches []models.Match
	if err := config.DB.Model(&models.Match{}).Where("tournament_team_id = ?", tournamentTeam.TournamentID).Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matches)
}

func CreateTeamMatchinTournament(c *gin.Context) {
	tournamentID := c.Param("tournamentID")
	if tournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tournament ID is required"})
		return
	}
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	// Check if tournament team exists
	var tournamentTeam models.TournamentTeam
	if err := config.DB.First(&tournamentTeam, "tournament_id = ? AND team_id = ?", tournamentID, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament or team not found"})
		return
	}

	if err := config.DB.Where("tournament_id = ? AND team_id = ?", tournamentID, teamID).First(&tournamentTeam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	input := struct {
		Week string `json:"week" binding:"required"`
		Day  string `json:"day" binding:"required"`
		Date int    `json:"date" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var match models.Match

	match.Week = input.Week
	match.Day = input.Day
	match.Date = input.Date
	match.TournamentTeamID = tournamentTeam.TournamentID

	if err := config.DB.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, match)
}

func UpdateTeamMatchinTournament(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// totalGames := match.TotalGames
	input := struct {
		OpponentTeamID uint   `json:"opponent_team_id"`
		Week           string `json:"week"`
		Day            string `json:"day"`
		Date           int    `json:"date"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if match.AwayTeamScore != 0 || match.HomeTeamScore != 0 {
	// 	match.TotalGames = match.AwayTeamScore + match.HomeTeamScore
	// }

	// if totalGames != match.TotalGames {
	// 	err := config.DB.Delete(&models.Game{}, "match_id = ?", matchID).Error
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// }

	if input.OpponentTeamID != 0 {
		match.OpponentTeamID = input.OpponentTeamID
	}
	if input.Week != "" {
		match.Week = input.Week
	}
	if input.Day != "" {
		match.Day = input.Day
	}
	if input.Date != 0 {
		match.Date = input.Date
	}

	if err := config.DB.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

func GetMatchByID(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	c.JSON(http.StatusOK, match)
}

func AddPlayerStatsToMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	input := struct {
		PlayerID  uint    `json:"player_id" binding:"required"`
		GameRate  float64 `json:"game_rate" binding:"required"`
		MatchRate float64 `json:"match_rate" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var player models.Player
	if err := config.DB.First(&player, input.PlayerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	var playerStats models.PlayerStats

	playerStats.PlayerID = input.PlayerID
	playerStats.GameRate = input.GameRate
	playerStats.MatchRate = input.MatchRate
	playerStats.MatchID = uint(matchIDInt)

	if err := config.DB.Create(&playerStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, playerStats)
}
func UpdatePlayerStats(c *gin.Context) {
	playerStatID := c.Param("playerStatID")
	if playerStatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PlayerStat ID is required"})
		return
	}

	var playerStats models.PlayerStats
	if err := config.DB.First(&playerStats, playerStatID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PlayerStats not found"})
		return
	}

	input := struct {
		PlayerID  uint    `json:"player_id"`
		GameRate  float64 `json:"game_rate"`
		MatchRate float64 `json:"match_rate"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var player models.Player

	if input.PlayerID != 0 {
		if err := config.DB.First(&player, input.PlayerID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			return
		}

		playerStats.PlayerID = input.PlayerID
	}

	if input.GameRate != 0 {
		playerStats.GameRate = input.GameRate
	}

	if input.MatchRate != 0 {
		playerStats.MatchRate = input.MatchRate
	}

	if err := config.DB.Save(&playerStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, playerStats)
}

func GetAllPlayerStatsinMatch(c *gin.Context) {
	matchID := c.Param("matchID")

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var playerStats []models.PlayerStats
	if err := config.DB.Model(&models.PlayerStats{}).Where("match_id = ?", matchIDInt).Find(&playerStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, playerStats)
}

func AddCoachStatsToMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	input := struct {
		CoachID   uint    `json:"coach_id" binding:"required"`
		GameRate  float64 `json:"game_rate" binding:"required"`
		MatchRate float64 `json:"match_rate" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coachStats := models.CoachStats{
		CoachID:   input.CoachID,
		MatchID:   uint(matchIDInt),
		GameRate:  input.GameRate,
		MatchRate: input.MatchRate,
	}

	if err := config.DB.Create(&coachStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, coachStats)
}

func UpdateCoachStats(c *gin.Context) {
	coachStatID := c.Param("coachStatID")
	if coachStatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CoachStat ID is required"})
		return
	}

	var coachStats models.CoachStats
	if err := config.DB.First(&coachStats, coachStatID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CoachStats not found"})
		return
	}

	input := struct {
		CoachID   uint    `json:"coach_id"`
		GameRate  float64 `json:"game_rate"`
		MatchRate float64 `json:"match_rate"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.CoachID != 0 {
		coachStats.CoachID = input.CoachID
	}

	if input.GameRate != 0 {
		coachStats.GameRate = input.GameRate
	}

	if input.MatchRate != 0 {
		coachStats.MatchRate = input.MatchRate
	}

	if err := config.DB.Save(&coachStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coachStats)
}

func GetAllCoachStatsinMatch(c *gin.Context) {
	matchID := c.Param("matchID")

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var coachStats []models.CoachStats
	if err := config.DB.Model(&models.CoachStats{}).Where("match_id = ?", matchIDInt).Find(&coachStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coachStats)
}

func AddFlexPicksToMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	input := struct {
		HeroID   uint    `json:"hero_id" binding:"required"`
		TeamID   uint    `json:"team_id" binding:"required"`
		Total    int     `json:"total" binding:"required"`
		Role     string  `json:"role" binding:"required,max=50"` // Role menjadi required
		RatePick float64 `json:"rate_pick" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hero models.Hero
	if err := config.DB.First(&hero, input.HeroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	var team models.Team
	if err := config.DB.First(&team, input.TeamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	flexPick := models.FlexPick{
		HeroID:   input.HeroID,
		MatchID:  uint(matchIDInt),
		TeamID:   input.TeamID,
		Total:    input.Total,
		Role:     input.Role,
		RatePick: input.RatePick,
	}

	if err := config.DB.Create(&flexPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, flexPick)
}

func UpdateFlexPick(c *gin.Context) {
	flexPickID := c.Param("flexPickID")
	if flexPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FlexPick ID is required"})
		return
	}

	var flexPick models.FlexPick
	if err := config.DB.First(&flexPick, flexPickID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FlexPick not found"})
		return
	}

	input := struct {
		HeroID   uint    `json:"hero_id"`
		TeamID   uint    `json:"team_id"`
		Total    int     `json:"total"`
		Role     string  `json:"role" binding:"max=50"`
		RatePick float64 `json:"rate_pick"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.HeroID != 0 {
		var hero models.Hero
		if err := config.DB.First(&hero, input.HeroID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
			return
		}
		flexPick.HeroID = input.HeroID
	}

	if input.TeamID != 0 {
		var team models.Team
		if err := config.DB.First(&team, input.TeamID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
			return
		}
		flexPick.TeamID = input.TeamID
	}

	if input.Total != 0 {
		flexPick.Total = input.Total
	}

	if input.Role != "" {
		flexPick.Role = input.Role
	}

	if input.RatePick != 0 {
		flexPick.RatePick = input.RatePick
	}

	if err := config.DB.Save(&flexPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flexPick)
}

func GetAllFlexPicksinMatch(c *gin.Context) {
	matchID := c.Param("matchID")

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var flexPicks []models.FlexPick
	if err := config.DB.Model(&models.FlexPick{}).Where("match_id = ?", matchIDInt).Find(&flexPicks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flexPicks)
}

func AddPriorityBansToMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchIDInt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	input := struct {
		HeroID  uint    `json:"hero_id" binding:"required"`
		TeamID  uint    `json:"team_id" binding:"required"`
		Total   int     `json:"total" binding:"required"`
		Role    string  `json:"role" binding:"required"`
		RateBan float64 `json:"rate_ban" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priorityBan := models.PriorityBan{
		HeroID:  input.HeroID,
		MatchID: uint(matchIDInt),
		TeamID:  input.TeamID,
		Total:   input.Total,
		Role:    input.Role,
		RateBan: input.RateBan,
	}

	// Validasi Hero
	var hero models.Hero
	if err := config.DB.First(&hero, priorityBan.HeroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Validasi Team
	var team models.Team
	if err := config.DB.First(&team, priorityBan.TeamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	if err := config.DB.Create(&priorityBan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, priorityBan)
}

func UpdatePriorityBan(c *gin.Context) {
	priorityBanID := c.Param("priorityBanID")
	if priorityBanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PriorityBan ID is required"})
		return
	}

	var priorityBan models.PriorityBan
	if err := config.DB.First(&priorityBan, priorityBanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PriorityBan not found"})
		return
	}

	input := struct {
		HeroID  *uint    `json:"hero_id"`
		TeamID  *uint    `json:"team_id"`
		Total   *int     `json:"total"`
		Role    *string  `json:"role" binding:"max=50"`
		RateBan *float64 `json:"rate_ban"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.HeroID != nil {
		var hero models.Hero
		if err := config.DB.First(&hero, *input.HeroID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
			return
		}
		priorityBan.HeroID = *input.HeroID
	}

	if input.TeamID != nil {
		var team models.Team
		if err := config.DB.First(&team, *input.TeamID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
			return
		}
		priorityBan.TeamID = *input.TeamID
	}

	if input.Total != nil {
		priorityBan.Total = *input.Total
	}

	if input.Role != nil {
		priorityBan.Role = *input.Role
	}

	if input.RateBan != nil {
		priorityBan.RateBan = *input.RateBan
	}

	if err := config.DB.Save(&priorityBan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, priorityBan)
}

func GetAllPriorityBansinMatch(c *gin.Context) {
	matchID := c.Param("matchID")

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var priorityBans []models.PriorityBan
	if err := config.DB.Model(&models.PriorityBan{}).Where("match_id = ?", matchIDInt).Find(&priorityBans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, priorityBans)
}

func AddPriorityPickToMatch(c *gin.Context) {
	matchID := c.Param("matchID")

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchIDInt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Struktur untuk input dari JSON
	input := struct {
		HeroID   uint    `json:"hero_id" binding:"required"`
		TeamID   uint    `json:"team_id" binding:"required"`
		Total    int     `json:"total" binding:"required"`
		Role     string  `json:"role" binding:"max=50"`
		RatePick float64 `json:"rate_pick" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priorityPick := models.PriorityPick{
		HeroID:   input.HeroID,
		MatchID:  uint(matchIDInt),
		TeamID:   input.TeamID,
		Total:    input.Total,
		Role:     input.Role,
		RatePick: input.RatePick,
	}

	// Validasi Hero
	var hero models.Hero
	if err := config.DB.First(&hero, priorityPick.HeroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Validasi Team
	var team models.Team
	if err := config.DB.First(&team, priorityPick.TeamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	if err := config.DB.Create(&priorityPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, priorityPick)
}

func UpdatePriorityPick(c *gin.Context) {
	priorityPickID := c.Param("priorityPickID")
	if priorityPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PriorityPick ID is required"})
		return
	}

	var priorityPick models.PriorityPick
	if err := config.DB.First(&priorityPick, priorityPickID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PriorityPick not found"})
		return
	}

	// Struktur untuk input dari JSON
	input := struct {
		HeroID   *uint    `json:"hero_id"` // Pointer untuk nullable fields
		TeamID   *uint    `json:"team_id"` // Pointer untuk nullable fields
		Total    *int     `json:"total"`   // Pointer untuk nullable fields
		Role     *string  `json:"role" binding:"max=50"`
		RatePick *float64 `json:"rate_pick"` // Pointer untuk nullable fields
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if they are provided
	if input.HeroID != nil {
		var hero models.Hero
		if err := config.DB.First(&hero, *input.HeroID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
			return
		}
		priorityPick.HeroID = *input.HeroID
	}

	if input.TeamID != nil {
		var team models.Team
		if err := config.DB.First(&team, *input.TeamID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
			return
		}
		priorityPick.TeamID = *input.TeamID
	}

	if input.Total != nil {
		priorityPick.Total = *input.Total
	}

	if input.Role != nil {
		priorityPick.Role = *input.Role
	}

	if input.RatePick != nil {
		priorityPick.RatePick = *input.RatePick
	}

	if err := config.DB.Save(&priorityPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, priorityPick)
}

func GetAllPriorityPicksinMatch(c *gin.Context) {
	matchID := c.Param("matchID")

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var priorityPicks []models.PriorityPick
	if err := config.DB.Model(&models.PriorityPick{}).Where("match_id = ?", matchIDInt).Find(&priorityPicks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, priorityPicks)
}
