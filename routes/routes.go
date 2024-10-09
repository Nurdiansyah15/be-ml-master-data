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
		protected.POST("/tournaments", controllers.CreateTournament)
		protected.PUT("/tournaments/:tournamentID", controllers.UpdateTournament)
		protected.DELETE("/tournaments/:tournamentID", controllers.DeleteTournament)

		protected.POST("/tournaments/:tournamentID/teams", controllers.CreateTeamInTournament)
		protected.GET("/tournaments/:tournamentID/teams", controllers.GetAllTeamsInTournament)

		protected.GET("/teams", controllers.GetAllTeams)
		protected.POST("/teams", controllers.CreateTeam)
		protected.PUT("/teams/:teamID", controllers.UpdateTeam)
	}

	return r
}
