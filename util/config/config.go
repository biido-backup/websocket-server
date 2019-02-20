package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var log = logrus.New()

func LoadConfig(){
	//viper.SetConfigFile("./config/config.yaml")

	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if (err != nil){
		panic(err)
	}
}