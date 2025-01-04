package config

import (
	"fmt"
	"os"
	"scraper/logger"
	"strconv"
)

// Default values for environment variables
const (
	defaultAppPort                           = "8080"
	defaultURLCheckPageSize                  = 10
	defaultOutgoingScrapeRequestTimeout      = 30
	defaultOutgoingAccessibilityCheckTimeout = 10
)

// Configuration variables initialized once
var (
	appPort                           string
	urlCheckPageSize                  int
	outgoingScrapeRequestTimeout      int
	outgoingAccessibilityCheckTimeout int
)

func init() {
	// Load environment variables and set defaults if necessary
	appPort = getEnv("APP_PORT", defaultAppPort)

	urlCheckPageSize = parseEnvAsInt("URL_STATUS_CHECK_PAGE_SIZE", defaultURLCheckPageSize)
	outgoingScrapeRequestTimeout = parseEnvAsInt("OUT_GOING_SCRAPE_REQ_TIMEOUT",
		defaultOutgoingScrapeRequestTimeout)
	outgoingAccessibilityCheckTimeout = parseEnvAsInt("OUT_GOING_URL_ACCESSIBILITY_CHECK_TIMEOUT",
		defaultOutgoingAccessibilityCheckTimeout)
}

// Helper function to get environment variable or return a default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to parse environment variable as int or return a default
func parseEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid value for %s: %v", key, err))
		return defaultValue
	}
	return parsedValue
}

// Exported getter functions
func GetAppPort() string {
	return appPort
}

func GetURLCheckPageSize() int {
	return urlCheckPageSize
}

func GetOutgoingScrapeRequestTimeout() int {
	return outgoingScrapeRequestTimeout
}

func GetOutgoingAccessibilityCheckTimeout() int {
	return outgoingAccessibilityCheckTimeout
}
