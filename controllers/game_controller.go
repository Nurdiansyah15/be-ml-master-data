package controllers

import (
	"ml-master-data/config"
	"ml-master-data/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	input := struct {
		FirstPickTeamID  uint   `json:"first_pick_team_id" binding:"required"`
		SecondPickTeamID uint   `json:"second_pick_team_id" binding:"required"`
		WinnerTeamID     uint   `json:"winner_team_id" binding:"required"`
		GameNumber       int    `json:"game_number" binding:"required"`
		VideoLink        string `json:"video_link"`
		FullDraftImage   string `json:"full_draft_image"`
	}{}

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

	input := struct {
		FirstPickTeamID  uint   `json:"first_pick_team_id"`
		SecondPickTeamID uint   `json:"second_pick_team_id"`
		WinnerTeamID     uint   `json:"winner_team_id"`
		GameNumber       int    `json:"game_number"`
		VideoLink        string `json:"video_link"`
		FullDraftImage   string `json:"full_draft_image"`
	}{}

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

	type GameResponse struct {
		GameID          uint `json:"game_id"`
		MatchID         uint `json:"match_id"`
		FirstPickTeamID uint `json:"first_pick_team_id"`
		FirstTeam       struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		SecondPickTeamID uint `json:"second_pick_team_id"`
		SecondTeam       struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		WinnerTeamID uint `json:"winner_team_id"`
		WinnerTeam   struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		GameNumber     int    `json:"game_number"`
		VideoLink      string `json:"video_link"`
		FullDraftImage string `json:"full_draft_image"`
	}

	var games []GameResponse

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

func GetGameByID(c *gin.Context) {
	gameID := c.Param("gameID")
	matchID := c.Param("matchID")

	if gameID == "" || matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID and Match ID are required"})
		return
	}

	type GameResponse struct {
		GameID          uint `json:"game_id"`
		MatchID         uint `json:"match_id"`
		FirstPickTeamID uint `json:"first_pick_team_id"`
		FirstTeam       struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		SecondPickTeamID uint `json:"second_pick_team_id"`
		SecondTeam       struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		WinnerTeamID uint `json:"winner_team_id"`
		WinnerTeam   struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		GameNumber     int    `json:"game_number"`
		VideoLink      string `json:"video_link"`
		FullDraftImage string `json:"full_draft_image"`
	}

	var game GameResponse

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

	var lordResult models.LordResult

	if err := c.ShouldBindJSON(&lordResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lordResult.GameID = game.GameID

	if err := config.DB.Create(&lordResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Lord result added successfully"})
}

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
	var input struct {
		TeamID   *uint   `json:"team_id"`
		Phase    *string `json:"phase"`
		Setup    *string `json:"setup"`
		Initiate *string `json:"initiate"`
		Result   *string `json:"result"`
	}

	// Validasi input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update hanya field yang tidak bernilai null
	if input.TeamID != nil {
		lordResult.TeamID = *input.TeamID
	}
	if input.Phase != nil {
		lordResult.Phase = *input.Phase
	}
	if input.Setup != nil {
		lordResult.Setup = *input.Setup
	}
	if input.Initiate != nil {
		lordResult.Initiate = *input.Initiate
	}
	if input.Result != nil {
		lordResult.Result = *input.Result
	}

	// Simpan perubahan ke database
	if err := config.DB.Save(&lordResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lordResult)
}

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

	type LordResultResponse struct {
		LordResultID uint `json:"lord_result_id"`
		GameID       uint `json:"game_id"`
		TeamID       uint `json:"team_id"`
		Team         struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		Phase    string `json:"phase"`
		Setup    string `json:"setup"`
		Initiate string `json:"initiate"`
		Result   string `json:"result"`
	}

	var results []LordResultResponse

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

	type LordResultResponse struct {
		LordResultID uint `json:"lord_result_id"`
		GameID       uint `json:"game_id"`
		TeamID       uint `json:"team_id"`
		Team         struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		}
		Phase    string `json:"phase"`
		Setup    string `json:"setup"`
		Initiate string `json:"initiate"`
		Result   string `json:"result"`
	}

	var result LordResultResponse

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

	var turtleResult models.TurtleResult
	if err := c.ShouldBindJSON(&turtleResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	turtleResult.GameID = game.GameID

	if err := config.DB.Create(&turtleResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Turtle result added successfully"})
}

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

	var input struct {
		TeamID   *uint   `json:"team_id"`
		Phase    *string `json:"phase"`
		Setup    *string `json:"setup"`
		Initiate *string `json:"initiate"`
		Result   *string `json:"result"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.TeamID != nil {
		turtleResult.TeamID = *input.TeamID
	}
	if input.Phase != nil {
		turtleResult.Phase = *input.Phase
	}
	if input.Setup != nil {
		turtleResult.Setup = *input.Setup
	}
	if input.Initiate != nil {
		turtleResult.Initiate = *input.Initiate
	}
	if input.Result != nil {
		turtleResult.Result = *input.Result
	}

	if err := config.DB.Save(&turtleResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, turtleResult)
}

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

	type TurtleResultResponse struct {
		TurtleResultID uint `json:"turtle_result_id"`
		GameID         uint `json:"game_id"`
		TeamID         uint `json:"team_id"`
		Team           struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		} `json:"team"`
		Phase    string `json:"phase"`
		Setup    string `json:"setup"`
		Initiate string `json:"initiate"`
		Result   string `json:"result"`
	}

	var results []TurtleResultResponse

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

	type TurtleResultResponse struct {
		TurtleResultID uint `json:"turtle_result_id"`
		GameID         uint `json:"game_id"`
		TeamID         uint `json:"team_id"`
		Team           struct {
			TeamID uint   `json:"team_id"`
			Name   string `json:"name"`
			Image  string `json:"image"`
		} `json:"team"`
		Phase    string `json:"phase"`
		Setup    string `json:"setup"`
		Initiate string `json:"initiate"`
		Result   string `json:"result"`
	}

	var result TurtleResultResponse

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

	var explaner models.Explaner
	if err := c.ShouldBindJSON(&explaner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	explaner.GameID = game.GameID

	if err := config.DB.Create(&explaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Explaner added successfully"})
}

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

	var input struct {
		TeamID      *uint   `json:"team_id"`
		HeroID      *uint   `json:"hero_id"`
		EarlyResult *string `json:"early_result"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.TeamID != nil {
		explaner.TeamID = *input.TeamID
	}
	if input.HeroID != nil {
		explaner.HeroID = *input.HeroID
	}
	if input.EarlyResult != nil {
		explaner.EarlyResult = *input.EarlyResult
	}

	if err := config.DB.Save(&explaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, explaner)
}

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

	type TeamResponse struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type HeroResponse struct {
		HeroID uint   `json:"hero_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type ExplanerResponse struct {
		ExplanerID  uint         `json:"explaner_id"`
		GameID      uint         `json:"game_id"`
		TeamID      uint         `json:"team_id"`
		Team        TeamResponse `json:"team"`
		HeroID      uint         `json:"hero_id"`
		Hero        HeroResponse `json:"hero"`
		EarlyResult string       `json:"early_result"`
	}

	var results []ExplanerResponse

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

	type TeamResponse struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type HeroResponse struct {
		HeroID uint   `json:"hero_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type ExplanerResponse struct {
		ExplanerID  uint         `json:"explaner_id"`
		GameID      uint         `json:"game_id"`
		TeamID      uint         `json:"team_id"`
		Team        TeamResponse `json:"team"`
		HeroID      uint         `json:"hero_id"`
		Hero        HeroResponse `json:"hero"`
		EarlyResult string       `json:"early_result"`
	}

	var result ExplanerResponse

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

	var goldlaner models.Goldlaner
	if err := c.ShouldBindJSON(&goldlaner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goldlaner.GameID = game.GameID

	if err := config.DB.Create(&goldlaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Goldlaner added successfully"})
}

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

	var input struct {
		TeamID      *uint   `json:"team_id"`
		HeroID      *uint   `json:"hero_id"`
		EarlyResult *string `json:"early_result"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.TeamID != nil {
		goldlaner.TeamID = *input.TeamID
	}
	if input.HeroID != nil {
		goldlaner.HeroID = *input.HeroID
	}
	if input.EarlyResult != nil {
		goldlaner.EarlyResult = *input.EarlyResult
	}

	if err := config.DB.Save(&goldlaner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goldlaner)
}

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

	type TeamResponse struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type HeroResponse struct {
		HeroID uint   `json:"hero_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type GoldlanerResponse struct {
		GoldlanerID uint         `json:"goldlaner_id"`
		GameID      uint         `json:"game_id"`
		TeamID      uint         `json:"team_id"`
		Team        TeamResponse `json:"team"`
		HeroID      uint         `json:"hero_id"`
		Hero        HeroResponse `json:"hero"`
		EarlyResult string       `json:"early_result"`
	}

	var results []GoldlanerResponse

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

	type TeamResponse struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type HeroResponse struct {
		HeroID uint   `json:"hero_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type GoldlanerResponse struct {
		GoldlanerID uint         `json:"goldlaner_id"`
		GameID      uint         `json:"game_id"`
		TeamID      uint         `json:"team_id"`
		Team        TeamResponse `json:"team"`
		HeroID      uint         `json:"hero_id"`
		Hero        HeroResponse `json:"hero"`
		EarlyResult string       `json:"early_result"`
	}

	var result GoldlanerResponse

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

	var trioMid models.TrioMid
	if err := c.ShouldBindJSON(&trioMid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trioMid.GameID = game.GameID

	if err := config.DB.Create(&trioMid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "TrioMid added successfully"})
}

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

	var input struct {
		TeamID      *uint   `json:"team_id"`
		HeroID      *uint   `json:"hero_id"`
		Role        *string `json:"role"`
		EarlyResult *string `json:"early_result"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.TeamID != nil {
		trioMid.TeamID = *input.TeamID
	}
	if input.HeroID != nil {
		trioMid.HeroID = *input.HeroID
	}
	if input.Role != nil {
		trioMid.Role = *input.Role
	}
	if input.EarlyResult != nil {
		trioMid.EarlyResult = *input.EarlyResult
	}

	if err := config.DB.Save(&trioMid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trioMid)
}

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

	type TeamResponse struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type HeroResponse struct {
		HeroID uint   `json:"hero_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type TrioMidResponse struct {
		TrioMidID   uint         `json:"trio_mid_id"`
		GameID      uint         `json:"game_id"`
		TeamID      uint         `json:"team_id"`
		Team        TeamResponse `json:"team"`
		HeroID      uint         `json:"hero_id"`
		Hero        HeroResponse `json:"hero"`
		EarlyResult string       `json:"early_result"`
	}

	var results []TrioMidResponse

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

	type TeamResponse struct {
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type HeroResponse struct {
		HeroID uint   `json:"hero_id"`
		Name   string `json:"name"`
		Image  string `json:"image"`
	}

	type TrioMidResponse struct {
		TrioMidID   uint         `json:"trio_mid_id"`
		GameID      uint         `json:"game_id"`
		TeamID      uint         `json:"team_id"`
		Team        TeamResponse `json:"team"`
		HeroID      uint         `json:"hero_id"`
		Hero        HeroResponse `json:"hero"`
		EarlyResult string       `json:"early_result"`
	}

	var result TrioMidResponse

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
