package controllers

import (
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

// @Summary Update user
// @Description Update username and password of the authenticated user
// @Tags Auth
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param update body dto.UpdateUserDto true "Update User"
// @Success 200 {object} models.User "Updated user data"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Old password does not match"
// @Failure 500 {string} string "Internal server error"
// @Router /user/update [put]
func UpdateUser(c *gin.Context) {
	// Ambil user dari context JWT
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userCtx.(models.User)

	// Bind input JSON ke struct DTO
	var updateDto dto.UpdateUserDto
	if err := c.ShouldBindJSON(&updateDto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Temukan user di database berdasarkan ID
	var user models.User
	if err := config.DB.First(&user, currentUser.UserID).Error; err != nil {
		c.JSON(500, gin.H{"error": "User not found"})
		return
	}

	// Validasi old_password jika password baru disediakan
	if updateDto.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updateDto.OldPassword)); err != nil {
			c.JSON(403, gin.H{"error": "Old password does not match"})
			return
		}

		// Hash password baru
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateDto.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	// Update username jika ada perubahan
	if updateDto.Username != "" {
		user.Username = updateDto.Username
	}

	// Simpan perubahan di database
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	// Berikan respon sukses dengan data user terbaru
	c.JSON(200, gin.H{"user": user})
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

	var user models.User
	user.UserID = userCtx.(models.User).UserID
	user.Username = userCtx.(models.User).Username

	c.JSON(200, gin.H{"user": user})
}
