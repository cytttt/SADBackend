package routers

import (
	"SADBackend/controllers"
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
	return router
}
