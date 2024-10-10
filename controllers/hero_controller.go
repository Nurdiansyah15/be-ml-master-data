package controllers

import (
	"ml-master-data/config"
	"ml-master-data/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllHeroes(c *gin.Context) {
	var heroes []models.Hero

	if err := config.DB.Find(&heroes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, heroes)
}

func CreateHero(c *gin.Context) {
	input := struct {
		Name      string `json:"name" binding:"required"`
		HeroImage string `json:"hero_image"`
	}{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.HeroImage == "" {
		input.HeroImage = "https://placehold.co/400x600"
	}

	hero := models.Hero{
		Name:      input.Name,
		HeroImage: input.HeroImage,
	}

	if err := config.DB.Create(&hero).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hero)
}

func GetHeroByID(c *gin.Context) {
	heroID := c.Param("heroID")

	var hero models.Hero
	if err := config.DB.First(&hero, heroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	c.JSON(http.StatusOK, hero)
}

func UpdateHero(c *gin.Context) {
	heroID := c.Param("heroID")
	if heroID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hero ID is required"})
		return
	}

	var hero models.Hero

	if err := config.DB.First(&hero, heroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	input := struct {
		Name      string `json:"name"`
		HeroImage string `json:"hero_image"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		hero.Name = input.Name
	}
	if input.HeroImage != "" {
		hero.HeroImage = input.HeroImage
	}

	if err := config.DB.Save(&hero).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hero)
}
