package controllers

import (
	"ml-master-data/config"
	"ml-master-data/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTournamentMatch(c *gin.Context) {
	tournamentID := c.Param("tournamentID")
	if tournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tournament ID is required"})
		return
	}

	var tournament models.Tournament
	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	input := struct {
		Week       int  `json:"week" binding:"required"`
		Day        int  `json:"day" binding:"required"`
		Date       int  `json:"date" binding:"required"`
		TeamAID    uint `json:"team_a_id" binding:"required"`
		TeamBID    uint `json:"team_b_id" binding:"required"`
		TeamAScore int  `json:"team_a_score" binding:"required"`
		TeamBScore int  `json:"team_b_score" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah Team A dan Team B valid
	var teamA, teamB models.Team
	if err := config.DB.First(&teamA, input.TeamAID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team A not found"})
		return
	}
	if err := config.DB.First(&teamB, input.TeamBID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team B not found"})
		return
	}

	tx := config.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}()

	match := models.Match{
		TournamentID: tournament.TournamentID,
		Week:         input.Week,
		Day:          input.Day,
		Date:         input.Date,
		TeamAID:      input.TeamAID,
		TeamBID:      input.TeamBID,
		TeamAScore:   input.TeamAScore,
		TeamBScore:   input.TeamBScore,
	}

	if err := tx.Create(&match).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	matchTeamADetails := models.MatchTeamDetail{
		MatchID: match.MatchID,
		TeamID:  teamA.TeamID,
	}

	if err := tx.Create(&matchTeamADetails).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	matchTeamBDetails := models.MatchTeamDetail{
		MatchID: match.MatchID,
		TeamID:  teamB.TeamID,
	}

	if err := tx.Create(&matchTeamBDetails).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, match)
}

// UpdateMatch untuk memperbarui informasi pertandingan
func UpdateMatch(c *gin.Context) {
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

	input := struct {
		Week       int  `json:"week"`
		Day        int  `json:"day"`
		Date       int  `json:"date"`
		TeamAID    uint `json:"team_a_id"`
		TeamBID    uint `json:"team_b_id"`
		TeamAScore int  `json:"team_a_score"`
		TeamBScore int  `json:"team_b_score"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah Team A dan Team B valid
	if input.TeamAID != 0 {
		var teamA models.Team
		if err := config.DB.First(&teamA, input.TeamAID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team A not found"})
			return
		}
		match.TeamAID = input.TeamAID
	}

	if input.TeamBID != 0 {
		var teamB models.Team
		if err := config.DB.First(&teamB, input.TeamBID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team B not found"})
			return
		}
		match.TeamBID = input.TeamBID
	}

	if input.Week != 0 {
		match.Week = input.Week
	}
	if input.Day != 0 {
		match.Day = input.Day
	}
	if input.Date != 0 {
		match.Date = input.Date
	}
	if input.TeamAScore >= 0 {
		match.TeamAScore = input.TeamAScore
	}
	if input.TeamBScore >= 0 {
		match.TeamBScore = input.TeamBScore
	}

	if err := config.DB.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

// GetMatchByID untuk mendapatkan informasi pertandingan berdasarkan ID
func GetMatchByID(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	type Team struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type Match struct {
		MatchID    uint `json:"match_id"`
		Week       int  `json:"week"`
		Day        int  `json:"day"`
		Date       int  `json:"date"`
		TeamAID    uint `json:"team_a_id"`
		TeamA      Team `gorm:"embedded;embeddedPrefix:team_a_" json:"team_a"`
		TeamBID    uint `json:"team_b_id"`
		TeamB      Team `gorm:"embedded;embeddedPrefix:team_b_" json:"team_b"`
		TeamAScore int  `json:"team_a_score"`
		TeamBScore int  `json:"team_b_score"`
	}

	var match Match

	query := `
		SELECT 
			m.match_id, m.week, m.day, m.date, m.team_a_id, m.team_b_id, 
			tA.team_id AS team_a_id, tA.name AS team_a_name, tA.image AS team_a_image,
			tB.team_id AS team_b_id, tB.name AS team_b_name, tB.image AS team_b_image,
			m.team_a_score, m.team_b_score
		FROM matches m
		JOIN teams tA ON m.team_a_id = tA.team_id
		JOIN teams tB ON m.team_b_id = tB.team_id
		WHERE m.match_id = ?
	`

	if err := config.DB.Raw(query, matchID).Scan(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	c.JSON(http.StatusOK, match)
}

// GetAllMatches untuk mendapatkan semua pertandingan
func GetMatchesByTournamentID(c *gin.Context) {

	tournamentID := c.Param("tournamentID")
	if tournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tournament ID is required"})
		return
	}

	// Cek apakah tournament valid
	var tournament models.Tournament
	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	type Team struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type Match struct {
		MatchID    uint `json:"match_id"`
		Week       int  `json:"week"`
		Day        int  `json:"day"`
		Date       int  `json:"date"`
		TeamAID    uint `json:"team_a_id"`
		TeamA      Team `gorm:"embedded;embeddedPrefix:team_a_" json:"team_a"`
		TeamBID    uint `json:"team_b_id"`
		TeamB      Team `gorm:"embedded;embeddedPrefix:team_b_" json:"team_b"`
		TeamAScore int  `json:"team_a_score"`
		TeamBScore int  `json:"team_b_score"`
	}

	var matches []Match

	query := `
		SELECT 
			m.match_id, m.week, m.day, m.date, m.team_a_id, m.team_b_id, 
			tA.team_id AS team_a_id, tA.name AS team_a_name, tA.image AS team_a_image,
			tB.team_id AS team_b_id, tB.name AS team_b_name, tB.image AS team_b_image,
			m.team_a_score, m.team_b_score
		FROM matches m
		JOIN teams tA ON m.team_a_id = tA.team_id
		JOIN teams tB ON m.team_b_id = tB.team_id
		WHERE m.tournament_id = ?
	`

	if err := config.DB.Raw(query, tournamentID).Scan(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matches)
}

func AddPlayerMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	matchTeamDetail := models.MatchTeamDetail{}

	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	input := struct {
		PlayerID uint `json:"player_id" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player := models.Player{}

	if err := config.DB.First(&player, input.PlayerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	playerMatch := models.PlayerMatch{
		MatchTeamDetailID: matchTeamDetail.MatchTeamDetailID,
		PlayerID:          player.PlayerID,
	}

	if err := config.DB.Create(&playerMatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Player match added successfully"})

}

func RemovePlayerMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	playerID := c.Param("playerID")
	if matchID == "" || teamID == "" || playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID and Player ID are required"})
		return
	}

	matchTeamDetail := models.MatchTeamDetail{}
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	playerMatch := models.PlayerMatch{}
	if err := config.DB.Where("match_team_detail_id = ? AND player_id = ?", matchTeamDetail.MatchTeamDetailID, playerID).First(&playerMatch).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player match not found"})
		return
	}

	if err := config.DB.Delete(&playerMatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player match removed successfully"})

}

func GetAllPlayersMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	matchTeamDetail := models.MatchTeamDetail{}

	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	type PlayerMatch struct {
		PlayerMatchID     uint `json:"player_match_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		PlayerID          uint `json:"player_id"`
		Player            struct {
			PlayerID uint   `json:"player_id"`
			TeamID   uint   `json:"team_id"`
			Name     string `json:"name"`
			Role     string `json:"role"`
			Image    string `json:"image"`
		} `json:"player"`
	}

	var players []PlayerMatch

	query := `
		SELECT 
			pm.player_match_id, pm.match_team_detail_id, pm.player_id,
			p.player_id AS player_player_id, p.team_id AS player_team_id, 
			p.name AS player_name, p.role AS player_role, p.image AS player_image
		FROM player_matches pm
		JOIN players p ON pm.player_id = p.player_id
		JOIN match_team_details mtd ON pm.match_team_detail_id = mtd.match_team_detail_id
		WHERE mtd.match_id = ? AND mtd.team_id = ?
	`

	// Eksekusi query
	if err := config.DB.Raw(query, matchID, teamID).Scan(&players).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Players not found"})
		return
	}

	c.JSON(http.StatusOK, players)

}

func AddCoachMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	input := struct {
		CoachID uint `json:"coach_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var coach models.Coach
	if err := config.DB.First(&coach, input.CoachID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	coachMatch := models.CoachMatch{
		MatchTeamDetailID: matchTeamDetail.MatchTeamDetailID,
		CoachID:           coach.CoachID,
	}

	if err := config.DB.Create(&coachMatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Coach match added successfully"})

}

func RemoveCoachMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	coachID := c.Param("coachID")

	if matchID == "" || teamID == "" || coachID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Coach ID are required"})
		return
	}

	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	var coachMatch models.CoachMatch
	if err := config.DB.Where("match_team_detail_id = ? AND coach_id = ?", matchTeamDetail.MatchTeamDetailID, coachID).First(&coachMatch).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach match not found"})
		return
	}

	if err := config.DB.Delete(&coachMatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coach match removed successfully"})

}

func GetAllCoachesMatch(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	type CoachMatchResponse struct {
		CoachMatchID      uint `json:"coach_match_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		CoachID           uint `json:"coach_id"`
		Coach             struct {
			CoachID uint   `json:"coach_id"`
			TeamID  uint   `json:"team_id"`
			Name    string `json:"name"`
			Role    string `json:"role"`
			Image   string `json:"image"`
		} `json:"coach"`
	}

	var coaches []CoachMatchResponse

	query := `
		SELECT 
			cm.coach_match_id, cm.match_team_detail_id, cm.coach_id,
			c.coach_id AS coach_coach_id, c.team_id AS coach_team_id, 
			c.name AS coach_name, c.role AS coach_role, c.image AS coach_image
		FROM coach_matches cm
		JOIN coaches c ON cm.coach_id = c.coach_id
		JOIN match_team_details mtd ON cm.match_team_detail_id = mtd.match_team_detail_id
		WHERE mtd.match_id = ? AND mtd.team_id = ?
	`

	if err := config.DB.Raw(query, matchID, teamID).Scan(&coaches).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coaches not found"})
		return
	}

	c.JSON(http.StatusOK, coaches)
}

func AddHeroPick(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	var input struct {
		HeroID       uint `json:"hero_id" binding:"required"`
		FirstPhase   int  `json:"first_phase" binding:"required"`
		SecondPhase  int  `json:"second_phase" binding:"required"`
		HeroPickGame []struct {
			GameNumber int  `json:"game_number" binding:"required"`
			IsPicked   bool `json:"is_picked" binding:"required"`
		}
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}()

	heroPick := models.HeroPick{
		MatchTeamDetailID: matchTeamDetail.MatchTeamDetailID,
		HeroID:            input.HeroID,
		FirstPhase:        input.FirstPhase,
		SecondPhase:       input.SecondPhase,
	}
	if err := tx.Create(&heroPick).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, game := range input.HeroPickGame {
		heroPickGame := models.HeroPickGame{
			HeroPickID: heroPick.HeroPickID,
			GameNumber: game.GameNumber,
			IsPicked:   game.IsPicked,
		}
		if err := tx.Create(&heroPickGame).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Hero pick added successfully"})

}

func UpdateHeroPick(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	heroPickID := c.Param("heroPickID")

	if matchID == "" || teamID == "" || heroPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	var heroPick models.HeroPick
	if err := config.DB.Where("match_team_detail_id = ? AND hero_pick_id = ?", matchTeamDetail.MatchTeamDetailID, heroPickID).First(&heroPick).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero pick not found"})
		return
	}

	var input struct {
		HeroID       uint `json:"hero_id" binding:"required"`
		FirstPhase   int  `json:"first_phase" binding:"required"`
		SecondPhase  int  `json:"second_phase" binding:"required"`
		HeroPickGame []struct {
			GameNumber int  `json:"game_number" binding:"required"`
			IsPicked   bool `json:"is_picked" binding:"required"`
		}
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}()

	heroPick.HeroID = input.HeroID
	heroPick.FirstPhase = input.FirstPhase
	heroPick.SecondPhase = input.SecondPhase

	if err := tx.Save(&heroPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, game := range input.HeroPickGame {
		var heroPickGame models.HeroPickGame
		if err := tx.Where("hero_pick_id = ? AND game_number = ?", heroPick.HeroPickID, game.GameNumber).First(&heroPickGame).Error; err != nil {
			heroPickGame = models.HeroPickGame{
				HeroPickID: heroPick.HeroPickID,
				GameNumber: game.GameNumber,
				IsPicked:   game.IsPicked,
			}
			if err := tx.Create(&heroPickGame).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			heroPickGame.HeroPickID = heroPick.HeroPickID
			heroPickGame.GameNumber = game.GameNumber
			heroPickGame.IsPicked = game.IsPicked
			if err := tx.Save(&heroPickGame).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Hero pick updated successfully"})

}

func RemoveHeroPick(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	heroPickID := c.Param("heroPickID")

	if matchID == "" || teamID == "" || heroPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	var heroPick models.HeroPick
	if err := config.DB.Where("match_team_detail_id = ? AND hero_pick_id = ?", matchTeamDetail.MatchTeamDetailID, heroPickID).First(&heroPick).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero pick not found"})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}()

	heroPicksGame := []models.HeroPickGame{}
	if err := tx.Where("hero_pick_id = ?", heroPick.HeroPickID).Find(&heroPicksGame).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, heroPickGame := range heroPicksGame {
		if err := tx.Delete(&heroPickGame).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Delete(&heroPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hero pick removed successfully"})

}
func GetAllHeroPicks(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	type HeroPickResponse struct {
		HeroPickID        uint `json:"hero_pick_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		HeroID            uint `json:"hero_id"`
		Hero              struct {
			HeroID uint   `json:"hero_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		FirstPhase  int `json:"first_phase"`
		SecondPhase int `json:"second_phase"`
		Total       int `json:"total"`
	}

	var picks []HeroPickResponse
	query := `
		SELECT 
			hp.hero_pick_id, hp.match_team_detail_id, hp.hero_id, 
			hp.first_phase, hp.second_phase, hp.total, 
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image
		FROM hero_picks hp
		JOIN heroes h ON hp.hero_id = h.hero_id
		JOIN match_team_details mtd ON hp.match_team_detail_id = mtd.match_team_detail_id
		WHERE mtd.match_id = ? AND mtd.team_id = ?
	`

	if err := config.DB.Raw(query, matchID, teamID).Scan(&picks).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero picks not found"})
		return
	}

	c.JSON(http.StatusOK, picks)
}

func AddHeroBan(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	var input struct {
		HeroID      uint `json:"hero_id" binding:"required"`
		FirstPhase  int  `json:"first_phase" binding:"required"`
		SecondPhase int  `json:"second_phase" binding:"required"`
		Total       int  `json:"total" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	heroBan := models.HeroBan{
		MatchTeamDetailID: matchTeamDetail.MatchTeamDetailID,
		HeroID:            input.HeroID,
		FirstPhase:        input.FirstPhase,
		SecondPhase:       input.SecondPhase,
		Total:             input.Total,
	}

	if err := config.DB.Create(&heroBan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Hero ban added successfully"})

}

func UpdateHeroBan(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	HeroBanID := c.Param("HeroBanID")

	if matchID == "" || teamID == "" || HeroBanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	var heroBan models.HeroBan
	if err := config.DB.Where("match_team_detail_id = ? AND hero_ban_id = ?", matchTeamDetail.MatchTeamDetailID, HeroBanID).First(&heroBan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero ban not found"})
		return
	}

	var input struct {
		HeroID      uint `json:"hero_id" binding:"required"`
		FirstPhase  int  `json:"first_phase" binding:"required"`
		SecondPhase int  `json:"second_phase" binding:"required"`
		HeroBanGame []struct {
			GameNumber int  `json:"game_number" binding:"required"`
			IsBanned   bool `json:"is_banned" binding:"required"`
		}
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}()

	heroBan.HeroID = input.HeroID
	heroBan.FirstPhase = input.FirstPhase
	heroBan.SecondPhase = input.SecondPhase

	if err := tx.Save(&heroBan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, game := range input.HeroBanGame {
		var heroBanGame models.HeroBanGame
		if err := tx.Where("hero_pick_id = ? AND game_number = ?", heroBan.HeroBanID, game.GameNumber).First(&heroBanGame).Error; err != nil {
			heroBanGame = models.HeroBanGame{
				HeroBanID:  heroBan.HeroBanID,
				GameNumber: game.GameNumber,
				IsBanned:   game.IsBanned,
			}
			if err := tx.Create(&heroBanGame).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			heroBanGame.HeroBanID = heroBan.HeroBanID
			heroBanGame.GameNumber = game.GameNumber
			heroBanGame.IsBanned = game.IsBanned
			if err := tx.Save(&heroBanGame).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Hero ban updated successfully"})

}
func RemoveHeroBan(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	HeroBanID := c.Param("HeroBanID")

	if matchID == "" || teamID == "" || HeroBanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	var heroBan models.HeroBan
	if err := config.DB.Where("match_team_detail_id = ? AND hero_ban_id = ?", matchTeamDetail.MatchTeamDetailID, HeroBanID).First(&heroBan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero pick not found"})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}()

	heroBansGame := []models.HeroBanGame{}
	if err := tx.Where("hero_ban_id = ?", heroBan.HeroBanID).Find(&heroBansGame).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, heroBanGame := range heroBansGame {
		if err := tx.Delete(&heroBanGame).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Delete(&heroBan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hero ban removed successfully"})
}
func GetAllHeroBans(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	type HeroBanResponse struct {
		HeroBanID         uint `json:"hero_ban_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		HeroID            uint `json:"hero_id"`
		Hero              struct {
			HeroID uint   `json:"hero_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		FirstPhase  int `json:"first_phase"`
		SecondPhase int `json:"second_phase"`
		Total       int `json:"total"`
	}

	var bans []HeroBanResponse
	query := `
		SELECT 
			hb.hero_ban_id, hb.match_team_detail_id, hb.hero_id, 
			hb.first_phase, hb.second_phase, hb.total, 
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image
		FROM hero_bans hb
		JOIN heroes h ON hb.hero_id = h.hero_id
		JOIN match_team_details mtd ON hb.match_team_detail_id = mtd.match_team_detail_id
		WHERE mtd.match_id = ? AND mtd.team_id = ?
	`

	if err := config.DB.Raw(query, matchID, teamID).Scan(&bans).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero bans not found"})
		return
	}

	c.JSON(http.StatusOK, bans)
}

func AddPriorityPick(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	// Validasi ID
	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Struct untuk menerima input JSON
	input := struct {
		HeroID   uint    `json:"hero_id" binding:"required"`
		Total    int     `json:"total" binding:"required"`
		Role     string  `json:"role" binding:"required,oneof=Gold Exp Roam Mid Jung"`
		PickRate float64 `json:"pick_rate" binding:"required"`
	}{}

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi keberadaan Hero
	var hero models.Hero
	if err := config.DB.First(&hero, input.HeroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Buat PriorityPick baru
	priorityPick := models.PriorityPick{
		MatchTeamDetailID: matchTeamDetail.MatchTeamDetailID,
		HeroID:            input.HeroID,
		Total:             input.Total,
		Role:              input.Role,
		PickRate:          input.PickRate,
	}

	// Simpan ke database
	if err := config.DB.Create(&priorityPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Priority pick added successfully"})
}

func UpdatePriorityPick(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	priorityPickID := c.Param("priorityPickID")

	// Validasi ID
	if matchID == "" || teamID == "" || priorityPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Priority Pick ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Cek keberadaan PriorityPick
	var priorityPick models.PriorityPick
	if err := config.DB.Where("priority_pick_id = ? AND match_team_detail_id = ?", priorityPickID, matchTeamDetail.MatchTeamDetailID).
		First(&priorityPick).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Priority pick not found"})
		return
	}
	// Struct untuk menerima input JSON
	input := struct {
		HeroID   uint    `json:"hero_id" binding:"required"`
		Total    int     `json:"total" binding:"required"`
		Role     string  `json:"role" binding:"required,oneof=Gold Exp Roam Mid Jung"`
		PickRate float64 `json:"pick_rate" binding:"required"`
	}{}

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi keberadaan Hero
	var hero models.Hero
	if err := config.DB.First(&hero, input.HeroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Perbarui data PriorityPick
	priorityPick.HeroID = input.HeroID
	priorityPick.Total = input.Total
	priorityPick.Role = input.Role
	priorityPick.PickRate = input.PickRate

	// Simpan perubahan ke database
	if err := config.DB.Save(&priorityPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Priority pick updated successfully"})
}

func GetAllPriorityPicks(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	// Validasi ID
	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	type PriorityPickResponse struct {
		PriorityPickID    uint `json:"priority_pick_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		HeroID            uint `json:"hero_id"`
		Hero              struct {
			HeroID uint   `json:"hero_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		Total    int     `json:"total"`
		Role     string  `json:"role"`
		PickRate float64 `json:"pick_rate"`
	}

	var priorityPicks []PriorityPickResponse

	// Query dengan WHERE untuk filter MatchTeamDetailID
	query := `
		SELECT 
			pp.priority_pick_id, pp.match_team_detail_id, pp.hero_id, 
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image, 
			pp.total, pp.role, pp.pick_rate
		FROM priority_picks pp
		JOIN heroes h ON pp.hero_id = h.hero_id
		WHERE pp.match_team_detail_id = ?
	`

	if err := config.DB.Raw(query, matchTeamDetail.MatchTeamDetailID).Scan(&priorityPicks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve priority picks"})
		return
	}

	c.JSON(http.StatusOK, priorityPicks)
}

func GetPriorityPickByID(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	priorityPickID := c.Param("priorityPickID")

	// Validasi ID
	if matchID == "" || teamID == "" || priorityPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	type PriorityPickResponse struct {
		PriorityPickID    uint `json:"priority_pick_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		HeroID            uint `json:"hero_id"`
		Hero              struct {
			HeroID uint   `json:"hero_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		Total    int     `json:"total"`
		Role     string  `json:"role"`
		PickRate float64 `json:"pick_rate"`
	}

	var priorityPick PriorityPickResponse

	query := `
		SELECT 
			pp.priority_pick_id, pp.match_team_detail_id, pp.hero_id, 
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image, 
			pp.total, pp.role, pp.pick_rate
		FROM priority_picks pp
		JOIN heroes h ON pp.hero_id = h.hero_id
		JOIN match_team_details mtd ON pp.match_team_detail_id = mtd.match_team_detail_id
		WHERE mtd.match_id = ? AND mtd.team_id = ? AND pp.priority_pick_id = ?
	`

	if err := config.DB.Raw(query, matchID, teamID, priorityPickID).Scan(&priorityPick).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Priority pick not found"})
		return
	}

	c.JSON(http.StatusOK, priorityPick)
}

func RemovePriorityPick(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	priorityPickID := c.Param("priorityPickID")

	// Validasi ID
	if matchID == "" || teamID == "" || priorityPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Priority Pick ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Cek keberadaan PriorityPick
	var priorityPick models.PriorityPick
	if err := config.DB.Where("priority_pick_id = ? AND match_team_detail_id = ?", priorityPickID, matchTeamDetail.MatchTeamDetailID).
		First(&priorityPick).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Priority pick not found"})
		return
	}

	// Hapus PriorityPick
	if err := config.DB.Delete(&priorityPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete priority pick"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Priority pick deleted successfully"})
}

func AddFlexPick(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	// Validasi ID
	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Struct untuk menerima input JSON
	input := struct {
		HeroID   uint    `json:"hero_id" binding:"required"`
		Total    int     `json:"total" binding:"required"`
		Role     string  `json:"role" binding:"required,oneof=Roam/Exp Jung/Gold Jung/Mid Jung/Exp"`
		PickRate float64 `json:"pick_rate" binding:"required"`
	}{}

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi keberadaan Hero
	var hero models.Hero
	if err := config.DB.First(&hero, input.HeroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Buat FlexPick baru
	flexPick := models.FlexPick{
		MatchTeamDetailID: matchTeamDetail.MatchTeamDetailID,
		HeroID:            input.HeroID,
		Total:             input.Total,
		Role:              input.Role,
		PickRate:          input.PickRate,
	}

	// Simpan ke database
	if err := config.DB.Create(&flexPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Flex pick added successfully"})
}

func UpdateFlexPick(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	flexPickID := c.Param("flexPickID")

	// Validasi ID
	if matchID == "" || teamID == "" || flexPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Flex Pick ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Cek keberadaan PriorityPick
	var flexPick models.FlexPick
	if err := config.DB.Where("flex_pick_id = ? AND match_team_detail_id = ?", flexPickID, matchTeamDetail.MatchTeamDetailID).
		First(&flexPick).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flex pick not found"})
		return
	}

	// Struct untuk menerima input JSON
	input := struct {
		HeroID   uint    `json:"hero_id" binding:"required"`
		Total    int     `json:"total" binding:"required"`
		Role     string  `json:"role" binding:"required,oneof=Roam/Exp Jung/Gold Jung/Mid Jung/Exp"`
		PickRate float64 `json:"pick_rate" binding:"required"`
	}{}

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi keberadaan Hero
	var hero models.Hero
	if err := config.DB.First(&hero, input.HeroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Perbarui data FlexPick
	flexPick.HeroID = input.HeroID
	flexPick.Total = input.Total
	flexPick.Role = input.Role
	flexPick.PickRate = input.PickRate

	// Simpan perubahan ke database
	if err := config.DB.Save(&flexPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flex pick updated successfully"})
}

func GetAllFlexPicks(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	// Validasi ID
	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	type FlexPickResponse struct {
		FlexPickID        uint `json:"flex_pick_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		HeroID            uint `json:"hero_id"`
		Hero              struct {
			HeroID uint   `json:"hero_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		Total    int     `json:"total"`
		Role     string  `json:"role"`
		PickRate float64 `json:"pick_rate"`
	}

	var flexPicks []FlexPickResponse

	// Query dengan WHERE untuk filter MatchTeamDetailID
	query := `
		SELECT 
			fp.flex_pick_id, fp.match_team_detail_id, fp.hero_id, 
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image, 
			fp.total, fp.role, fp.pick_rate
		FROM flex_picks fp
		JOIN heroes h ON fp.hero_id = h.hero_id
		WHERE fp.match_team_detail_id = ?
	`

	if err := config.DB.Raw(query, matchTeamDetail.MatchTeamDetailID).Scan(&flexPicks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve flex picks"})
		return
	}

	c.JSON(http.StatusOK, flexPicks)
}

func GetFlexPickByID(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	flexPickID := c.Param("flexPickID")

	// Validasi ID
	if matchID == "" || teamID == "" || flexPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	type FlexPickResponse struct {
		FlexPickID        uint `json:"flex_pick_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		HeroID            uint `json:"hero_id"`
		Hero              struct {
			HeroID uint   `json:"hero_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		Total    int     `json:"total"`
		Role     string  `json:"role"`
		PickRate float64 `json:"pick_rate"`
	}

	var flexPick FlexPickResponse

	query := `
		SELECT 
			fp.flex_pick_id, fp.match_team_detail_id, fp.hero_id, 
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image, 
			fp.total, fp.role, fp.pick_rate
		FROM flex_picks fp
		JOIN heroes h ON fp.hero_id = h.hero_id
		JOIN match_team_details mtd ON fp.match_team_detail_id = mtd.match_team_detail_id
		WHERE mtd.match_id = ? AND mtd.team_id = ? AND fp.flex_pick_id = ?
	`

	if err := config.DB.Raw(query, matchID, teamID, flexPickID).Scan(&flexPick).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flex pick not found"})
		return
	}

	c.JSON(http.StatusOK, flexPick)
}

func DeleteFlexPick(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	flexPickID := c.Param("flexPickID")

	// Validasi ID
	if matchID == "" || teamID == "" || flexPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and  Flex Pick ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Cek keberadaan PriorityPick
	var flexPick models.FlexPick
	if err := config.DB.Where("flex_pick_id = ? AND match_team_detail_id = ?", flexPickID, matchTeamDetail.MatchTeamDetailID).
		First(&flexPick).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flex pick not found"})
		return
	}

	// Hapus FlexPick dari database
	if err := config.DB.Delete(&flexPick).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flex pick deleted successfully"})
}

func AddPriorityBan(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	// Validasi ID
	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Struct untuk menerima input JSON
	input := struct {
		HeroID  uint    `json:"hero_id" binding:"required"`
		Total   int     `json:"total" binding:"required"`
		Role    string  `json:"role" binding:"required,oneof=Gold Exp Roam Mid Jung"`
		BanRate float64 `json:"ban_rate" binding:"required"`
	}{}

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi keberadaan Hero
	var hero models.Hero
	if err := config.DB.First(&hero, input.HeroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Buat PriorityBan baru
	priorityBan := models.PriorityBan{
		MatchTeamDetailID: matchTeamDetail.MatchTeamDetailID,
		HeroID:            input.HeroID,
		Total:             input.Total,
		Role:              input.Role,
		BanRate:           input.BanRate,
	}

	// Simpan ke database
	if err := config.DB.Create(&priorityBan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Priority ban added successfully"})
}

func UpdatePriorityBan(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	priorityBanID := c.Param("priorityBanID")

	// Validasi ID
	if matchID == "" || teamID == "" || priorityBanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Priority Ban ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Cek keberadaan PriorityBan
	var priorityBan models.PriorityBan
	if err := config.DB.Where("priority_ban_id = ? AND match_team_detail_id = ?", priorityBanID, matchTeamDetail.MatchTeamDetailID).
		First(&priorityBan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Priority ban not found"})
		return
	}

	// Struct untuk menerima input JSON
	input := struct {
		HeroID  uint    `json:"hero_id" binding:"required"`
		Total   int     `json:"total" binding:"required"`
		Role    string  `json:"role" binding:"required,oneof=Gold Exp Roam Mid Jung"`
		BanRate float64 `json:"ban_rate" binding:"required"`
	}{}

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi keberadaan Hero
	var hero models.Hero
	if err := config.DB.First(&hero, input.HeroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Perbarui data PriorityBan
	priorityBan.HeroID = input.HeroID
	priorityBan.Total = input.Total
	priorityBan.Role = input.Role
	priorityBan.BanRate = input.BanRate

	// Simpan perubahan ke database
	if err := config.DB.Save(&priorityBan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Priority ban updated successfully"})
}

func GetAllPriorityBans(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	// Validasi ID
	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	type PriorityBanResponse struct {
		PriorityBanID     uint `json:"priority_ban_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		HeroID            uint `json:"hero_id"`
		Hero              struct {
			HeroID uint   `json:"hero_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		Total   int     `json:"total"`
		Role    string  `json:"role"`
		BanRate float64 `json:"ban_rate"`
	}

	var priorityBans []PriorityBanResponse

	// Query dengan WHERE untuk filter MatchTeamDetailID
	query := `
		SELECT 
			pb.priority_ban_id, pb.match_team_detail_id, pb.hero_id, 
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image, 
			pb.total, pb.role, pb.ban_rate
		FROM priority_bans pb
		JOIN heroes h ON pb.hero_id = h.hero_id
		WHERE pb.match_team_detail_id = ?
	`

	if err := config.DB.Raw(query, matchTeamDetail.MatchTeamDetailID).Scan(&priorityBans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve priority bans"})
		return
	}

	c.JSON(http.StatusOK, priorityBans)
}

func GetPriorityBanByID(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	priorityBanID := c.Param("priorityBanID")

	// Validasi ID
	if matchID == "" || teamID == "" || priorityBanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	type PriorityBanResponse struct {
		PriorityBanID     uint `json:"priority_ban_id"`
		MatchTeamDetailID uint `json:"match_team_detail_id"`
		HeroID            uint `json:"hero_id"`
		Hero              struct {
			HeroID uint   `json:"hero_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		Total   int     `json:"total"`
		Role    string  `json:"role"`
		BanRate float64 `json:"ban_rate"`
	}

	var priorityBan PriorityBanResponse

	query := `
		SELECT 
			pb.priority_ban_id, pb.match_team_detail_id, pb.hero_id, 
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image, 
			pb.total, pb.role, pb.ban_rate
		FROM priority_bans pb
		JOIN heroes h ON pb.hero_id = h.hero_id
		JOIN match_team_details mtd ON pb.match_team_detail_id = mtd.match_team_detail_id
		WHERE pb.priority_ban_id = ? AND mtd.match_id = ? AND mtd.team_id = ?
	`

	if err := config.DB.Raw(query, priorityBanID, matchID, teamID).Scan(&priorityBan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Priority ban not found"})
		return
	}

	c.JSON(http.StatusOK, priorityBan)
}

func DeletePriorityBan(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	priorityBanID := c.Param("priorityBanID")

	// Validasi ID
	if matchID == "" || teamID == "" || priorityBanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	// Cek keberadaan MatchTeamDetail
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).
		First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Cek keberadaan PriorityBan
	var priorityBan models.PriorityBan
	if err := config.DB.Where("priority_ban_id = ? AND match_team_detail_id = ?", priorityBanID, matchTeamDetail.MatchTeamDetailID).
		First(&priorityBan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Priority ban not found"})
		return
	}

	// Hapus PriorityBan dari database
	if err := config.DB.Delete(&priorityBan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete priority ban"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Priority ban deleted successfully"})
}
