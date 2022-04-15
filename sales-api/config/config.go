package config

import (
	"sales-api/constants"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"MYSQL_HOST"`
	DBPort     string `mapstructure:"MYSQL_PORT"`
	DBUser     string `mapstructure:"MYSQL_USER"`
	DBPassword string `mapstructure:"MYSQL_PASSWORD"`
	DBName     string `mapstructure:"MYSQL_DBNAME"`
	BindAddr   string `mapstructure:"BIND_ADDR"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.SetDefault(constants.JWTSecret, "01234567890123456789012345678901")
	viper.SetDefault(constants.DBDriver, "mysql")
	viper.SetDefault(constants.BindAddr, ":3030")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
