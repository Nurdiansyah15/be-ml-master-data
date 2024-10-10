package controllers

import (
	"ml-master-data/config"
	"ml-master-data/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTeams(c *gin.Context) {
	var teams []models.Team

	if err := config.DB.Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func CreateTeam(c *gin.Context) {

	input := struct {
		Name string `json:"name" binding:"required"`
		Logo string `json:"logo"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Logo == "" {
		input.Logo = "https://placehold.co/400x600"
	}

	team := models.Team{
		Name: input.Name,
		Logo: input.Logo,
	}

	if err := config.DB.Create(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func UpdateTeam(c *gin.Context) {
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	var team models.Team

	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	input := struct {
		Name string `json:"name"`
		Logo string `json:"logo"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		team.Name = input.Name
	}
	if input.Logo != "" {
		team.Logo = input.Logo
	}

	if err := config.DB.Save(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

func CreatePlayerInTeam(c *gin.Context) {
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}
	var requestBody struct {
		Name  string `json:"name" binding:"required"`
		Role  string `json:"role" binding:"required"`
		Image string `json:"image" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	if requestBody.Image == "" {
		requestBody.Image = "https://placehold.co/400x600"
	}

	player := models.Player{
		Name:   requestBody.Name,
		Role:   requestBody.Role,
		Image:  requestBody.Image,
		TeamID: team.TeamID,
	}

	if err := config.DB.Create(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, player)
}

func CreateCoachInTeam(c *gin.Context) {
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	var requestBody struct {
		Name  string `json:"name" binding:"required"`
		Role  string `json:"role" binding:"required"`
		Image string `json:"image" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	if requestBody.Image == "" {
		requestBody.Image = "https://placehold.co/400x600"
	}

	coach := models.Coach{
		Name:   requestBody.Name,
		Role:   requestBody.Role,
		Image:  requestBody.Image,
		TeamID: team.TeamID,
	}

	if err := config.DB.Create(&coach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, coach)
}

func UpdatePlayerInTeam(c *gin.Context) {
	playerID := c.Param("playerID")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player ID is required"})
		return
	}
	var requestBody struct {
		Name  string `json:"name"`
		Role  string `json:"role"`
		Image string `json:"image"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var player models.Player
	if err := config.DB.First(&player, playerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	if requestBody.Name != "" {
		player.Name = requestBody.Name
	}
	if requestBody.Role != "" {
		player.Role = requestBody.Role
	}
	if requestBody.Image != "" {
		player.Image = requestBody.Image
	}

	if err := config.DB.Save(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, player)
}

func UpdateCoachInTeam(c *gin.Context) {
	coachID := c.Param("coachID")

	if coachID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coach ID is required"})
		return
	}

	var requestBody struct {
		Name  string `json:"name"`
		Role  string `json:"role"`
		Image string `json:"image"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var coach models.Coach
	if err := config.DB.First(&coach, coachID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	if requestBody.Name != "" {
		coach.Name = requestBody.Name
	}
	if requestBody.Role != "" {
		coach.Role = requestBody.Role
	}
	if requestBody.Image != "" {
		coach.Image = requestBody.Image
	}

	if err := config.DB.Save(&coach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coach)
}

func GetAllPlayersInTeam(c *gin.Context) {
	teamID := c.Param("teamID")

	var players []models.Player
	if err := config.DB.Where("team_id = ?", teamID).Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}

func GetAllCoachesInTeam(c *gin.Context) {
	teamID := c.Param("teamID")

	var coaches []models.Coach
	if err := config.DB.Where("team_id = ?", teamID).Find(&coaches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coaches)
}
