package constant

import (
	"log"

	"github.com/spf13/viper"
)

func ReadConfig(configPath string) {
	viper.SetConfigFile(configPath)
	viper.AddConfigPath(".")

	viper.SetDefault("PORT", ":8888")
	envs := []string{
		"MONGO_DB_CONNECTION",
		"MONGO_DB_NAME",
	}

	for _, env := range envs {
		ErrorHandleLogger(viper.BindEnv(env))
	}

	ErrorHandleLogger(viper.ReadInConfig())
}

func ErrorHandleLogger(err error) {
	if err != nil {
		log.Println(err)
	}
}
