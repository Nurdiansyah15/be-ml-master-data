package controllers

import (
	"fmt"
	"io"
	"ml-master-data/config"
	"ml-master-data/dto"
	"ml-master-data/models"
	"ml-master-data/services"
	"ml-master-data/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Tags Game
// @Summary Create a new game
// @Description Create a new game for the specified match with additional information including the first pick team ID, second pick team ID, winner team ID, game number, video link, and optionally a full draft image.
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param matchID path string true "Match ID"
// @Param first_pick_team_id formData integer true "First Pick Team ID"
// @Param second_pick_team_id formData integer true "Second Pick Team ID"
// @Param winner_team_id formData integer true "Winner Team ID"
// @Param game_number formData integer true "Game Number"
// @Param video_link formData string false "Video Link"
// @Param full_draft_image formData file false "Full Draft Image"
// @Success 201 {object} models.Game "Game created successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Match not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games [post]
func CreateGame(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	// Cari match berdasarkan ID
	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Parse form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Inisialisasi struct untuk menerima input
	input := dto.GameRequestDto{}

	// Mengisi struct dari form data dengan penanganan error
	firstPickTeamID, err := strconv.ParseUint(c.Request.FormValue("first_pick_team_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid first_pick_team_id"})
		return
	}
	input.FirstPickTeamID = uint(firstPickTeamID)

	secondPickTeamID, err := strconv.ParseUint(c.Request.FormValue("second_pick_team_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid second_pick_team_id"})
		return
	}
	input.SecondPickTeamID = uint(secondPickTeamID)

	winnerTeamID, err := strconv.ParseUint(c.Request.FormValue("winner_team_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid winner_team_id"})
		return
	}
	input.WinnerTeamID = uint(winnerTeamID)

	gameNumber, err := strconv.Atoi(c.Request.FormValue("game_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game_number"})
		return
	}
	input.GameNumber = gameNumber

	input.VideoLink = c.Request.FormValue("video_link")

	// Validasi input
	if input.FirstPickTeamID == 0 || input.SecondPickTeamID == 0 || input.WinnerTeamID == 0 || input.GameNumber == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	var fullDraftImagePath string

	// Cek apakah ada file gambar
	file, header, err := c.Request.FormFile("full_draft_image")
	if err == nil {
		// Memeriksa ukuran file
		if header.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		// Simpan file gambar
		newFileName := utils.GenerateUniqueFileName("draft") + ext
		fullDraftImagePath = fmt.Sprintf("public/images/%s", newFileName)

		dst, err := os.Create(fullDraftImagePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
			return
		}
		defer dst.Close()

		if _, err = io.Copy(dst, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		// Atur path lengkap dengan BASE_URL
		fullDraftImagePath = os.Getenv("BASE_URL") + "/" + fullDraftImagePath
	}

	// Buat instance Game
	game := models.Game{
		MatchID:          match.MatchID,
		FirstPickTeamID:  input.FirstPickTeamID,
		SecondPickTeamID: input.SecondPickTeamID,
		WinnerTeamID:     input.WinnerTeamID,
		GameNumber:       input.GameNumber,
		VideoLink:        input.VideoLink,
		FullDraftImage:   fullDraftImagePath, // Path file atau kosong jika tidak ada gambar
	}

	// Simpan game ke database
	if err := config.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusCreated, game)
}

// @Tags Game
// @Summary Update a game
// @Description Update a game with the given game ID and match ID with the given information
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Param first_pick_team_id formData integer false "First Pick Team ID"
// @Param second_pick_team_id formData integer false "Second Pick Team ID"
// @Param winner_team_id formData integer false "Winner Team ID"
// @Param game_number formData integer false "Game Number"
// @Param video_link formData string false "Video Link"
// @Param full_draft_image formData file false "Full Draft Image"
// @Success 200 {object} models.Game "Game updated successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID} [put]
func UpdateGame(c *gin.Context) {
	gameID := c.Param("gameID")
	matchID := c.Param("matchID")

	if gameID == "" || matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID and Match ID are required"})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	var game models.Game
	if err := config.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// Parse form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Tangani file gambar jika ada
	file, header, err := c.Request.FormFile("full_draft_image")
	if err == nil {
		// Memeriksa ukuran file
		if header.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		// Jika ada gambar sebelumnya, hapus
		if game.FullDraftImage != "" {
			oldImagePath := strings.Replace(game.FullDraftImage, os.Getenv("BASE_URL")+"/", "", 1)
			if _, err := os.Stat(oldImagePath); err == nil {
				if err := os.Remove(oldImagePath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old image"})
					return
				}
			}
		}

		// Simpan gambar baru
		newFileName := utils.GenerateUniqueFileName("draft") + ext
		newImagePath := fmt.Sprintf("public/images/%s", newFileName)

		dst, err := os.Create(newImagePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
			return
		}
		defer dst.Close()

		if _, err = io.Copy(dst, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new image"})
			return
		}

		game.FullDraftImage = os.Getenv("BASE_URL") + "/" + newImagePath
	}

	// Update game fields if provided in the form
	if firstPickTeamID, err := strconv.ParseUint(c.Request.FormValue("first_pick_team_id"), 10, 32); err == nil {
		game.FirstPickTeamID = uint(firstPickTeamID)
	}
	if secondPickTeamID, err := strconv.ParseUint(c.Request.FormValue("second_pick_team_id"), 10, 32); err == nil {
		game.SecondPickTeamID = uint(secondPickTeamID)
	}
	if winnerTeamID, err := strconv.ParseUint(c.Request.FormValue("winner_team_id"), 10, 32); err == nil {
		game.WinnerTeamID = uint(winnerTeamID)
	}
	if gameNumber, err := strconv.Atoi(c.Request.FormValue("game_number")); err == nil {
		game.GameNumber = gameNumber
	}
	if videoLink := c.Request.FormValue("video_link"); videoLink != "" {
		game.VideoLink = videoLink
	}

	// Simpan perubahan ke database
	if err := config.DB.Save(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusOK, game)
}

// @Tags Game
// @Summary Get all games for a match
// @Description Get all games for a match with the given match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param matchID path string true "Match ID"
// @Success 200 {array} dto.GameResponseDto
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games [get]
func GetAllGames(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}

	var match models.Match
	if err := config.DB.Where("match_id = ?", matchID).Find(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var games = []dto.GameResponseDto{}

	query := `
		SELECT 
			g.game_id, g.match_id, g.first_pick_team_id, 
			t1.team_id AS first_team_team_id, t1.name AS first_team_name, t1.image AS first_team_image,
			g.second_pick_team_id, 
			t2.team_id AS second_team_team_id, t2.name AS second_team_name, t2.image AS second_team_image,
			g.winner_team_id, 
			t3.team_id AS winner_team_team_id, t3.name AS winner_team_name, t3.image AS winner_team_image,
			g.game_number, g.video_link, g.full_draft_image
		FROM games g
		JOIN teams t1 ON g.first_pick_team_id = t1.team_id
		JOIN teams t2 ON g.second_pick_team_id = t2.team_id
		JOIN teams t3 ON g.winner_team_id = t3.team_id
		WHERE g.match_id = ?
	`

	if err := config.DB.Raw(query, match.MatchID).Scan(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, games)
}

// @Tags Game
// @Summary Get a game by ID
// @Description Get a game by ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Success 200 {object} dto.GameResponseDto
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID} [get]
func GetGameByID(c *gin.Context) {
	gameID := c.Param("gameID")
	matchID := c.Param("matchID")

	if gameID == "" || matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID and Match ID are required"})
		return
	}

	var game dto.GameResponseDto

	query := `
		SELECT 
			g.game_id, g.match_id, g.first_pick_team_id, 
			t1.team_id AS first_team_id, t1.name AS first_team_name, t1.image AS first_team_image,
			g.second_pick_team_id, 
			t2.team_id AS second_team_id, t2.name AS second_team_name, t2.image AS second_team_image,
			g.winner_team_id, 
			t3.team_id AS winner_team_id, t3.name AS winner_team_name, t3.image AS winner_team_image,
			g.game_number, g.video_link, g.full_draft_image
		FROM games g
		JOIN teams t1 ON g.first_pick_team_id = t1.team_id
		JOIN teams t2 ON g.second_pick_team_id = t2.team_id
		JOIN teams t3 ON g.winner_team_id = t3.team_id
		WHERE g.match_id = ? AND g.game_id = ?
	`

	if err := config.DB.Raw(query, matchID, gameID).Scan(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// @Tags Game
// @Summary Delete a game
// @Description Delete a game with the given game ID and match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param matchID path string true "Match ID"
// @Param gameID path string true "Game ID"
// @Success 200 {string} string "Game deleted successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID} [delete]
func RemoveGame(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID is required"})
		return
	}
	gameID := c.Param("gameID")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID is required"})
		return
	}

	match := models.Match{}
	if err := config.DB.Where("match_id = ?", matchID).First(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	game := models.Game{}
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	if err := services.DeleteGame(config.DB, game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game deleted successfully"})
}

// @Tags Game
// @Summary Add a lord result
// @Description Add a lord result for a game with the given game ID and match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Param lordResult body dto.LordResultRequestDto true "Lord result data"
// @Success 201 {string} string "Lord result added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/lord-results [post]
func AddLordResult(c *gin.Context) {

	gameID := c.Param("gameID")
	matchID := c.Param("matchID")

	if gameID == "" || matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID and Match ID are required"})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	var game models.Game
	if err := config.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var input dto.LordResultRequestDto

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lordResult := models.LordResult{
		GameID:   game.GameID,
		TeamID:   input.TeamID,
		Phase:    input.Phase,
		Setup:    input.Setup,
		Initiate: input.Initiate,
		Result:   input.Result,
	}

	if err := config.DB.Create(&lordResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Lord result added successfully"})
}

// @Tags Game
// @Summary Update a LordResult
// @Description Update a LordResult with the given game ID, match ID, and Lord Result ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Param lordResultID path string true "Lord Result ID"
// @Param lordResult body dto.LordResultRequestDto true "Lord result data"
// @Success 200 {object} models.LordResult "Lord result updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Lord result not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/lord-results/{lordResultID} [put]
func UpdateLordResult(c *gin.Context) {
	gameID := c.Param("gameID")
	matchID := c.Param("matchID")
	lordResultID := c.Param("lordResultID")

	// Validasi parameter
	if gameID == "" || matchID == "" || lordResultID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID, Match ID, and Lord Result ID are required"})
		return
	}

	// Validasi keberadaan Match dan Game
	if err := config.DB.First(&models.Match{}, "game_id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "game_id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// Cek apakah LordResult tersedia
	var lordResult models.LordResult
	if err := config.DB.First(&lordResult, "lord_result_id = ?", lordResultID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lord result not found"})
		return
	}

	// Struct input untuk pembaruan
	input := dto.LordResultRequestDto{}

	// Validasi input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update hanya field yang tidak bernilai null
	lordResult.TeamID = input.TeamID
	lordResult.Phase = input.Phase
	lordResult.Setup = input.Setup
	lordResult.Initiate = input.Initiate
	lordResult.Result = input.Result

	// Simpan perubahan ke database
	if err := config.DB.Save(&lordResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lordResult)
}

// @Tags Game
// @Summary Delete a LordResult
// @Description Delete a LordResult with the given game ID, match ID, and Lord Result ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Param lordResultID path string true "Lord Result ID"
// @Success 200 {string} string "Lord result deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Lord result not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/lord-results/{lordResultID} [delete]
func RemoveLordResult(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	lordResultID := c.Param("lordResultID")

	// Validasi parameter
	if matchID == "" || gameID == "" || lordResultID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Game ID, and Lord Result ID are required"})
		return
	}

	// Validasi keberadaan Match dan Game
	if err := config.DB.First(&models.Match{}, "match_id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "game_id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// Validasi keberadaan LordResult
	var lordResult models.LordResult
	if err := config.DB.First(&lordResult, "lord_result_id = ? AND game_id = ?", lordResultID, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lord result not found"})
		return
	}

	// Hapus LordResult
	if err := config.DB.Delete(&lordResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lord result deleted successfully"})
}

// @Tags Game
// @Summary Get all lord results for a game
// @Description Get all lord results for a game with the given game ID and match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param gameID path string true "Game ID"
// @Success 200 {array} dto.LordResultResponseDto
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/lord-results [get]
func GetAllLordResults(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")

	// Validasi keberadaan match dan game
	if err := config.DB.First(&models.Match{}, "match_id = ?", teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "game_id = ?", gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var results []dto.LordResultResponseDto

	query := `
		SELECT 
			l.lord_result_id, l.game_id,
			t.team_id AS team_team_id, t.name AS team_name, t.image AS team_image,
			l.phase, l.setup, l.initiate, l.result
		FROM lord_results l
		JOIN teams t ON l.team_id = t.team_id
		WHERE l.game_id = ? AND l.team_id = ?
	`

	if err := config.DB.Raw(query, gameID, teamID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// @Tags Game
// @Summary Get a LordResult by ID
// @Description Get a LordResult by the given game ID, match ID, and Lord Result ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param lordResultID path string true "Lord Result ID"
// @Success 200 {object} dto.LordResultResponseDto
// @Failure 404 {string} string "Match or game or Lord result not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/lord-results/{lordResultID} [get]
func GetLordResultByID(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")
	lordResultID := c.Param("lordResultID")

	// Validasi keberadaan match dan game
	if err := config.DB.First(&models.Match{}, "match_id = ?", teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "game_id = ?", gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var result dto.LordResultResponseDto

	query := `
		SELECT 
			l.lord_result_id, l.game_id,
			t.team_id AS team_team_id, t.name AS team_name, t.image AS team_image,
			l.phase, l.setup, l.initiate, l.result
		FROM lord_results l
		JOIN teams t ON l.team_id = t.team_id
		WHERE l.lord_result_id = ? AND l.game_id = ? AND l.team_id = ?
	`

	if err := config.DB.Raw(query, lordResultID, gameID, teamID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lord result not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Tags Game
// @Summary Add a turtle result
// @Description Add a turtle result for a game with the given game ID and match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Param turtleResult body dto.TurtleResultRequestDto true "Turtle result data"
// @Success 201 {string} string "Turtle result added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/turtle-results [post]
func AddTurtleResult(c *gin.Context) {
	gameID := c.Param("gameID")
	matchID := c.Param("matchID")

	if gameID == "" || matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID and Match ID are required"})
		return
	}

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	var game models.Game
	if err := config.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	input := dto.TurtleResultRequestDto{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	turtleResult := models.TurtleResult{
		TeamID:   input.TeamID,
		Phase:    input.Phase,
		Setup:    input.Setup,
		Initiate: input.Initiate,
		Result:   input.Result,
	}

	turtleResult.GameID = game.GameID

	if err := config.DB.Create(&turtleResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Turtle result added successfully"})
}

// @Tags Game
// @Summary Update a turtle result
// @Description Update a turtle result with the given game ID, match ID, and turtle result ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Param turtleResultID path string true "Turtle Result ID"
// @Param turtleResult body dto.TurtleResultRequestDto true "Turtle result data"
// @Success 200 {object} models.TurtleResult "Turtle result updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Turtle result not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/turtle-results/{turtleResultID} [put]
func UpdateTurtleResult(c *gin.Context) {
	gameID := c.Param("gameID")
	matchID := c.Param("matchID")
	turtleResultID := c.Param("turtleResultID")

	if gameID == "" || matchID == "" || turtleResultID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID, Match ID, and Turtle Result ID are required"})
		return
	}

	if err := config.DB.First(&models.Match{}, "match_id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "game_id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var turtleResult models.TurtleResult
	if err := config.DB.First(&turtleResult, "turtle_result_id = ?", turtleResultID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turtle result not found"})
		return
	}

	input := dto.TurtleResultRequestDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	turtleResult.TeamID = input.TeamID
	turtleResult.Phase = input.Phase
	turtleResult.Setup = input.Setup
	turtleResult.Initiate = input.Initiate
	turtleResult.Result = input.Result

	if err := config.DB.Save(&turtleResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, turtleResult)
}

// @Tags Game
// @Summary Delete a TurtleResult
// @Description Delete a TurtleResult with the given game ID, match ID, and Turtle Result ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Param turtleResultID path string true "Turtle Result ID"
// @Success 200 {string} string "Turtle result deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Turtle result not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/turtle-results/{turtleResultID} [delete]
func RemoveTurtleResult(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	turtleResultID := c.Param("turtleResultID")

	if matchID == "" || gameID == "" || turtleResultID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Game ID, and Turtle Result ID are required"})
		return
	}

	if err := config.DB.First(&models.Match{}, "match_id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "game_id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var turtleResult models.TurtleResult
	if err := config.DB.First(&turtleResult, "turtle_result_id = ? AND game_id = ?", turtleResultID, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turtle result not found"})
		return
	}

	if err := config.DB.Delete(&turtleResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Turtle result deleted successfully"})
}

// @Tags Game
// @Summary Get all TurtleResults
// @Description Get all TurtleResults for a game with the given game ID and match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Match ID"
// @Success 200 {array} dto.TurtleResultResponseDto
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{teamID}/games/{gameID}/turtle-results [get]
func GetAllTurtleResults(c *gin.Context) {

	teamID := c.Param("teamID")
	gameID := c.Param("gameID")

	// Validasi keberadaan match dan game
	if err := config.DB.First(&models.Match{}, "match_id = ?", teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "game_id = ?", gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var results []dto.TurtleResultResponseDto

	query := `
		SELECT 
			tr.turtle_result_id, tr.game_id,
			t.team_id AS team_team_id, t.name AS team_name, t.image AS team_image,
			tr.phase, tr.setup, tr.initiate, tr.result
		FROM turtle_results tr
		JOIN teams t ON tr.team_id = t.team_id
		WHERE tr.game_id = ? AND tr.team_id  = ?
	`

	if err := config.DB.Raw(query, gameID, teamID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// @Tags Game
// @Summary Get a TurtleResult by ID
// @Description Get a TurtleResult by ID for a game with the given game ID and match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param gameID path string true "Game ID"
// @Param turtleResultID path string true "Turtle result ID"
// @Success 200 {object} dto.TurtleResultResponseDto
// @Failure 404 {string} string "Match, Game, or Turtle result not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/turtle-results/{turtleResultID} [get]
func GetTurtleResultByID(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")
	turtleResultID := c.Param("turtleResultID")

	// Validasi keberadaan match dan game
	if err := config.DB.First(&models.Match{}, "match_id = ?", teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "game_id = ?", gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var result dto.TurtleResultResponseDto

	query := `
		SELECT 
			tr.turtle_result_id, tr.game_id,
			t.team_id AS team_team_id, t.name AS team_name, t.image AS team_image,
			tr.phase, tr.setup, tr.initiate, tr.result
		FROM turtle_results tr
		JOIN teams t ON tr.team_id = t.team_id
		WHERE tr.turtle_result_id = ? AND tr.game_id = ? AND tr.team_id = ?
	`

	if err := config.DB.Raw(query, turtleResultID, gameID, teamID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turtle result not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Tags Game
// @Summary Add an explaner
// @Description Add an explaner for a game with the given game ID and match ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param explaner body dto.ExplanerRequestDto true "Explaner data"
// @Success 201 {string} string "Explaner added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/explaners [post]
func AddExplaner(c *gin.Context) {
	gameID := c.Param("gameID")
	teamID := c.Param("teamID")

	if teamID == "" || gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID and Team ID are required"})
		return
	}

	uintTeamID, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	var game = models.Game{}
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	match := models.Match{}
	if err := config.DB.Where("match_id = ?", game.MatchID).First(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	if match.TeamAID != uint(uintTeamID) && match.TeamBID != uint(uintTeamID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is not part of the match"})
		return
	}

	var input = dto.ExplanerRequestDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// cek heriID duplicated
	var explanerExists models.Explaner
	if err := config.DB.Where("game_id = ? AND team_id = ? AND hero_id = ?", game.GameID, teamID, input.HeroID).First(&explanerExists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Explaner already exists"})
		return
	}

	explaner := models.Explaner{
		GameID:      game.GameID,
		TeamID:      uint(uintTeamID),
		HeroID:      input.HeroID,
		EarlyResult: input.EarlyResult,
	}

	if err := config.DB.Create(&explaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Explaner added successfully"})
}

// @Tags Game
// @Summary Update an Explaner
// @Description Update an Explaner with the given game ID, match ID, and Explaner ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param explanerID path string true "Explaner ID"
// @Param explaner body dto.ExplanerRequestDto true "Explaner data"
// @Success 200 {object} models.Explaner "Explaner updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Explaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/explaners/{explanerID} [put]
func UpdateExplaner(c *gin.Context) {
	gameID := c.Param("gameID")
	teamID := c.Param("teamID")
	explanerID := c.Param("explanerID")

	if teamID == "" || gameID == "" || explanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID, Team ID, and Explaner ID are required"})
		return
	}

	uintTeamID, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	var game = models.Game{}
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	match := models.Match{}
	if err := config.DB.Where("match_id = ?", game.MatchID).First(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	if match.TeamAID != uint(uintTeamID) && match.TeamBID != uint(uintTeamID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is not part of the match"})
		return
	}

	var explaner models.Explaner
	if err := config.DB.First(&explaner, "explaner_id = ?", explanerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Explaner not found"})
		return
	}

	input := dto.ExplanerRequestDto{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// cek hero are duplicated
	var explanerExists models.Explaner
	if err := config.DB.Where("game_id = ? AND team_id = ? AND hero_id = ? AND explaner_id != ?", game.GameID, teamID, input.HeroID, explanerID).First(&explanerExists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Explaner already exists"})
		return
	}

	explaner.TeamID = uint(uintTeamID)
	explaner.HeroID = input.HeroID
	explaner.EarlyResult = input.EarlyResult

	if err := config.DB.Save(&explaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, explaner)
}

// @Tags Game
// @Summary Delete an Explaner
// @Description Delete an Explaner with the given game ID, match ID, and Explaner ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param explanerID path string true "Explaner ID"
// @Success 200 {string} string "Explaner deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Explaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/explaners/{explanerID} [delete]
func RemoveExplaner(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")
	explanerID := c.Param("explanerID")

	if teamID == "" || gameID == "" || explanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID, Team ID, and Explaner ID are required"})
		return
	}

	var game = models.Game{}
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	var explaner models.Explaner
	if err := config.DB.First(&explaner, "explaner_id = ? AND game_id = ?", explanerID, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Explaner not found"})
		return
	}

	if err := config.DB.Delete(&explaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Explaner deleted successfully"})
}

// @Tags Game
// @Summary Get all Explaners for a game
// @Description Get all Explaners for a game with the given game ID and match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param gameID path string true "Game ID"
// @Success 200 {array} dto.ExplanerResponseDto
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/explaners [get]
func GetAllExplaners(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")

	if teamID == "" || gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID and Team ID are required"})
		return
	}

	var game = models.Game{}
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	var results []dto.ExplanerResponseDto

	query := `
		SELECT 
			e.explaner_id, e.game_id, e.team_id, e.hero_id, e.early_result,
			t.team_id AS team_team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image
		FROM explaners e
		JOIN teams t ON e.team_id = t.team_id
		JOIN heros h ON e.hero_id = h.hero_id
		WHERE e.game_id = ? AND e.team_id = ?
	`

	if err := config.DB.Raw(query, gameID, teamID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// @Tags Game
// @Summary Get an explaner by ID
// @Description Get an explaner by ID for a game with the given game ID and match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param gameID path string true "Game ID"
// @Param explanerID path string true "Explaner ID"
// @Success 200 {object} dto.ExplanerResponseDto
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match, game, or Explaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/explaners/{explanerID} [get]
func GetExplanerByID(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")
	explanerID := c.Param("explanerID")

	if teamID == "" || gameID == "" || explanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID, Game ID, and Explaner ID are required"})
		return
	}

	game := models.Game{}
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	explaner := models.Explaner{}
	if err := config.DB.First(&explaner, "explaner_id = ? AND game_id = ?", explanerID, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Explaner not found"})
		return
	}

	var result dto.ExplanerResponseDto

	query := `
		SELECT 
			e.explaner_id, e.game_id, e.team_id, e.hero_id, e.early_result,
			t.team_id AS team_team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image
		FROM explaners e
		JOIN teams t ON e.team_id = t.team_id
		JOIN heros h ON e.hero_id = h.hero_id
		WHERE e.game_id = ? AND e.team_id = ? AND e.explaner_id = ?
	`

	if err := config.DB.Raw(query, gameID, teamID, explanerID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Explaner not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Tags Game
// @Summary Add a goldlaner
// @Description Add a goldlaner for a game with the given game ID and match ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param goldlaner body dto.GoldlanerRequestDto true "Goldlaner data"
// @Success 201 {string} string "Goldlaner added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/goldlaners [post]
func AddGoldlaner(c *gin.Context) {
	gameID := c.Param("gameID")
	teamID := c.Param("teamID")

	if gameID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID and Team ID are required"})
		return
	}

	uintTeamID, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	var game models.Game
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	match := models.Match{}
	if err := config.DB.Where("match_id = ?", game.MatchID).First(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	if match.TeamAID != uint(uintTeamID) && match.TeamBID != uint(uintTeamID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is not part of the match"})
		return
	}

	var input dto.GoldlanerRequestDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var goldlanerExists models.Goldlaner
	if err := config.DB.Where("game_id = ? AND team_id = ? AND hero_id = ?", gameID, teamID, input.HeroID).
		First(&goldlanerExists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Goldlaner already exists"})
		return
	}

	goldlaner := models.Goldlaner{
		GameID:      game.GameID,
		TeamID:      uint(uintTeamID),
		HeroID:      input.HeroID,
		EarlyResult: input.EarlyResult,
	}

	if err := config.DB.Create(&goldlaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Goldlaner added successfully"})
}

// @Tags Game
// @Summary Update a Goldlaner
// @Description Update a Goldlaner with the given game ID, match ID, and Goldlaner ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param goldlanerID path string true "Goldlaner ID"
// @Param goldlaner body dto.GoldlanerRequestDto true "Goldlaner data"
// @Success 200 {object} models.Goldlaner "Goldlaner updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game, Match, or Goldlaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/goldlaners/{goldlanerID} [put]
func UpdateGoldlaner(c *gin.Context) {
	gameID := c.Param("gameID")
	teamID := c.Param("teamID")
	goldlanerID := c.Param("goldlanerID")

	if gameID == "" || teamID == "" || goldlanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID, Team ID, and Goldlaner ID are required"})
		return
	}

	uintTeamID, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	var game models.Game
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	match := models.Match{}
	if err := config.DB.Where("match_id = ?", game.MatchID).First(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	if match.TeamAID != uint(uintTeamID) && match.TeamBID != uint(uintTeamID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is not part of the match"})
		return
	}

	var goldlaner models.Goldlaner
	if err := config.DB.First(&goldlaner, "goldlaner_id = ?", goldlanerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goldlaner not found"})
		return
	}

	var input dto.GoldlanerRequestDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var goldlanerExists models.Goldlaner
	if err := config.DB.Where("game_id = ? AND team_id = ? AND hero_id = ? AND goldlaner_id != ?", gameID, teamID, input.HeroID, goldlanerID).
		First(&goldlanerExists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Goldlaner already exists"})
		return
	}

	goldlaner.TeamID = uint(uintTeamID)
	goldlaner.HeroID = input.HeroID
	goldlaner.EarlyResult = input.EarlyResult

	if err := config.DB.Save(&goldlaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goldlaner)
}

// @Tags Game
// @Summary Delete a Goldlaner
// @Description Delete a Goldlaner with the given game ID, match ID, and Goldlaner ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param goldlanerID path string true "Goldlaner ID"
// @Success 200 {string} string "Goldlaner deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game, Match, or Goldlaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/goldlaners/{goldlanerID} [delete]
func RemoveGoldlaner(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")
	goldlanerID := c.Param("goldlanerID")

	if teamID == "" || gameID == "" || goldlanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID, Game ID, and Goldlaner ID are required"})
		return
	}

	var game models.Game
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	var goldlaner models.Goldlaner
	if err := config.DB.First(&goldlaner, "goldlaner_id = ? AND game_id = ?", goldlanerID, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goldlaner not found"})
		return
	}

	if err := config.DB.Delete(&goldlaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goldlaner deleted successfully"})
}

// @Tags Game
// @Summary Get all Goldlaners for a game
// @Description Get all Goldlaners for a game with the given game ID and match ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param gameID path string true "Game ID"
// @Success 200 {array} dto.GoldlanerResponseDto
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/goldlaners [get]
func GetAllGoldlaners(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")

	if teamID == "" || gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Game ID are required"})
		return
	}

	var game models.Game
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	var results []dto.GoldlanerResponseDto
	query := `
		SELECT 
			g.goldlaner_id, g.game_id, g.team_id, g.hero_id, g.early_result,
			t.team_id AS team_team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image
		FROM goldlaners g
		JOIN teams t ON g.team_id = t.team_id
		JOIN heros h ON g.hero_id = h.hero_id
		WHERE g.game_id = ? AND g.team_id = ?
	`

	if err := config.DB.Raw(query, gameID, teamID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// @Tags Game
// @Summary Get a Goldlaner by ID
// @Description Get a Goldlaner by the given game ID, match ID, and Goldlaner ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param gameID path string true "Game ID"
// @Param goldlanerID path string true "Goldlaner ID"
// @Success 200 {object} dto.GoldlanerResponseDto
// @Failure 404 {string} string "Match or game or Goldlaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/goldlaners/{goldlanerID} [get]
func GetGoldlanerByID(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")
	goldlanerID := c.Param("goldlanerID")

	if teamID == "" || gameID == "" || goldlanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Game ID, and Goldlaner ID are required"})
		return
	}

	game := models.Game{}
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	var result dto.GoldlanerResponseDto

	query := `
		SELECT 
			g.goldlaner_id, g.game_id, g.team_id, g.hero_id, g.early_result,
			t.team_id AS team_team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id AS hero_hero_id, h.name AS hero_name, h.image AS hero_image
		FROM goldlaners g
		JOIN teams t ON g.team_id = t.team_id
		JOIN heros h ON g.hero_id = h.hero_id
		WHERE g.game_id = ? AND g.team_id = ? AND g.goldlaner_id = ?
	`

	if err := config.DB.Raw(query, gameID, teamID, goldlanerID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goldlaner not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Tags Game
// @Summary Add a TrioMid
// @Description Add a TrioMid for a game with the given game ID and match ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param trioMid body dto.TrioMidRequestDto true "TrioMid data"
// @Success 201 {string} string "TrioMid added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/trio-mids [post]
func AddTrioMid(c *gin.Context) {
	gameID, teamID := c.Param("gameID"), c.Param("teamID")

	if gameID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Game ID are required"})
		return
	}

	uintTeamID, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	var game models.Game
	if err := tx.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		tx.Rollback() // Rollback if any error occurs
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	match := models.Match{}
	if err := tx.Where("match_id = ?", game.MatchID).First(&match).Error; err != nil {
		tx.Rollback() // Rollback if any error occurs
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	if match.TeamAID != uint(uintTeamID) && match.TeamBID != uint(uintTeamID) {
		tx.Rollback() // Rollback if any error occurs
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is not part of the match"})
		return
	}

	var input dto.TrioMidRequestDto
	if err := c.ShouldBindJSON(&input); err != nil {
		tx.Rollback() // Rollback if any error occurs
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var trioMid models.TrioMid

	// Cek apakah trioMid ada
	if err := tx.Where("game_id = ? AND team_id = ?", gameID, teamID).First(&trioMid).Error; err != nil {
		// Jika trioMid tidak ada, buat baru dengan EarlyResult kosong
		if err == gorm.ErrRecordNotFound {
			trioMid = models.TrioMid{
				GameID: game.GameID,
				TeamID: uint(uintTeamID),
			}
			if err := tx.Create(&trioMid).Error; err != nil {
				tx.Rollback() // Rollback if any error occurs
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			tx.Rollback() // Rollback if any error occurs
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Cek apakah TrioMidHero sudah ada
	existingTrioMidHero := models.TrioMidHero{}
	if err := tx.Where("trio_mid_id = ? AND hero_id = ?", trioMid.TrioMidID, input.HeroID).First(&existingTrioMidHero).Error; err == nil {
		tx.Rollback() // Rollback if any error occurs
		c.JSON(http.StatusBadRequest, gin.H{"error": "TrioMidHero already exists"})
		return
	}

	// Buat TrioMidHero baru
	trioMidHero := models.TrioMidHero{
		TrioMidID:   trioMid.TrioMidID,
		HeroID:      input.HeroID,
		Role:        input.Role,
		EarlyResult: input.EarlyResult,
	}

	if err := tx.Create(&trioMidHero).Error; err != nil {
		tx.Rollback() // Rollback if any error occurs
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "TrioMid added successfully"})
}

// @Tags Game
// @Summary Update a TrioMid
// @Description Update a TrioMid with the given game ID, match ID, and TrioMid ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param trioMidHeroID path string true "TrioMid ID"
// @Param trioMid body dto.TrioMidRequestDto true "Trio mid data"
// @Success 200 {object} models.TrioMid "Trio mid updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Trio mid not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/trio-mids/{trioMidHeroID} [put]
func UpdateTrioMid(c *gin.Context) {
	gameID, teamID, trioMidHeroID := c.Param("gameID"), c.Param("teamID"), c.Param("trioMidHeroID")

	// Validasi parameter
	if gameID == "" || teamID == "" || trioMidHeroID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID, Team ID, and TrioMidHero ID are required"})
		return
	}

	uintTeamID, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	var game models.Game
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	match := models.Match{}
	if err := config.DB.Where("match_id = ?", game.MatchID).First(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	if match.TeamAID != uint(uintTeamID) && match.TeamBID != uint(uintTeamID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is not part of the match"})
		return
	}

	trioMidHero := models.TrioMidHero{}
	if err := config.DB.Where("trio_mid_hero_id = ?", trioMidHeroID).First(&trioMidHero).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TrioMidHero not found"})
		return
	}

	// Bind input dari JSON request
	var input dto.TrioMidRequestDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingTrioMidHero := models.TrioMidHero{}
	if err := config.DB.Where("trio_mid_id = ? AND hero_id = ? AND trio_mid_hero_id != ?", trioMidHero.TrioMidID, input.HeroID, trioMidHeroID).First(&existingTrioMidHero).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TrioMidHero already exists"})
		return
	}

	trioMidHero.EarlyResult = input.EarlyResult
	trioMidHero.Role = input.Role
	trioMidHero.HeroID = input.HeroID

	if err := config.DB.Save(&trioMidHero).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TrioMid updated successfully"})
}

// @Tags Game
// @Summary Delete a TrioMid
// @Description Delete a TrioMid with the given game ID, match ID, and TrioMid ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param trioMidHeroID path string true "TrioMid ID"
// @Success 200 {string} string "TrioMid deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match, game, or TrioMid not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/trio-mids/{trioMidHeroID} [delete]
func RemoveTrioMid(c *gin.Context) {
	teamID, gameID, trioMidHeroID := c.Param("teamID"), c.Param("gameID"), c.Param("trioMidHeroID")

	// Validasi parameter
	if teamID == "" || gameID == "" || trioMidHeroID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Game ID, and TrioMid ID are required"})
		return
	}

	uintTeamID, err := strconv.ParseUint(teamID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	var game models.Game
	if err := config.DB.Where("game_id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match or game not found"})
		return
	}

	match := models.Match{}
	if err := config.DB.Where("match_id = ?", game.MatchID).First(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	if match.TeamAID != uint(uintTeamID) && match.TeamBID != uint(uintTeamID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is not part of the match"})
		return
	}

	trioMidHero := models.TrioMidHero{}
	if err := config.DB.Where("trio_mid_hero_id = ?", trioMidHeroID).First(&trioMidHero).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TrioMidHero not found"})
		return
	}

	// Cek jumlah TrioMidHero
	var trioMidHeroCount int64
	if err := config.DB.Model(&models.TrioMidHero{}).Where("trio_mid_id = ?", trioMidHero.TrioMidID).Count(&trioMidHeroCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Hapus TrioMidHero atau TrioMid sesuai dengan jumlahnya
	if trioMidHeroCount > 1 {
		// Hapus TrioMidHero
		if err := config.DB.Where("trio_mid_hero_id = ?", trioMidHero.TrioMidHeroID).Delete(&models.TrioMidHero{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "TrioMidHero deleted successfully"})
		return
	}

	// Hapus TrioMidHero
	if err := config.DB.Where("trio_mid_hero_id = ?", trioMidHero.TrioMidHeroID).Delete(&models.TrioMidHero{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var trioMid models.TrioMid
	if err := config.DB.First(&trioMid, "trio_mid_id = ?", trioMidHero.TrioMidID).Error; err != nil {
		config.DB.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "TrioMid not found"})
		return
	}

	if err := config.DB.Delete(&trioMid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TrioMid deleted successfully"})

}

// @Tags Game
// @Summary Get all TrioMids for a game
// @Description Get all TrioMids for a game with the given game ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param gameID path string true "Game ID"
// @Success 200 {array} dto.TrioMidResponseDto
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/trio-mids [get]
func GetAllTrioMids(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")

	if teamID == "" || gameID == "" { // Validasi parameter
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Game ID are required"})
		return
	}

	game := models.Game{}
	if err := config.DB.First(&game, "game_id = ?", gameID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var results = []dto.TrioMidResponseDto{}

	// Query untuk mengambil TrioMidHero dengan JOIN
	query := `
		SELECT 
			tm.trio_mid_id, 
			tm.game_id, 
			tmh.trio_mid_hero_id,
			tmh.role, 
			tmh.early_result,
			t.team_id AS team_team_id, t.name AS team_name, t.image AS team_image,
			th.hero_id AS hero_hero_id, th.name AS hero_name, th.image AS hero_image
		FROM trio_mids tm
		JOIN teams t ON tm.team_id = t.team_id
		JOIN trio_mid_heros tmh ON tm.trio_mid_id = tmh.trio_mid_id
		JOIN heros th ON tmh.hero_id = th.hero_id
		WHERE tm.game_id = ? AND t.team_id = ?
	`

	if err := config.DB.Raw(query, gameID, teamID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Struktur hasil
	for i := range results {
		results[i].Hero = struct {
			HeroID uint   `json:"hero_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}{
			HeroID: results[i].Hero.HeroID,
			Name:   results[i].Hero.Name,
			Image:  results[i].Hero.Image,
		}
		results[i].Team = struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}{
			TeamID: results[i].Team.TeamID,
			Name:   results[i].Team.Name,
			Image:  results[i].Team.Image,
		}
	}

	c.JSON(http.StatusOK, results)
}

// @Tags Game
// @Summary Get a TrioMid by ID
// @Description Get a TrioMid with the given game ID, match ID, and TrioMid ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param trioMidID path string true "TrioMid ID"
// @Success 200 {object} dto.TrioMidResponseDto "Trio mid found successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Trio mid not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/trio-mids/{trioMidID} [get]
func GetTrioMidByID(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")
	trioMidID := c.Param("trioMidID")

	if teamID == "" || gameID == "" { // Validasi parameter
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID and Game ID are required"})
		return
	}

	game := models.Game{}
	if err := config.DB.First(&game, "game_id = ?", gameID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result = dto.TrioMidResponseDto{}

	// Query untuk mengambil TrioMidHero dengan JOIN
	query := `
		SELECT 
			tm.trio_mid_id, 
			tm.game_id, 
			tmh.trio_mid_hero_id,
			tmh.role, 
			tmh.early_result,
			t.team_id AS team_team_id , t.name AS team_name, t.image AS team_image,
			th.hero_id AS hero_hero_id, th.name AS hero_name, th.image AS hero_image
		FROM trio_mids tm
		JOIN teams t ON tm.team_id = t.team_id
		JOIN trio_mid_heros tmh ON tm.trio_mid_id = tmh.trio_mid_id
		JOIN heros th ON tmh.hero_id = th.hero_id
		WHERE tm.trio_mid_id = ? AND tm.game_id = ? AND t.team_id = ?
	`

	if err := config.DB.Raw(query, trioMidID, gameID, teamID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TrioMid not found"})
		return
	}

	// Menyusun data untuk Hero dan Team
	result.Hero = struct {
		HeroID uint   `json:"hero_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}{
		HeroID: result.Hero.HeroID,
		Name:   result.Hero.Name,
		Image:  result.Hero.Image,
	}

	result.Team = struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}{
		TeamID: result.Team.TeamID,
		Name:   result.Team.Name,
		Image:  result.Team.Image,
	}

	c.JSON(http.StatusOK, result)
}

type TrioMidResultDto struct {
	TeamID      uint   `json:"team_id" binding:"required"`
	EarlyResult string `gorm:"type:enum('win', 'draw', 'lose')" json:"early_result"`
}

// @Tags Game
// @Summary Update a TrioMid result
// @Description Update a TrioMid with the given team ID and game ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param teamID path string true "Team ID"
// @Param trioMid body TrioMidResultDto true "Trio mid data"
// @Success 200 {string} string "Trio mid result updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Trio mid not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/trio-mid-results/{trioMidID} [put]
func UpdateTrioMidResult(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")
	trioMidID := c.Param("trioMidID")
	if teamID == "" || gameID == "" || trioMidID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Game ID and Team ID are required"})
		return
	}

	var result TrioMidResultDto
	if err := c.BindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result.EarlyResult != "win" && result.EarlyResult != "draw" && result.EarlyResult != "lose" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid early result"})
	}

	trioMid := models.TrioMid{}
	if err := config.DB.First(&trioMid, "game_id = ? AND team_id = ? AND trio_mid_id = ?", gameID, result.TeamID, trioMidID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TrioMid not found"})
		return
	}

	trioMid.EarlyResult = &result.EarlyResult
	if err := config.DB.Save(&trioMid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TrioMid result updated successfully"})
}

// @Tags Game
// @Summary Get a TrioMid result by ID
// @Description Get a TrioMid result with the given game ID, match ID, and TrioMid ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param gameID path string true "Game ID"
// @Param trioMidID path string true "TrioMid ID"
// @Success 200 {object} models.TrioMid "Trio mid result found successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Trio mid not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/trio-mid-results/{trioMidID} [get]
func GetTrioMidResultByID(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")
	trioMidID := c.Param("trioMidID")

	if teamID == "" || gameID == "" || trioMidID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID, Game ID, and TrioMid ID are required"})
		return
	}

	trioMid := models.TrioMid{}

	if err := config.DB.First(&trioMid, "game_id = ? AND trio_mid_id = ?", gameID, trioMidID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TrioMid not found"})
		return
	}

	c.JSON(http.StatusOK, trioMid)
}

type GameResultDto struct {
	Win    int    `json:"win"`
	Draw   int    `json:"draw"`
	Lose   int    `json:"lose"`
	Result string `json:"result"`
}

// @Tags Game
// @Summary Get all game results for a team in a game
// @Description Get all game results for a team in a game with the given game ID and team ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param gameID path string true "Game ID"
// @Success 200 {array} GameResultDto "All game results found successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /games/{gameID}/teams/{teamID}/game-results [get]
func GetAllGameResults(c *gin.Context) {
	teamID := c.Param("teamID")
	gameID := c.Param("gameID")

	if teamID == "" || gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID and Game ID are required"})
		return
	}

	// Initialize counters for wins, draws, and losses
	winCount := 0
	drawCount := 0
	loseCount := 0

	// Query Explaners
	var explaners []models.Explaner
	if err := config.DB.Where("game_id = ? AND team_id = ?", gameID, teamID).Find(&explaners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Count results for Explaners
	for _, explaner := range explaners {
		switch explaner.EarlyResult {
		case "win":
			winCount++
		case "draw":
			drawCount++
		case "lose":
			loseCount++
		}
	}

	// Query Goldlaners
	var goldlaners []models.Goldlaner
	if err := config.DB.Where("game_id = ? AND team_id = ?", gameID, teamID).Find(&goldlaners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Count results for Goldlaners
	for _, goldlaner := range goldlaners {
		switch goldlaner.EarlyResult {
		case "win":
			winCount++
		case "draw":
			drawCount++
		case "lose":
			loseCount++
		}
	}

	// Query TrioMids
	var trioMids []models.TrioMid
	if err := config.DB.Where("game_id = ? AND team_id = ?", gameID, teamID).Find(&trioMids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Count results for TrioMids
	for _, trioMid := range trioMids {
		if trioMid.EarlyResult != nil {
			switch *trioMid.EarlyResult {
			case "win":
				winCount++
			case "draw":
				drawCount++
			case "lose":
				loseCount++
			}
		}
	}

	// Determine Result based on win, draw, and lose counts
	var result string
	switch {
	case winCount == 3:
		result = "Good Early"
	case winCount == 2 && drawCount == 1:
		result = "Good Early"
	case winCount == 1 && drawCount == 2:
		result = "Good Early"
	case winCount == 2 && loseCount == 1:
		result = "Ok Early"
	case winCount == 1 && drawCount == 1 && loseCount == 1:
		result = "Ok Early"
	case winCount == 1 && loseCount == 2:
		result = "Bad Early"
	case drawCount == 2 && loseCount == 1:
		result = "Bad Early"
	case drawCount == 1 && loseCount == 2:
		result = "Bad Early"
	case loseCount == 3:
		result = "Bad Early"
	default:
		result = "No Result"
	}

	// Prepare the response DTO
	response := GameResultDto{
		Win:    winCount,
		Draw:   drawCount,
		Lose:   loseCount,
		Result: result,
	}

	c.JSON(http.StatusOK, response)
}
