package routes

import (
	"ml-master-data/controllers"
	"ml-master-data/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	r.Use(cors.New(config))

	// Public routes
	r.POST("/login", controllers.Login)

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

		protected.POST("/tournaments/:tournamentID/teams", controllers.CreateTeamInTournament) //ok
		protected.GET("/tournaments/:tournamentID/teams", controllers.GetAllTeamsInTournament)

		protected.GET("/tournaments/:tournamentID/teams/:teamID/matches", controllers.GetAllTeamMatchesinTournament)
		protected.GET("/matches/:matchID", controllers.GetMatchByID)
		protected.POST("/tournaments/:tournamentID/teams/:teamID/matches", controllers.CreateTeamMatchinTournament) //ok
		protected.PUT("/matches/:matchID", controllers.UpdateTeamMatchinTournament)                                 //ok

		protected.POST("matches/:matchID/games", controllers.CreateMatchGame) //ok
		protected.PUT("games/:gameID", controllers.UpdateMatchGame)           //ok
		protected.GET("matches/:matchID/games", controllers.GetAllGameMatches)

		protected.GET("matches/:matchID/teams", controllers.GetAllTeamsInMatch)

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

		protected.POST("matches/:matchID/player-stats", controllers.AddPlayerStatsToMatch) //ok
		protected.PUT("player-stats/:playerStatID", controllers.UpdatePlayerStats)         //ok
		protected.GET("matches/:matchID/player-stats", controllers.GetAllPlayerStatsinMatch)

		protected.POST("matches/:matchID/coach-stats", controllers.AddCoachStatsToMatch) //ok
		protected.PUT("coach-stats/:coachStatID", controllers.UpdateCoachStats)          //ok
		protected.GET("matches/:matchID/coach-stats", controllers.GetAllCoachStatsinMatch)

		protected.POST("matches/:matchID/priority-picks", controllers.AddPriorityPickToMatch) //ok
		protected.PUT("priority-picks/:priorityPickID", controllers.UpdatePriorityPick)       //ok
		protected.GET("matches/:matchID/priority-picks", controllers.GetAllPriorityPicksinMatch)

		protected.POST("matches/:matchID/priority-bans", controllers.AddPriorityBansToMatch) //ok
		protected.PUT("priority-bans/:priorityBanID", controllers.UpdatePriorityBan)         //ok
		protected.GET("matches/:matchID/priority-bans", controllers.GetAllPriorityBansinMatch)

		protected.POST("matches/:matchID/flex-picks", controllers.AddFlexPicksToMatch) //ok
		protected.PUT("flex-picks/:flexPickID", controllers.UpdateFlexPick)            //ok
		protected.GET("matches/:matchID/flex-picks", controllers.GetAllFlexPicksinMatch)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Resource not found"})
	})

	return r
}
