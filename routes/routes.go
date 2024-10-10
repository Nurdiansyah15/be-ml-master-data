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
		protected.GET("/tournaments", controllers.GetAllTournaments)
		protected.POST("/tournaments", controllers.CreateTournament)              //ok
		protected.PUT("/tournaments/:tournamentID", controllers.UpdateTournament) //ok
		protected.DELETE("/tournaments/:tournamentID", controllers.DeleteTournament)

		protected.POST("/tournaments/:tournamentID/teams", controllers.CreateTeamInTournament) //ok
		protected.GET("/tournaments/:tournamentID/teams", controllers.GetAllTeamsInTournament)

		protected.GET("/teams", controllers.GetAllTeams)
		protected.POST("/teams", controllers.CreateTeam)        //ok
		protected.PUT("/teams/:teamID", controllers.UpdateTeam) //ok

		protected.GET("/tournaments/:tournamentID/teams/:teamID/matches", controllers.GetAllTeamMatchesinTournament)
		protected.POST("/tournaments/:tournamentID/teams/:teamID/matches", controllers.CreateTeamMatchinTournament) //ok
		protected.PUT("/matches/:matchID", controllers.UpdateTeamMatchinTournament)                                 //ok

		protected.POST("matches/:matchID/games", controllers.CreateMatchGame) //ok
		protected.PUT("games/:gameID", controllers.UpdateMatchGame)           //ok
		protected.GET("matches/:matchID/games", controllers.GetAllGameMatches)

		protected.GET("teams/:teamID/coaches", controllers.GetAllCoachesInTeam)
		protected.POST("teams/:teamID/coaches", controllers.CreateCoachInTeam) //ok
		protected.PUT("coaches/:coachID", controllers.UpdateCoachInTeam)       //ok

		protected.GET("teams/:teamID/players", controllers.GetAllPlayersInTeam)
		protected.POST("teams/:teamID/players", controllers.CreatePlayerInTeam) //ok
		protected.PUT("players/:playerID", controllers.UpdatePlayerInTeam)      //ok

		protected.GET("heroes", controllers.GetAllHeroes)
		protected.POST("heroes", controllers.CreateHero)        //ok
		protected.PUT("heroes/:heroID", controllers.UpdateHero) //ok

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
