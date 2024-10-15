package controllers

import (
	"ml-master-data/config"
	"ml-master-data/dto"
	"ml-master-data/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags Game
// @Summary Create a new game
// @Description Create a new game with the given match ID and additional information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param matchID path string true "Match ID"
// @Param game body dto.GameRequestDto true "Game data"
// @Success 201 {object} models.Game "Game created successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games [post]
func CreateGame(c *gin.Context) {
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

	input := dto.GameRequestDto{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game := models.Game{
		MatchID:          match.MatchID,
		FirstPickTeamID:  input.FirstPickTeamID,
		SecondPickTeamID: input.SecondPickTeamID,
		WinnerTeamID:     input.WinnerTeamID,
		GameNumber:       input.GameNumber,
		VideoLink:        input.VideoLink,
		FullDraftImage:   input.FullDraftImage,
	}

	if err := config.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, game)
}

// @Tags Game
// @Summary Update a game
// @Description Update a game with the given game ID and match ID with the given information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Param game body dto.GameRequestDto true "Game data"
// @Success 200 {object} models.Game "Game updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game not found"
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

	input := dto.GameRequestDto{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
	if input.GameNumber != 0 {
		game.GameNumber = input.GameNumber
	}
	if input.VideoLink != "" {
		game.VideoLink = input.VideoLink
	}
	if input.FullDraftImage != "" {
		game.FullDraftImage = input.FullDraftImage
	}

	if err := config.DB.Save(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	if err := config.DB.Where("match_id = ?", matchID).Find(&models.Game{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var games []dto.GameResponseDto

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
		WHERE g.match_id = ?
	`

	if err := config.DB.Raw(query, matchID).Scan(&games).Error; err != nil {
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
	if err := config.DB.First(&lordResult, "id = ?", lordResultID).Error; err != nil {
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
	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// Validasi keberadaan LordResult
	var lordResult models.LordResult
	if err := config.DB.First(&lordResult, "id = ? AND game_id = ?", lordResultID, gameID).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param gameID path string true "Game ID"
// @Success 200 {array} dto.LordResultResponseDto
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/lord-results [get]
func GetAllLordResults(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")

	// Validasi keberadaan match dan game
	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var results []dto.LordResultResponseDto

	query := `
		SELECT 
			l.lord_result_id, l.game_id, l.team_id, 
			t.team_id, t.name AS team_name, t.image AS team_image,
			l.phase, l.setup, l.initiate, l.result
		FROM lord_results l
		JOIN teams t ON l.team_id = t.team_id
		WHERE l.game_id = ?
	`

	if err := config.DB.Raw(query, gameID).Scan(&results).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param lordResultID path string true "Lord Result ID"
// @Success 200 {object} dto.LordResultResponseDto
// @Failure 404 {string} string "Match or game or Lord result not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/lord-results/{lordResultID} [get]
func GetLordResultByID(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	lordResultID := c.Param("lordResultID")

	// Validasi keberadaan match dan game
	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var result dto.LordResultResponseDto

	query := `
		SELECT 
			l.lord_result_id, l.game_id, l.team_id, 
			t.team_id, t.name AS team_name, t.image AS team_image,
			l.phase, l.setup, l.initiate, l.result
		FROM lord_results l
		JOIN teams t ON l.team_id = t.team_id
		WHERE l.lord_result_id = ? AND l.game_id = ?
	`

	if err := config.DB.Raw(query, lordResultID, gameID).Scan(&result).Error; err != nil {
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

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var turtleResult models.TurtleResult
	if err := config.DB.First(&turtleResult, "id = ?", turtleResultID).Error; err != nil {
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

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var turtleResult models.TurtleResult
	if err := config.DB.First(&turtleResult, "id = ? AND game_id = ?", turtleResultID, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turtle result not found"})
		return
	}

	if err := config.DB.Delete(&turtleResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Turtle result deleted successfully"})
}

func GetAllTurtleResults(c *gin.Context) {

	matchID := c.Param("matchID")
	gameID := c.Param("gameID")

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var results []dto.TurtleResultResponseDto

	query := `
		SELECT 
			tr.turtle_result_id, tr.game_id, tr.team_id, 
			t.team_id AS team_id, t.name AS team_name, t.image AS team_image,
			tr.phase, tr.setup, tr.initiate, tr.result
		FROM turtle_results tr
		JOIN teams t ON tr.team_id = t.team_id
		WHERE tr.game_id = ?
	`

	if err := config.DB.Raw(query, gameID).Scan(&results).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param gameID path string true "Game ID"
// @Param turtleResultID path string true "Turtle result ID"
// @Success 200 {object} dto.TurtleResultResponseDto
// @Failure 404 {string} string "Match, Game, or Turtle result not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/turtle-results/{turtleResultID} [get]
func GetTurtleResultByID(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	turtleResultID := c.Param("turtleResultID")

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var result dto.TurtleResultResponseDto

	query := `
		SELECT 
			tr.turtle_result_id, tr.game_id, tr.team_id, 
			t.team_id AS team_id, t.name AS team_name, t.image AS team_image,
			tr.phase, tr.setup, tr.initiate, tr.result
		FROM turtle_results tr
		JOIN teams t ON tr.team_id = t.team_id
		WHERE tr.turtle_result_id = ? AND tr.game_id = ?
	`

	if err := config.DB.Raw(query, turtleResultID, gameID).Scan(&result).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param explaner body dto.ExplanerRequestDto true "Explaner data"
// @Success 201 {string} string "Explaner added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/explaners [post]
func AddExplaner(c *gin.Context) {
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

	var input dto.ExplanerRequestDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	explaner := models.Explaner{
		GameID:      game.GameID,
		TeamID:      input.TeamID,
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
// @Param matchID path string true "Match ID"
// @Param explanerID path string true "Explaner ID"
// @Param explaner body dto.ExplanerRequestDto true "Explaner data"
// @Success 200 {object} models.Explaner "Explaner updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Explaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/explaners/{explanerID} [put]
func UpdateExplaner(c *gin.Context) {
	gameID := c.Param("gameID")
	matchID := c.Param("matchID")
	explanerID := c.Param("explanerID")

	if gameID == "" || matchID == "" || explanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID, Match ID, and Explaner ID are required"})
		return
	}

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var explaner models.Explaner
	if err := config.DB.First(&explaner, "id = ?", explanerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Explaner not found"})
		return
	}

	input := dto.ExplanerRequestDto{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	explaner.TeamID = input.TeamID
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
// @Param matchID path string true "Match ID"
// @Param explanerID path string true "Explaner ID"
// @Success 200 {string} string "Explaner deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Explaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/explaners/{explanerID} [delete]
func RemoveExplaner(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	explanerID := c.Param("explanerID")

	if matchID == "" || gameID == "" || explanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Game ID, and Explaner ID are required"})
		return
	}

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var explaner models.Explaner
	if err := config.DB.First(&explaner, "id = ? AND game_id = ?", explanerID, gameID).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param gameID path string true "Game ID"
// @Success 200 {array} dto.ExplanerResponseDto
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/explaners [get]
func GetAllExplaners(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var results []dto.ExplanerResponseDto

	query := `
		SELECT 
			e.explaner_id, e.game_id, e.team_id, e.hero_id, e.early_result,
			t.team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id, h.name AS hero_name, h.image AS hero_image
		FROM explaners e
		JOIN teams t ON e.team_id = t.team_id
		JOIN heroes h ON e.hero_id = h.hero_id
		WHERE e.game_id = ?
	`

	if err := config.DB.Raw(query, gameID).Scan(&results).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param gameID path string true "Game ID"
// @Param explanerID path string true "Explaner ID"
// @Success 200 {object} dto.ExplanerResponseDto
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match, game, or Explaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/explaners/{explanerID} [get]
func GetExplanerByID(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	explanerID := c.Param("explanerID")

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var result dto.ExplanerResponseDto

	query := `
		SELECT 
			e.explaner_id, e.game_id, e.team_id, e.hero_id, e.early_result,
			t.team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id, h.name AS hero_name, h.image AS hero_image
		FROM explaners e
		JOIN teams t ON e.team_id = t.team_id
		JOIN heroes h ON e.hero_id = h.hero_id
		WHERE e.explaner_id = ? AND e.game_id = ?
	`

	if err := config.DB.Raw(query, explanerID, gameID).Scan(&result).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param goldlaner body dto.GoldlanerRequestDto true "Goldlaner data"
// @Success 201 {string} string "Goldlaner added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/goldlaners [post]
func AddGoldlaner(c *gin.Context) {
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

	var input dto.GoldlanerRequestDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var goldlaner models.Goldlaner
	goldlaner.TeamID = input.TeamID
	goldlaner.HeroID = input.HeroID
	goldlaner.EarlyResult = input.EarlyResult
	goldlaner.GameID = game.GameID

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
// @Param matchID path string true "Match ID"
// @Param goldlanerID path string true "Goldlaner ID"
// @Param goldlaner body dto.GoldlanerRequestDto true "Goldlaner data"
// @Success 200 {object} models.Goldlaner "Goldlaner updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game, Match, or Goldlaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/goldlaners/{goldlanerID} [put]
func UpdateGoldlaner(c *gin.Context) {
	gameID := c.Param("gameID")
	matchID := c.Param("matchID")
	goldlanerID := c.Param("goldlanerID")

	if gameID == "" || matchID == "" || goldlanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID, Match ID, and Goldlaner ID are required"})
		return
	}

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var goldlaner models.Goldlaner
	if err := config.DB.First(&goldlaner, "id = ?", goldlanerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goldlaner not found"})
		return
	}

	input := dto.GoldlanerRequestDto{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goldlaner.TeamID = input.TeamID
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
// @Param matchID path string true "Match ID"
// @Param goldlanerID path string true "Goldlaner ID"
// @Success 200 {string} string "Goldlaner deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game, Match, or Goldlaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/goldlaners/{goldlanerID} [delete]
func RemoveGoldlaner(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	goldlanerID := c.Param("goldlanerID")

	if matchID == "" || gameID == "" || goldlanerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Game ID, and Goldlaner ID are required"})
		return
	}

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var goldlaner models.Goldlaner
	if err := config.DB.First(&goldlaner, "id = ? AND game_id = ?", goldlanerID, gameID).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param gameID path string true "Game ID"
// @Success 200 {array} dto.GoldlanerResponseDto
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/goldlaners [get]
func GetAllGoldlaners(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var results []dto.GoldlanerResponseDto

	query := `
		SELECT 
			g.goldlaner_id, g.game_id, g.team_id, g.hero_id, g.early_result,
			t.team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id, h.name AS hero_name, h.image AS hero_image
		FROM goldlaners g
		JOIN teams t ON g.team_id = t.team_id
		JOIN heroes h ON g.hero_id = h.hero_id
		WHERE g.game_id = ?
	`

	if err := config.DB.Raw(query, gameID).Scan(&results).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param gameID path string true "Game ID"
// @Param goldlanerID path string true "Goldlaner ID"
// @Success 200 {object} dto.GoldlanerResponseDto
// @Failure 404 {string} string "Match or game or Goldlaner not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/goldlaners/{goldlanerID} [get]
func GetGoldlanerByID(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	goldlanerID := c.Param("goldlanerID")

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var result dto.GoldlanerResponseDto

	query := `
		SELECT 
			g.goldlaner_id, g.game_id, g.team_id, g.hero_id, g.early_result,
			t.team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id, h.name AS hero_name, h.image AS hero_image
		FROM goldlaners g
		JOIN teams t ON g.team_id = t.team_id
		JOIN heroes h ON g.hero_id = h.hero_id
		WHERE g.goldlaner_id = ? AND g.game_id = ?
	`

	if err := config.DB.Raw(query, goldlanerID, gameID).Scan(&result).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param trioMid body dto.TrioMidRequestDto true "TrioMid data"
// @Success 201 {string} string "TrioMid added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match or game not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/trio-mids [post]
func AddTrioMid(c *gin.Context) {
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

	var input dto.TrioMidRequestDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trioMid := models.TrioMid{
		GameID:      game.GameID,
		TeamID:      input.TeamID,
		HeroID:      input.HeroID,
		Role:        input.Role,
		EarlyResult: input.EarlyResult,
	}

	if err := config.DB.Create(&trioMid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Param matchID path string true "Match ID"
// @Param trioMidID path string true "TrioMid ID"
// @Param trioMid body dto.TrioMidRequestDto true "Trio mid data"
// @Success 200 {object} models.TrioMid "Trio mid updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Trio mid not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/trio-mids/{trioMidID} [put]
func UpdateTrioMid(c *gin.Context) {
	gameID := c.Param("gameID")
	matchID := c.Param("matchID")
	trioMidID := c.Param("trioMidID")

	if gameID == "" || matchID == "" || trioMidID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID, Match ID, and TrioMid ID are required"})
		return
	}

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var trioMid models.TrioMid
	if err := config.DB.First(&trioMid, "id = ?", trioMidID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TrioMid not found"})
		return
	}

	input := dto.TrioMidRequestDto{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trioMid.TeamID = input.TeamID
	trioMid.HeroID = input.HeroID
	trioMid.Role = input.Role
	trioMid.EarlyResult = input.EarlyResult

	if err := config.DB.Save(&trioMid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trioMid)
}

// @Tags Game
// @Summary Delete a TrioMid
// @Description Delete a TrioMid with the given game ID, match ID, and TrioMid ID
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param gameID path string true "Game ID"
// @Param matchID path string true "Match ID"
// @Param trioMidID path string true "TrioMid ID"
// @Success 200 {string} string "TrioMid deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Match, game, or TrioMid not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/trio-mids/{trioMidID} [delete]
func RemoveTrioMid(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	trioMidID := c.Param("trioMidID")

	if matchID == "" || gameID == "" || trioMidID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match ID, Game ID, and TrioMid ID are required"})
		return
	}

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var trioMid models.TrioMid
	if err := config.DB.First(&trioMid, "id = ? AND game_id = ?", trioMidID, gameID).Error; err != nil {
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
// @Param matchID path string true "Match ID"
// @Param gameID path string true "Game ID"
// @Success 200 {array} dto.TrioMidResponseDto
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/trio-mids [get]
func GetAllTrioMids(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var results []dto.TrioMidResponseDto

	query := `
		SELECT 
			tm.trio_mid_id, tm.game_id, tm.team_id, tm.hero_id, tm.early_result,
			t.team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id, h.name AS hero_name, h.image AS hero_image
		FROM trio_mids tm
		JOIN teams t ON tm.team_id = t.team_id
		JOIN heroes h ON tm.hero_id = h.hero_id
		WHERE tm.game_id = ?
	`

	if err := config.DB.Raw(query, gameID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
// @Param matchID path string true "Match ID"
// @Param trioMidID path string true "TrioMid ID"
// @Success 200 {object} dto.TrioMidResponseDto "Trio mid found successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Game or Trio mid not found"
// @Failure 500 {string} string "Internal server error"
// @Router /matches/{matchID}/games/{gameID}/trio-mids/{trioMidID} [get]
func GetTrioMidByID(c *gin.Context) {
	matchID := c.Param("matchID")
	gameID := c.Param("gameID")
	trioMidID := c.Param("trioMidID")

	if err := config.DB.First(&models.Match{}, "id = ?", matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.First(&models.Game{}, "id = ? AND match_id = ?", gameID, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var result dto.TrioMidResponseDto

	query := `
		SELECT 
			tm.trio_mid_id, tm.game_id, tm.team_id, tm.hero_id, tm.early_result,
			t.team_id, t.name AS team_name, t.image AS team_image,
			h.hero_id, h.name AS hero_name, h.image AS hero_image
		FROM trio_mids tm
		JOIN teams t ON tm.team_id = t.team_id
		JOIN heroes h ON tm.hero_id = h.hero_id
		WHERE tm.trio_mid_id = ? AND tm.game_id = ?
	`

	if err := config.DB.Raw(query, trioMidID, gameID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TrioMid not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}
