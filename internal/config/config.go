package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init(env string, service string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath(fmt.Sprintf("../docker/%s", service))
	config.AddConfigPath(fmt.Sprintf("internal/config/%s", service))
	config.AddConfigPath(fmt.Sprintf("config/%s", service))
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
	log.Printf("config file used : %v \n", config.ConfigFileUsed())

}

func InitTest() {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("test")
	config.AddConfigPath("../../docker/sms_server/aibot")
	config.AddConfigPath("../../../docker/sms_server/aibot")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

func GetConfig() *viper.Viper {
	return config
}

func IsProduction() bool {
	env := config.GetString("env")

	return env == "prod"
}

func IsNotProduction() bool {

	return !IsProduction()
}
func IsDevelopment() bool {
	env := config.GetString("env")

	return env == "dev"
}
func IsLocal() bool {
	env := config.GetString("env")

	return env == "local"
}
func IsNotLocal() bool {

	return !IsLocal()
}
