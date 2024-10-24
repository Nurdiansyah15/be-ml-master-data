package routes

import (
	"ml-master-data/controllers"
	"ml-master-data/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "ml-master-data/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	// swagger
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	r.Use(cors.New(config))

	// Public routes
	r.POST("/api/login", controllers.Login)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{

		protected.GET("/me", controllers.Me)
		protected.PUT(("/me"), controllers.UpdateUser)

		protected.GET("/tournaments", controllers.GetAllTournaments)
		protected.GET(("/tournaments/:tournamentID"), controllers.GetTournamentByID)
		protected.POST("/tournaments", controllers.CreateTournament)              //ok
		protected.PUT("/tournaments/:tournamentID", controllers.UpdateTournament) //ok
		protected.DELETE("/tournaments/:tournamentID", controllers.DeleteTournament)

		protected.GET("/tournaments/:tournamentID/matches", controllers.GetMatchesByTournamentID)
		protected.POST("/tournaments/:tournamentID/matches", controllers.CreateTournamentMatch) //ok
		protected.GET("/matches/:matchID", controllers.GetMatchByID)
		protected.PUT("/matches/:matchID", controllers.UpdateMatch) //ok
		protected.DELETE("/matches/:matchID", controllers.DeleteMatch)

		protected.GET("/matches/:matchID/teams", controllers.GetTeamsByMatchID)

		protected.POST("matches/:matchID/teams/:teamID/players", controllers.AddPlayerMatch) //ok
		protected.PUT("matches/:matchID/teams/:teamID/players/:playerID", controllers.UpdatePlayerMatch)
		protected.DELETE("matches/:matchID/teams/:teamID/players/:playerID", controllers.RemovePlayerMatch)
		protected.GET("matches/:matchID/teams/:teamID/players", controllers.GetAllPlayersMatch)

		protected.POST("matches/:matchID/teams/:teamID/coaches", controllers.AddCoachMatch) //ok
		protected.PUT("matches/:matchID/teams/:teamID/coaches/:coachID", controllers.UpdateCoachMatch)
		protected.DELETE("matches/:matchID/teams/:teamID/coaches/:coachID", controllers.RemoveCoachMatch)
		protected.GET("matches/:matchID/teams/:teamID/coaches", controllers.GetAllCoachesMatch)

		protected.POST("matches/:matchID/teams/:teamID/hero-picks", controllers.AddHeroPick) //ok
		protected.DELETE("matches/:matchID/teams/:teamID/hero-picks/:heroPickID", controllers.RemoveHeroPick)
		protected.PUT("matches/:matchID/teams/:teamID/hero-picks/:heroPickID", controllers.UpdateHeroPick)
		protected.GET("matches/:matchID/teams/:teamID/hero-picks", controllers.GetAllHeroPicks)

		protected.GET("matches/:matchID/teams/:teamID/hero-picks-first-phase-more-than-zero", controllers.GetAllHeroPicksWithFirstPhaseMoreThanZero)

		protected.POST("matches/:matchID/teams/:teamID/hero-bans", controllers.AddHeroBan) //ok
		protected.DELETE("matches/:matchID/teams/:teamID/hero-bans/:HeroBanID", controllers.RemoveHeroBan)
		protected.PUT("matches/:matchID/teams/:teamID/hero-bans/:HeroBanID", controllers.UpdateHeroBan)
		protected.GET("matches/:matchID/teams/:teamID/hero-bans", controllers.GetAllHeroBans)

		protected.GET("matches/:matchID/teams/:teamID/hero-bans-first-phase-more-than-zero", controllers.GetAllHeroBansWithFirstPhaseMoreThanZero)

		protected.POST("matches/:matchID/teams/:teamID/priority-picks", controllers.AddPriorityPick) //ok
		protected.DELETE("matches/:matchID/teams/:teamID/priority-picks/:priorityPickID", controllers.RemovePriorityPick)
		protected.PUT("matches/:matchID/teams/:teamID/priority-picks/:priorityPickID", controllers.UpdatePriorityPick)
		protected.GET("matches/:matchID/teams/:teamID/priority-picks/:priorityPickID", controllers.GetPriorityPickByID)
		protected.GET("matches/:matchID/teams/:teamID/priority-picks", controllers.GetAllPriorityPicks)

		protected.POST("matches/:matchID/teams/:teamID/flex-picks", controllers.AddFlexPick) //ok
		protected.DELETE("matches/:matchID/teams/:teamID/flex-picks/:flexPickID", controllers.DeleteFlexPick)
		protected.PUT("matches/:matchID/teams/:teamID/flex-picks/:flexPickID", controllers.UpdateFlexPick)
		protected.GET("matches/:matchID/teams/:teamID/flex-picks/:flexPickID", controllers.GetFlexPickByID)
		protected.GET("matches/:matchID/teams/:teamID/flex-picks", controllers.GetAllFlexPicks)

		protected.POST("/matches/:matchID/teams/:teamID/priority-bans", controllers.AddPriorityBan)
		protected.PUT("/matches/:matchID/teams/:teamID/priority-bans/:priorityBanID", controllers.UpdatePriorityBan)
		protected.GET("/matches/:matchID/teams/:teamID/priority-bans", controllers.GetAllPriorityBans)
		protected.GET("/matches/:matchID/teams/:teamID/priority-bans/:priorityBanID", controllers.GetPriorityBanByID)
		protected.DELETE("/matches/:matchID/teams/:teamID/priority-bans/:priorityBanID", controllers.DeletePriorityBan)

		protected.POST("matches/:matchID/games", controllers.CreateGame)        //ok
		protected.PUT("matches/:matchID/games/:gameID", controllers.UpdateGame) //ok
		protected.GET("matches/:matchID/games", controllers.GetAllGames)
		protected.GET("matches/:matchID/games/:gameID", controllers.GetGameByID)
		protected.DELETE("matches/:matchID/games/:gameID", controllers.RemoveGame)

		protected.POST("matches/:matchID/games/:gameID/lord-results", controllers.AddLordResult)                 //ok
		protected.PUT("matches/:matchID/games/:gameID/lord-results/:lordResultID", controllers.UpdateLordResult) //ok
		protected.GET("matches/:matchID/games/:gameID/lord-results", controllers.GetAllLordResults)
		protected.GET("matches/:matchID/games/:gameID/lord-results/:lordResultID", controllers.GetLordResultByID)
		protected.DELETE("matches/:matchID/games/:gameID/lord-results/:lordResultID", controllers.RemoveLordResult)

		protected.POST("matches/:matchID/games/:gameID/turtle-results", controllers.AddTurtleResult)
		protected.PUT("matches/:matchID/games/:gameID/turtle-results/:turtleResultID", controllers.UpdateTurtleResult)
		protected.GET("matches/:matchID/games/:gameID/turtle-results", controllers.GetAllTurtleResults)
		protected.GET("matches/:matchID/games/:gameID/turtle-results/:turtleResultID", controllers.GetTurtleResultByID)
		protected.DELETE("matches/:matchID/games/:gameID/turtle-results/:turtleResultID", controllers.RemoveTurtleResult)

		protected.POST("games/:gameID/teams/:teamID/explaners", controllers.AddExplaner)
		protected.PUT("games/:gameID/teams/:teamID/explaners/:explanerID", controllers.UpdateExplaner)
		protected.GET("games/:gameID/teams/:teamID/explaners", controllers.GetAllExplaners)
		protected.GET("games/:gameID/teams/:teamID/explaners/:explanerID", controllers.GetExplanerByID)
		protected.DELETE("games/:gameID/teams/:teamID/explaners/:explanerID", controllers.RemoveExplaner)

		protected.POST("games/:gameID/teams/:teamID/goldlaners", controllers.AddGoldlaner)
		protected.PUT("games/:gameID/teams/:teamID/goldlaners/:goldlanerID", controllers.UpdateGoldlaner)
		protected.GET("games/:gameID/teams/:teamID/goldlaners", controllers.GetAllGoldlaners)
		protected.GET("games/:gameID/teams/:teamID/goldlaners/:goldlanerID", controllers.GetGoldlanerByID)
		protected.DELETE("games/:gameID/teams/:teamID/goldlaners/:goldlanerID", controllers.RemoveGoldlaner)

		protected.POST("games/:gameID/teams/:teamID/trio-mids", controllers.AddTrioMid)
		protected.PUT("games/:gameID/teams/:teamID/trio-mids/:trioMidHeroID", controllers.UpdateTrioMid)
		protected.GET("games/:gameID/teams/:teamID/trio-mids", controllers.GetAllTrioMids)
		protected.GET("games/:gameID/teams/:teamID/trio-mids/:trioMidHeroID", controllers.GetTrioMidByID)
		protected.DELETE("games/:gameID/teams/:teamID/trio-mids/:trioMidHeroID", controllers.RemoveTrioMid)

		protected.PUT("games/:gameID/teams/:teamID/trio-mid-results/:trioMidID", controllers.UpdateTrioMidResult)
		protected.GET("games/:gameID/teams/:teamID/trio-mid-results/:trioMidID", controllers.GetTrioMidResultByID)

		protected.GET("games/:gameID/teams/:teamID/game-results", controllers.GetAllGameResults)

		protected.GET("/tournaments/:tournamentID/teams/:teamID/team-statistics", controllers.GetTeamStatistics)
		protected.GET("/teams", controllers.GetAllTeams)
		protected.GET("/teams/:teamID", controllers.GetTeamByID)
		protected.POST("/teams", controllers.CreateTeam)        //ok image ok
		protected.PUT("/teams/:teamID", controllers.UpdateTeam) //ok image ok
		protected.DELETE("/teams/:teamID", controllers.DeleteTeam)

		protected.GET("/tournaments/:tournamentID/coachs/:coachID/coach-statistics", controllers.CoachStatistics)
		protected.GET("teams/:teamID/coaches", controllers.GetAllCoachesInTeam)
		protected.GET("coaches/:coachID", controllers.GetCoachByID)
		protected.POST("teams/:teamID/coaches", controllers.CreateCoachInTeam) //ok image ok
		protected.PUT("coaches/:coachID", controllers.UpdateCoachInTeam)       //ok image ok
		protected.DELETE("coaches/:coachID", controllers.DeleteCoachInTeam)

		protected.GET("/tournaments/:tournamentID/players/:playerID/player-statistics", controllers.PlayerStatistics)
		protected.GET("teams/:teamID/players", controllers.GetAllPlayersInTeam)
		protected.GET("players/:playerID", controllers.GetPlayerByID)
		protected.POST("teams/:teamID/players", controllers.CreatePlayerInTeam) //ok image ok
		protected.PUT("players/:playerID", controllers.UpdatePlayerInTeam)      //ok image ok
		protected.DELETE("players/:playerID", controllers.DeletePlayerInTeam)

		protected.GET("heroes", controllers.GetAllHeroes)
		protected.GET("heroes/:heroID", controllers.GetHeroByID)
		protected.POST("heroes", controllers.CreateHero)        //ok image ok
		protected.PUT("heroes/:heroID", controllers.UpdateHero) //ok image ok
		protected.DELETE("heroes/:heroID", controllers.DeleteHero)

	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Resource not found"})
	})

	return r
}
