package main

import (
	"fmt"
	"log"

	"scraper/config"
	"scraper/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/scrape", handlers.ScrapeHandler)
	router.GET("/scrape/:session_id/:id/:page", handlers.PageHandler)

	log.Fatal(router.Run(fmt.Sprintf(":%s", config.GetAppPort())))
}
