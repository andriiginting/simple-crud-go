package config

import (
	"os"
	"fmt"

	"github.com/spf13/viper"
)

var (
	port int
	username string
	dbName string
)

const (
	POSTGRE_USERNAME_KEY = "PGUSERNAME"
	POSTGRE_PORT_KEY     = "PGPORT"
	POSTGRE_PASSWORD_KEY = "PGPASSWORD"
	POSTGRE_DB_NAME_KEY  = "DBNAME"
)

func ConnectionString() string{
	username,port,_,dbName = GetConfig()
	return fmt.Sprintf("user=%s port=%d dbname=%s sslmode=disable", username, port, dbName)
}

func GetConfig() (string, int, string, string) {
	deployEnv := os.Getenv("DEPLOYENV")
	config := viper.New()
	config.SetConfigName("application")
	config.AddConfigPath("./")
	config.AddConfigPath("../")
	config.AddConfigPath("../../")
	config.SetConfigType("yaml")
	config.ReadInConfig()
	
	return config.GetString(deployEnv + "." + POSTGRE_USERNAME_KEY), config.GetInt(deployEnv + "." + POSTGRE_PORT_KEY), config.GetString(deployEnv + "." + POSTGRE_PASSWORD_KEY), config.GetString(deployEnv + "." + POSTGRE_DB_NAME_KEY)
}

