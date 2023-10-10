package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/b-turchyn/lifetracker-backend/database"
	"github.com/b-turchyn/lifetracker-backend/endpoint"
	"github.com/b-turchyn/lifetracker-backend/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := initRouter(db)

	router.Run(":8080")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func initRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	bucketService := service.NewBucketService(db)
	endpoint.BucketEndpoints(router, bucketService)

	return router
}
