package controllers

import (
	"fmt"
	"ml-master-data/config"
	"ml-master-data/dto"
	"ml-master-data/models"
	"ml-master-data/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Login
// @Description Login to get JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param login body dto.LoginDto true "Login"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginDto dto.LoginDto

	if err := c.ShouldBindJSON(&loginDto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", loginDto.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password)); err != nil {
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

// @Summary Get user data
// @Description Get user data from JWT token
// @Tags Auth
// @Produce  json
// @Security Bearer
// @Success 200 {object} models.User "User data"
// @Failure 500 {string} string "Internal server error"
// @Router /me [get]
func Me(c *gin.Context) {
	userCtx, _ := c.Get("user")

	fmt.Println(userCtx)

	var user models.User
	user.UserID = userCtx.(models.User).UserID
	user.Username = userCtx.(models.User).Username

	c.JSON(200, gin.H{"user": user})
}
