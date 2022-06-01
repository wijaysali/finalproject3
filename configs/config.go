package configs

import (
	"MyGram/app"
	"MyGram/configs/rsaconfig"
	"fmt"

	"github.com/spf13/viper"
)

func LoadEnvConfig(path string) {
	//viper.SetConfigName("t.env")
	//viper.Debug()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func LoadConfig() {

	LoadEnvConfig(".")
	LoadPort()
	ConfigDatabase()
	rsaconfig.ConfigRSA()
	ConfigJwtware()
	ConfigValidator()
}
func LoadPort() {
	app.Port.WriteString(":")
	app.Port.WriteString(viper.GetString("APP_PORT"))
}
