package controllers

import (
	"ml-master-data/config"
	"ml-master-data/dto"
	"ml-master-data/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTournamentMatch godoc
// @Summary Create a match for a tournament
// @Description Create a match for a tournament and save its data
// @Security Bearer
// @Tags Match
// @Produce json
// @Param tournamentID path string true "Tournament ID"
// @Param dto body dto.MatchRequestDto true "Match request"
// @Success 201 {object} models.Match
// @Router /tournaments/{tournamentID}/matches [post]
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

	input := dto.MatchRequestDto{}

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

// @Summary Update a match
// @Description Update a match with the given match ID with the given information
// @Accept  json
// @Security Bearer
// @Tags Match
// @Produce  json
// @Param matchID path string true "Match ID"
// @Param match body dto.MatchRequestDto true "Match data"
// @Success 200 {object} models.Match "Match updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID} [put]
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

	input := dto.MatchRequestDto{}
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

// @Summary Get a match by ID
// @Description Get a match by ID
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Success 200 {object} dto.MatchResponseDto
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID} [get]
func GetMatchByID(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	var match dto.MatchResponseDto

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

// @Summary Get all matches for a tournament
// @Description Get all matches for a tournament with the given tournament ID
// @Security Bearer
// @Tags Match
// @Produce json
// @Param tournamentID path string true "Tournament ID"
// @Success 200 {array} dto.MatchResponseDto
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Tournament not found"
// @Failure 500 {string} string "Internal server error"
// @Router /tournaments/{tournamentID}/matches [get]
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

	var matches []dto.MatchResponseDto

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

// @Summary Add a player to a match
// @Description Add a player to a match by specifying the match ID, team ID, and player ID
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param dto body dto.PlayerMatchRequestDto true "Player match request"
// @Success 201 {string} string "Player match added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/players [post]
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

	input := dto.PlayerMatchRequestDto{}

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

// @Summary Remove a player match
// @Description Remove a player match with the given match ID, team ID and player ID
// @Accept  json
// @Security Bearer
// @Tags Match
// @Produce  json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param playerID path string true "Player ID"
// @Success 200 {string} string "Player match removed successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/players/{playerID} [delete]
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

// @Summary Get all players for a match and team
// @Description Get all players for a match and team with the given match ID and team ID
// @Accept  json
// @Security Bearer
// @Tags Match
// @Produce  json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Success 200 {array} dto.PlayerMatchResponseDto
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/players [get]
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

	var players []dto.PlayerMatchResponseDto

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

// @Summary Add a coach match
// @Description Add a coach match with the given match ID, team ID and coach ID
// @Accept  json
// @Security Bearer
// @Tags Match
// @Produce  json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param coachID body int true "Coach ID"
// @Success 201 {string} string "Coach match added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/coaches [post]
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

	input := dto.CoachMatchRequestDto{}
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

// @Summary Remove a coach match
// @Description Remove a coach match with the given match ID, team ID, and coach ID
// @Accept  json
// @Security Bearer
// @Tags Match
// @Produce  json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param coachID path string true "Coach ID"
// @Success 200 {string} string "Coach match removed successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or coach not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/coaches/{coachID} [delete]
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

// @Summary Get all coaches match
// @Description Get all coaches match with the given match ID and team ID
// @Accept  json
// @Security Bearer
// @Tags Match
// @Produce  json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Success 200 {array} dto.CoachMatchResponseDto "Coaches match found"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/coaches [get]
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

	var coaches []dto.CoachMatchResponseDto

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

// @Summary Add hero pick
// @Description Add hero pick to match
// @ID add-hero-pick
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param heroPick body dto.HeroPickRequestDto true "Hero pick"
// @Success 200 {string} string "Hero pick added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/hero-picks [post]
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

	input := dto.HeroPickRequestDto{}
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

// @Summary Update hero pick
// @Description Update hero pick in match
// @Accept  json
// @Security Bearer
// @Tags Match
// @Produce  json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param heroPickID path string true "Hero pick ID"
// @Param heroPick body dto.HeroPickRequestDto true "Hero pick"
// @Success 200 {string} string "Hero pick updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/hero-picks/{heroPickID} [put]
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

	input := dto.HeroPickRequestDto{}
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

// @Summary Remove hero pick
// @Description Remove hero pick from match
// @ID remove-hero-pick
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param heroPickID path string true "Hero pick ID"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match, team, or hero pick not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/hero-picks/{heroPickID} [delete]
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

// @Summary Get all hero picks
// @Description Get all hero picks in a match with the given team ID
// @ID get-all-hero-picks
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Success 200 {array} dto.HeroPickResponseDto
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Hero picks not found"
// @Router /matches/{matchID}/teams/{teamID}/hero-picks [get]
func GetAllHeroPicks(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	var picks []dto.HeroPickResponseDto
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

// @Summary Add hero ban
// @Description Add hero ban to match
// @ID add-hero-ban
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param heroBan body dto.HeroBanRequestDto true "Hero ban"
// @Success 200 {string} string "Hero ban added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/hero-bans [post]
func AddHeroBan(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	// Validasi keberadaan match dan team
	var matchTeamDetail models.MatchTeamDetail
	if err := config.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&matchTeamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or team not found"})
		return
	}

	// Binding input JSON
	input := dto.HeroBanRequestDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Memulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}()

	heroBan := models.HeroBan{
		MatchTeamDetailID: matchTeamDetail.MatchTeamDetailID,
		HeroID:            input.HeroID,
		FirstPhase:        input.FirstPhase,
		SecondPhase:       input.SecondPhase,
	}
	if err := tx.Create(&heroBan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Membuat HeroBan dan menyimpannya
	for _, ban := range input.HeroBanGame {
		heroBanGame := models.HeroBanGame{
			HeroBanID:  heroBan.HeroBanID,
			GameNumber: ban.GameNumber,
			IsBanned:   ban.IsBanned,
		}

		if err := tx.Create(&heroBanGame).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit jika semua operasi sukses
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Respon sukses
	c.JSON(http.StatusCreated, gin.H{"message": "Hero ban added successfully"})
}

// @Summary Update hero ban
// @Description Update hero ban in match
// @ID update-hero-ban
// @Accept  json
// @Security Bearer
// @Tags Match
// @Produce  json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param HeroBanID path string true "Hero ban ID"
// @Param heroBan body dto.HeroBanRequestDto true "Hero ban"
// @Success 200 {string} string "Hero ban updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/hero-bans/{HeroBanID} [put]
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

	input := dto.HeroBanRequestDto{}
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

// @Summary Remove hero ban
// @Description Remove hero ban from match
// @ID remove-hero-ban
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param HeroBanID path string true "HeroBan ID"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/hero-bans/{HeroBanID} [delete]
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

// @Summary Get all hero bans
// @Description Get all hero bans of a match by team
// @ID get-all-hero-bans
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Success 200 {array} dto.HeroBanResponseDto "Hero bans"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/hero-bans [get]
func GetAllHeroBans(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")

	if matchID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Team ID are required"})
		return
	}

	var bans []dto.HeroBanResponseDto
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

// @Summary Add priority pick
// @Description Add priority pick to match
// @ID add-priority-pick
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param priorityPick body dto.PriorityPickRequestDto true "Priority pick"
// @Success 201 {string} string "Priority pick added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/priority-picks [post]
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
	input := dto.PriorityPickRequestDto{}

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

// @Summary Update priority pick
// @Description Update priority pick in match
// @Accept  json
// @Security Bearer
// @Tags Match
// @Produce  json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param priorityPickID path string true "Priority pick ID"
// @Param priorityPick body dto.PriorityPickRequestDto true "Priority pick"
// @Success 200 {string} string "Priority pick updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/priority-picks/{priorityPickID} [put]
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
	input := dto.PriorityPickRequestDto{}

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

// @Summary Get all priority picks
// @Description Get all priority picks of a match by team
// @ID get-all-priority-picks
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Success 200 {array} dto.PriorityPickResponseDto "Priority pick list"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/priority-picks [get]
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

	var priorityPicks []dto.PriorityPickResponseDto

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

// @Summary Get priority pick by ID
// @Description Get a priority pick by its ID
// @ID get-priority-pick-by-id
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param priorityPickID path string true "Priority pick ID"
// @Success 200 {object} dto.PriorityPickResponseDto
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Priority pick not found"
// @Router /matches/{matchID}/teams/{teamID}/priority-picks/{priorityPickID} [get]
func GetPriorityPickByID(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	priorityPickID := c.Param("priorityPickID")

	// Validasi ID
	if matchID == "" || teamID == "" || priorityPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	var priorityPick dto.PriorityPickResponseDto

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

// @Summary Delete priority pick
// @Description Delete a priority pick by its ID
// @ID delete-priority-pick
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param priorityPickID path string true "Priority Pick ID"
// @Success 200 {string} string "Priority pick deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Priority pick not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/priority-picks/{priorityPickID} [delete]
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

// @Summary Add flex pick
// @Description Add flex pick to match
// @ID add-flex-pick
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param flexPick body dto.FlexPickRequestDto true "Flex pick"
// @Success 201 {string} string "Flex pick added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/flex-picks [post]
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
	input := dto.FlexPickRequestDto{}

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

// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param flexPickID path string true "Flex Pick ID"
// @Param flexPick body dto.FlexPickRequestDto true "Flex pick"
// @Success 200 {string} string "Flex pick updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match, team, or flex pick not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/flex-picks/{flexPickID} [put]
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
	input := dto.FlexPickRequestDto{}

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

// @Summary Get all flex picks
// @Description Get all flex picks of a match by team
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Success 200 {array} dto.FlexPickResponseDto "Flex pick list"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/flex-picks [get]
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

	var flexPicks []dto.FlexPickResponseDto

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

// @Summary Get flex pick by ID
// @Description Get a flex pick by its ID
// @ID get-flex-pick-by-id
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param flexPickID path string true "Flex pick ID"
// @Success 200 {object} dto.FlexPickResponseDto
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Flex pick not found"
// @Router /matches/{matchID}/teams/{teamID}/flex-picks/{flexPickID} [get]
func GetFlexPickByID(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	flexPickID := c.Param("flexPickID")

	// Validasi ID
	if matchID == "" || teamID == "" || flexPickID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	var flexPick dto.FlexPickResponseDto

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

// @Summary Delete flex pick
// @Description Delete flex pick by ID
// @ID delete-flex-pick
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param flexPickID path string true "Flex pick ID"
// @Success 200 {string} string "Flex pick deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match, team, or flex pick not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/flex-picks/{flexPickID} [delete]
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

// @Summary Add priority ban
// @Description Add priority ban to match
// @ID add-priority-ban
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param priorityBan body dto.PriorityBanRequestDto true "Priority ban"
// @Success 201 {string} string "Priority ban added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/priority-bans [post]
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
	input := dto.PriorityBanRequestDto{}

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

// @Summary Update priority ban
// @Description Update priority ban in match
// @ID update-priority-ban
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param priorityBanID path string true "Priority Ban ID"
// @Param priorityBan body dto.PriorityBanRequestDto true "Priority Ban"
// @Success 200 {string} string "Priority ban updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match, team, or priority ban not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/priority-bans/{priorityBanID} [put]
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
	input := dto.PriorityBanRequestDto{}

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

// @Summary Get all priority bans
// @Description Get all priority bans in a match with specific team
// @ID get-all-priority-bans
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Success 200 {array} dto.PriorityBanResponseDto "Priority bans"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/priority-bans [get]
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

	var priorityBans []dto.PriorityBanResponseDto

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

// @Summary Get priority ban by ID
// @Description Get priority ban by ID in a match
// @ID get-priority-ban-by-id
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param priorityBanID path string true "Priority Ban ID"
// @Success 200 {object} dto.PriorityBanResponseDto "Priority Ban"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Priority ban not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/priority-bans/{priorityBanID} [get]
func GetPriorityBanByID(c *gin.Context) {
	matchID := c.Param("matchID")
	teamID := c.Param("teamID")
	priorityBanID := c.Param("priorityBanID")

	// Validasi ID
	if matchID == "" || teamID == "" || priorityBanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Team ID, and Hero ID are required"})
		return
	}

	var priorityBan dto.PriorityBanResponseDto

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

// @Summary Delete priority ban
// @Description Delete priority ban in match
// @ID delete-priority-ban
// @Accept json
// @Security Bearer
// @Tags Match
// @Produce json
// @Param matchID path string true "Match ID"
// @Param teamID path string true "Team ID"
// @Param priorityBanID path string true "Priority Ban ID"
// @Success 200 {string} string "Priority ban deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match, team, or priority ban not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/teams/{teamID}/priority-bans/{priorityBanID} [delete]
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
