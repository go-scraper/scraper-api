package main

import (
	"log"

	"scraper/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/scrape", handlers.ScrapeHandler)
	router.GET("/scrape/:id/:page", handlers.PageHandler)

	log.Fatal(router.Run(":8080"))
}
