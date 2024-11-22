package bootstrap

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type Container struct {
	db *pgx.Conn
}

func InitContainer() *Container {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to read config with error: %v", err))
	}

	return &Container{}
}
