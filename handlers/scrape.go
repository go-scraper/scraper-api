package handlers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"

	"scraper/config"
	"scraper/logger"
	"scraper/services"
	"scraper/storage"
	"scraper/utils"

	"github.com/gin-gonic/gin"
)

func ScrapeHandler(context *gin.Context) {
	baseURL := context.Query("url")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable TLS verification
		},
	}

	if baseURL == "" {
		logger.Debug("URL query parameter is required")
		context.JSON(http.StatusBadRequest, gin.H{"error": "url query parameter is required"})
		return
	}

	pageInfo, err := services.FetchPageInfo(client, baseURL)
	if err != nil {
		logger.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch page info"})
		return
	}

	requestID := storage.StorePageInfo(pageInfo)
	inaccessibleCount := services.CheckURLStatus(client, pageInfo.URLs, 0, min(config.PageSize, len(pageInfo.URLs)))
	totalPages := utils.CalculateTotalPages(len(pageInfo.URLs), config.PageSize)

	context.JSON(http.StatusOK, utils.BuildPageResponse(requestID, 1, totalPages, pageInfo, inaccessibleCount, 0, min(config.PageSize, len(pageInfo.URLs))))
}

func PageHandler(context *gin.Context) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable TLS verification
		},
	}
	requestID := context.Param("id")
	pageNumStr := context.Param("page")

	pageInfo, exists := storage.RetrievePageInfo(requestID)
	if !exists {
		logger.Debug(fmt.Sprintf("Requested ID [%s] not found in the local storage", requestID))
		context.JSON(http.StatusNotFound, gin.H{"error": "request ID not found"})
		return
	}

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil || pageNum < 1 {
		logger.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	start, end := utils.CalculatePageBounds(pageNum, len(pageInfo.URLs), config.PageSize)
	if start >= len(pageInfo.URLs) {
		logger.Debug(fmt.Sprintf("Requested page [%d] not found", pageNum))
		context.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}

	inaccessibleCount := services.CheckURLStatus(client, pageInfo.URLs, start, end)
	totalPages := utils.CalculateTotalPages(len(pageInfo.URLs), config.PageSize)

	context.JSON(http.StatusOK, utils.BuildPageResponse(requestID, pageNum, totalPages, pageInfo, inaccessibleCount, start, end))
}
