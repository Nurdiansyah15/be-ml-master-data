package controllers

import (
	"fmt"
	"ml-master-data/config"
	"ml-master-data/models"
	"ml-master-data/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func Me(c *gin.Context) {
	userCtx, _ := c.Get("user")

	fmt.Println(userCtx)

	var user models.User
	user.UserID = userCtx.(models.User).UserID
	user.Username = userCtx.(models.User).Username

	c.JSON(200, gin.H{"user": user})
}
