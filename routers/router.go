package routers

import (
	"SADBackend/controllers"
	v1 "SADBackend/controllers/v1"
	_ "SADBackend/docs"
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

func InitRouters() *gin.Engine {
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
			gymAPI.GET("/list", v1.GetGymList)
			gymAPI.GET("/machine", v1.GetMachineList)
		}
		machineAPI := apiV1.Group("/machine")
		{
			machineAPI.PUT("/status", v1.UpdateMachineStatus)
		}
	}
	return router
}
