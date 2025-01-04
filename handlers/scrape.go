package handlers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"scraper/config"
	"scraper/logger"
	"scraper/services"
	"scraper/storage"
	"scraper/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/publicsuffix"
)

// This handles the initial scraping request received from the client.
func ScrapeHandler(context *gin.Context) {
	baseURL := context.Query("url")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable TLS verification
		},
		Timeout: time.Duration(config.GetOutgoingScrapeRequestTimeout()) * time.Second,
	}

	if baseURL == "" {
		logger.Debug("URL query parameter is required")
		context.JSON(http.StatusBadRequest,
			utils.BuildErrorResponse("url query parameter is required"))
		return
	} else {
		if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
			baseURL = "http://" + baseURL
		}

		baseUrlParsed, _ := url.Parse(baseURL)
		_, err := publicsuffix.EffectiveTLDPlusOne(baseUrlParsed.Host)
		if err != nil {
			logger.Error(err)
			context.JSON(http.StatusBadRequest,
				utils.BuildErrorResponse("Invalid URL format, please provide a valid URL."))
			return
		}
	}

	pageInfo, err := services.FetchPageInfo(client, baseURL)
	if err != nil {
		logger.Error(err)
		context.JSON(http.StatusInternalServerError,
			utils.BuildErrorResponse("Failed to fetch page info"))
		return
	}

	// We store scraped page info in-memory to use with pagination later.
	// Stored page infomation mapped to the returned request ID.
	requestID := storage.StorePageInfo(pageInfo)
	// Here we check the status of 10 (config.PageSize) scraped URLs.
	inaccessibleCount := services.CheckURLStatus(client, pageInfo.URLs, 0,
		min(config.GetURLCheckPageSize(), len(pageInfo.URLs)))
	totalPages := utils.CalculateTotalPages(len(pageInfo.URLs), config.GetURLCheckPageSize())

	context.JSON(http.StatusOK, utils.BuildPageResponse(requestID, 1, totalPages, pageInfo,
		inaccessibleCount, 0, min(config.GetURLCheckPageSize(), len(pageInfo.URLs))))
}

// This handles subsequent pagination requests to check status of URLs.
func PageHandler(context *gin.Context) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable TLS verification
		},
		Timeout: time.Duration(config.GetOutgoingAccessibilityCheckTimeout()) * time.Second,
	}
	// Request ID is required to fetch infromation from the in-memory storage.
	requestID := context.Param("id")
	pageNumStr := context.Param("page")

	// Retrieve page information from in-memory storage using the request ID.
	pageInfo, exists := storage.RetrievePageInfo(requestID)
	if !exists {
		logger.Debug(fmt.Sprintf("Requested ID [%s] not found in the local storage", requestID))
		context.JSON(http.StatusNotFound, utils.BuildErrorResponse("request ID not found"))
		return
	}

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil || pageNum < 1 {
		logger.Error(err)
		context.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid page number"))
		return
	}

	start, end := utils.CalculatePageBounds(pageNum, len(pageInfo.URLs), config.GetURLCheckPageSize())
	if start >= len(pageInfo.URLs) {
		logger.Debug(fmt.Sprintf("Requested page [%d] not found", pageNum))
		context.JSON(http.StatusNotFound, utils.BuildErrorResponse("page not found"))
		return
	}

	// Check the URL status for URLs on the given pagination page.
	inaccessibleCount := services.CheckURLStatus(client, pageInfo.URLs, start, end)
	totalPages := utils.CalculateTotalPages(len(pageInfo.URLs), config.GetURLCheckPageSize())

	context.JSON(http.StatusOK, utils.BuildPageResponse(requestID, pageNum, totalPages, pageInfo,
		inaccessibleCount, start, end))
}
