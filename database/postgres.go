package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func NewConnection() (*sql.DB, error) {
	viper.SetDefault("database.port", "5432")
	password := viper.GetString("database.password")
	if password != "" {
		password = fmt.Sprintf(":%s", password)
	}
	connStr := fmt.Sprintf(
		"postgres://%s%s@%s:%d/%s?sslmode=disable",
		viper.GetString("database.username"),
		password,
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.database"),
	)
	log.Println(connStr)

	return sql.Open("postgres", connStr)
}
