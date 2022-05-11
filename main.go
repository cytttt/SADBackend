package main

import (
	"SADBackend/constant"
	"SADBackend/pkg/mongodb"
	"SADBackend/routers"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	constant.ReadConfig(".env")
	mongodb.Init()
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := routers.InitRouters()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGILL, syscall.SIGFPE)

	// graceful terminate
	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			mongodb.Dispose()
			os.Exit(0)
		}
	}()

	router.Run(viper.GetString("PORT"))
}
