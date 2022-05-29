package routes

import (
	"SADBackend/controllers"
	v1 "SADBackend/controllers/v1"
	_ "SADBackend/docs"
	"SADBackend/repo"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func initConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Upgrade", "Connection", "Accept", "Accept-Encoding", "Accept-Language", "Host", "Cookie", "Referer", "User-Agent"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Token"}
	config.AllowWildcard = true
	config.AllowOriginFunc = func(origin string) bool {
		allowOriginDomains := []string{"localhost", "127.0.0.1"}
		for _, d := range allowOriginDomains {
			if strings.Contains(origin, d) {
				return true
			}
		}
		return false
	}
	return config
}

func InitRouters(dbInstance repo.AllRepo) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := initConfig()
	router.Use(cors.New(config))

	// swag
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/heartBeat", controllers.HeartBeat)
	apiV1 := router.Group("/api/v1")
	{
		gymAPI := apiV1.Group("/gym")
		{
			gymAPI.GET("/list", func(c *gin.Context) {
				v1.GetGymList(c, dbInstance.Gym)
			})
			gymAPI.GET("/machine", func(c *gin.Context) {
				v1.GetMachineList(c, dbInstance.Machine)
			})
			gymAPI.GET("/machine/category/:gym_id", func(c *gin.Context) {
				v1.GetMachineListByCategory(c, dbInstance.Machine)
			})
		}
		machineAPI := apiV1.Group("/machine")
		{
			machineAPI.PUT("/status", func(c *gin.Context) {
				v1.UpdateMachineStatus(c, dbInstance.Machine)
			})
		}
		userAPI := apiV1.Group("/user")
		{
			userAPI.GET("/info", func(c *gin.Context) {
				v1.GetClientInfo(c, dbInstance.Client)
			})
			userAPI.GET("/reservation/:account", func(c *gin.Context) {
				v1.GetClientReservation(c, dbInstance.Reservation)
			})
			userAPI.GET("/stat/:account", func(c *gin.Context) {
				v1.GetClientStat(c, dbInstance.Client)
			})
			userAPI.GET("/staff/stat", func(c *gin.Context) {
				v1.GetCompanyStat(c, dbInstance.Attendance)
			})
			userAPI.GET("/available", func(c *gin.Context) {
				v1.GetAvailableTime(c, dbInstance.Machine, dbInstance.Reservation)
			})
			userAPI.POST("/login", func(c *gin.Context) {
				v1.Login(c, dbInstance.Client, dbInstance.Staff)
			})
			userAPI.POST("/signup", func(c *gin.Context) {
				v1.Signup(c, dbInstance.Client)
			})
			userAPI.POST("/reservation", func(c *gin.Context) {
				v1.MakeReservation(c, dbInstance.Reservation)
			})
			userAPI.PUT("/info", func(c *gin.Context) {
				v1.UpdateClientInfo(c, dbInstance.Client)
			})

		}
		// apiV1.GET("/test", v1.TTTTT)
	}
	return router
}
