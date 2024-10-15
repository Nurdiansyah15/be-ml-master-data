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
	r.Use(cors.New(config))

	// Public routes
	r.POST("/api/login", controllers.Login)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{

		protected.GET("/me", controllers.Me)

		protected.GET("/tournaments", controllers.GetAllTournaments)
		protected.GET(("/tournaments/:tournamentID"), controllers.GetTournamentByID)
		protected.POST("/tournaments", controllers.CreateTournament)              //ok
		protected.PUT("/tournaments/:tournamentID", controllers.UpdateTournament) //ok
		protected.DELETE("/tournaments/:tournamentID", controllers.DeleteTournament)

		// protected.POST("/tournaments/:tournamentID/teams", controllers.CreateTeamInTournament) //ok
		// protected.GET("/tournaments/:tournamentID/teams", controllers.GetAllTeamsInTournament)

		protected.GET("/tournaments/:tournamentID/matches", controllers.GetMatchesByTournamentID)
		protected.POST("/tournaments/:tournamentID/matches", controllers.CreateTournamentMatch) //ok
		protected.GET("/matches/:matchID", controllers.GetMatchByID)
		protected.PUT("/matches/:matchID", controllers.UpdateMatch) //ok

		protected.POST("matches/:matchID/teams/:teamID/players", controllers.AddPlayerMatch) //ok
		protected.DELETE("matches/:matchID/teams/:teamID/players/:playerID", controllers.RemovePlayerMatch)
		protected.GET("matches/:matchID/teams/:teamID/players", controllers.GetAllPlayersMatch)

		protected.POST("matches/:matchID/teams/:teamID/coaches", controllers.AddCoachMatch) //ok
		protected.DELETE("matches/:matchID/teams/:teamID/coaches/:coachID", controllers.RemoveCoachMatch)
		protected.GET("matches/:matchID/teams/:teamID/coaches", controllers.GetAllCoachesMatch)

		protected.POST("matches/:matchID/teams/:teamID/hero-picks", controllers.AddHeroPick) //ok
		protected.DELETE("matches/:matchID/teams/:teamID/hero-picks/:heroPickID", controllers.RemoveHeroPick)
		protected.PUT("matches/:matchID/teams/:teamID/hero-picks/:heroPickID", controllers.UpdateHeroPick)
		protected.GET("matches/:matchID/teams/:teamID/hero-picks", controllers.GetAllHeroPicks)

		protected.POST("matches/:matchID/teams/:teamID/hero-bans", controllers.AddHeroBan) //ok
		protected.DELETE("matches/:matchID/teams/:teamID/hero-bans/:HeroBanID", controllers.RemoveHeroBan)
		protected.PUT("matches/:matchID/teams/:teamID/hero-bans/:HeroBanID", controllers.UpdateHeroBan)
		protected.GET("matches/:matchID/teams/:teamID/hero-bans", controllers.GetAllHeroBans)

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

		protected.POST("matches/:matchID/games/:gameID/lord-result", controllers.AddLordResult)                 //ok
		protected.PUT("matches/:matchID/games/:gameID/lord-result/:lordResultID", controllers.UpdateLordResult) //ok
		protected.GET("matches/:matchID/games/:gameID/lord-result", controllers.GetAllLordResults)
		protected.GET("matches/:matchID/games/:gameID/lord-result/:lordResultID", controllers.GetLordResultByID)
		protected.DELETE("matches/:matchID/games/:gameID/lord-result/:lordResultID", controllers.RemoveLordResult)

		protected.POST("matches/:matchID/games/:gameID/turtle-result", controllers.AddTurtleResult)
		protected.PUT("matches/:matchID/games/:gameID/turtle-result/:turtleResultID", controllers.UpdateTurtleResult)
		protected.GET("matches/:matchID/games/:gameID/turtle-result", controllers.GetAllTurtleResults)
		protected.GET("matches/:matchID/games/:gameID/turtle-result/:turtleResultID", controllers.GetTurtleResultByID)
		protected.DELETE("matches/:matchID/games/:gameID/turtle-result/:turtleResultID", controllers.RemoveTurtleResult)

		protected.POST("matches/:matchID/games/:gameID/explaner", controllers.AddExplaner)
		protected.PUT("matches/:matchID/games/:gameID/explaner/:explanerID", controllers.UpdateExplaner)
		protected.GET("matches/:matchID/games/:gameID/explaner", controllers.GetAllExplaners)
		protected.GET("matches/:matchID/games/:gameID/explaner/:explanerID", controllers.GetExplanerByID)
		protected.DELETE("matches/:matchID/games/:gameID/explaner/:explanerID", controllers.RemoveExplaner)

		protected.POST("matches/:matchID/games/:gameID/goldlaner", controllers.AddGoldlaner)
		protected.PUT("matches/:matchID/games/:gameID/goldlaner/:goldlanerID", controllers.UpdateGoldlaner)
		protected.GET("matches/:matchID/games/:gameID/goldlaner", controllers.GetAllGoldlaners)
		protected.GET("matches/:matchID/games/:gameID/goldlaner/:goldlanerID", controllers.GetGoldlanerByID)
		protected.DELETE("matches/:matchID/games/:gameID/goldlaner/:goldlanerID", controllers.RemoveGoldlaner)

		protected.POST("matches/:matchID/games/:gameID/trio_mid", controllers.AddTrioMid)
		protected.PUT("matches/:matchID/games/:gameID/trio_mid/:trioMidID", controllers.UpdateTrioMid)
		protected.GET("matches/:matchID/games/:gameID/trio_mid", controllers.GetAllTrioMids)
		protected.GET("matches/:matchID/games/:gameID/trio_mid/:trioMidID", controllers.GetTrioMidByID)
		protected.DELETE("matches/:matchID/games/:gameID/trio_mid/:trioMidID", controllers.RemoveTrioMid)

		// protected.GET("matches/:matchID/teams", controllers.GetAllTeamsInMatch)

		protected.GET("/teams", controllers.GetAllTeams)
		protected.GET("/teams/:teamID", controllers.GetTeamByID)
		protected.POST("/teams", controllers.CreateTeam)        //ok image ok
		protected.PUT("/teams/:teamID", controllers.UpdateTeam) //ok image ok

		protected.GET("teams/:teamID/coaches", controllers.GetAllCoachesInTeam)
		protected.GET("coaches/:coachID", controllers.GetCoachByID)
		protected.POST("teams/:teamID/coaches", controllers.CreateCoachInTeam) //ok image ok
		protected.PUT("coaches/:coachID", controllers.UpdateCoachInTeam)       //ok image ok

		protected.GET("teams/:teamID/players", controllers.GetAllPlayersInTeam)
		protected.GET("players/:playerID", controllers.GetPlayerByID)
		protected.POST("teams/:teamID/players", controllers.CreatePlayerInTeam) //ok image ok
		protected.PUT("players/:playerID", controllers.UpdatePlayerInTeam)      //ok image ok

		protected.GET("heroes", controllers.GetAllHeroes)
		protected.GET("heroes/:heroID", controllers.GetHeroByID)
		protected.POST("heroes", controllers.CreateHero)        //ok image ok
		protected.PUT("heroes/:heroID", controllers.UpdateHero) //ok image ok

		// protected.POST("matches/:matchID/player-stats", controllers.AddPlayerStatsToMatch) //ok
		// protected.PUT("player-stats/:playerStatID", controllers.UpdatePlayerStats)         //ok
		// protected.GET("matches/:matchID/player-stats", controllers.GetAllPlayerStatsinMatch)

		// protected.POST("matches/:matchID/coach-stats", controllers.AddCoachStatsToMatch) //ok
		// protected.PUT("coach-stats/:coachStatID", controllers.UpdateCoachStats)          //ok
		// protected.GET("matches/:matchID/coach-stats", controllers.GetAllCoachStatsinMatch)

		// protected.POST("matches/:matchID/priority-picks", controllers.AddPriorityPickToMatch) //ok
		// protected.PUT("priority-picks/:priorityPickID", controllers.UpdatePriorityPick)       //ok
		// protected.GET("matches/:matchID/priority-picks", controllers.GetAllPriorityPicksinMatch)

		// protected.POST("matches/:matchID/priority-bans", controllers.AddPriorityBansToMatch) //ok
		// protected.PUT("priority-bans/:priorityBanID", controllers.UpdatePriorityBan)         //ok
		// protected.GET("matches/:matchID/priority-bans", controllers.GetAllPriorityBansinMatch)

		// protected.POST("matches/:matchID/flex-picks", controllers.AddFlexPicksToMatch) //ok
		// protected.PUT("flex-picks/:flexPickID", controllers.UpdateFlexPick)            //ok
		// protected.GET("matches/:matchID/flex-picks", controllers.GetAllFlexPicksinMatch)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Resource not found"})
	})

	return r
}
