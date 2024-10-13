package config

import (
	"github.com/spf13/viper"
	"log"
)

type appConfig struct {
	LoggingEnabled bool   `mapstructure:"logging_enabled"`
	LogFile        string `mapstructure:"log_file"`
	OktaIssuer     string `mapstructure:"OKTA_ISSUER"`
	OktaClientID   string `mapstructure:"OKTA_CLIENT_ID"`
	EnableOktaAuth bool   `mapstructure:"ENABLE_OKTA_AUTH"`
	RateLimit      int    `mapstructure:"RATE_LIMIT"`
	BackendURL     string `mapstructure:"BACKEND_URL"`
	MainMp3Track   string `mapstructure:"MAIN_MP3_TRACK"`
}

type dbConfig struct {
	LoggingEnabled bool   `mapstructure:"logging_enabled"`
	LogFile        string `mapstructure:"log_file"`
	DbType         string `mapstructure:"db_type"`
	DbDsn          string `mapstructure:"db_dsn"`
	OktaIssuer     string `mapstructure:"okta_issuer"`
	OktaClientID   string `mapstructure:"okta_client_id"`
	EnableOktaAuth bool   `mapstructure:"enable_okta_auth"`
	RateLimit      int    `mapstructure:"rate_limit"`
}

var AppConfig appConfig

var DbConfig dbConfig

func LoadConfig(filePath string, db bool) error {
	log.Println("Loading config file...", "db:", db, " and filePath:", filePath)
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("err loading config file:", err)
		return err
	}

	if !db {
		err = viper.Unmarshal(&AppConfig)
		log.Println("AppConfig:", AppConfig)
	} else {
		err = viper.Unmarshal(&DbConfig)
		log.Println("DbConfig:", DbConfig)
	}

	if err != nil {
		log.Println("err unmarshalling config file:", err)
		return err
	}
	return nil
}
