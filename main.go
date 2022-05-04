package main

import (
	"github.com/spf13/viper"
	"log"
	"upload2/global"
	"upload2/router"
)

func init() {
	err := Setting()
	if err != nil {
		log.Println("fuck it", err)
	}
}

func main() {
	router.NewRouter()
}

func Setting() error {
	viper.SetConfigFile("configs/application.yml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("App", &global.AppSetting); err != nil {
		return err
	}
	return nil

}
