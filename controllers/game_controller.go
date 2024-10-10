package controllers

import (
	"ml-master-data/config"
	"ml-master-data/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateMatchGame(c *gin.Context) {
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
		GameNumber       int  `json:"game_number" binding:"required"`
		FirstPickTeamID  uint `json:"first_pick_team_id" binding:"required"`
		SecondPickTeamID uint `json:"second_pick_team_id" binding:"required"`
		WinnerTeamID     uint `json:"winner_team_id" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var game models.Game

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game.GameNumber = input.GameNumber
	game.FirstPickTeamID = input.FirstPickTeamID
	game.SecondPickTeamID = input.SecondPickTeamID
	game.WinnerTeamID = input.WinnerTeamID
	game.MatchID = uint(matchIDInt)

	if err := config.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, game)
}

func UpdateMatchGame(c *gin.Context) {
	gameID := c.Param("gameID")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID is required"})
		return
	}

	var game models.Game
	if err := config.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	input := struct {
		GameNumber           int    `json:"game_number"`
		FirstPickTeamID      uint   `json:"first_pick_team_id"`
		SecondPickTeamID     uint   `json:"second_pick_team_id"`
		WinnerTeamID         uint   `json:"winner_team_id"`
		TrioMidOverallResult string `json:"trio_mid_overall_result"`
		EarlyGameResult      string `json:"early_game_result"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.GameNumber != 0 {
		game.GameNumber = input.GameNumber
	}
	if input.FirstPickTeamID != 0 {
		game.FirstPickTeamID = input.FirstPickTeamID
	}
	if input.SecondPickTeamID != 0 {
		game.SecondPickTeamID = input.SecondPickTeamID
	}
	if input.WinnerTeamID != 0 {
		game.WinnerTeamID = input.WinnerTeamID
	}
	if input.TrioMidOverallResult != "" {
		game.TrioMidOverallResult = input.TrioMidOverallResult
	}
	if input.EarlyGameResult != "" {
		game.EarlyGameResult = input.EarlyGameResult
	}

	if err := config.DB.Save(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, game)
}

func GetAllGameMatches(c *gin.Context) {
	matchID := c.Param("matchID")

	var games []models.Game
	if err := config.DB.Where("match_id = ?", matchID).Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, games)
}
