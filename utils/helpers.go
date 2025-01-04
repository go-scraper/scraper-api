package utils

import (
	"fmt"
	"math"
	"scraper/config"
	"scraper/models"

	"github.com/gin-gonic/gin"
)

func CalculateTotalPages(totalItems, pageSize int) int {
	return int(math.Ceil(float64(totalItems) / float64(pageSize)))
}

func CalculatePageBounds(pageNum, totalItems, pageSize int) (int, int) {
	start := (pageNum - 1) * pageSize
	end := min(start+pageSize, totalItems)
	return start, end
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// This is to build the response after a successful scraping.
func BuildPageResponse(requestID string, pageNum, totalPages int, pageInfo *models.PageInfo, inaccessible, start, end int) models.PageResponse {
	var prevPage, nextPage *string
	if pageNum > 1 {
		prev := fmt.Sprintf("/scrape/%s/%d", requestID, pageNum-1)
		prevPage = &prev
	}
	if end < len(pageInfo.URLs) {
		next := fmt.Sprintf("/scrape/%s/%d", requestID, pageNum+1)
		nextPage = &next
	}

	return models.PageResponse{
		RequestID: requestID,
		Pagination: models.Pagination{
			PageSize:    config.GetURLCheckPageSize(),
			CurrentPage: pageNum,
			TotalPages:  totalPages,
			PrevPage:    prevPage,
			NextPage:    nextPage,
		},
		Scraped: models.ScrapedData{
			HTMLVersion:       pageInfo.HTMLVersion,
			Title:             pageInfo.Title,
			Headings:          pageInfo.HeadingCounts,
			ContainsLoginForm: pageInfo.ContainsLoginForm,
			TotalURLs:         len(pageInfo.URLs),
			InternalURLs:      pageInfo.InternalURLsCount,
			ExternalURLs:      pageInfo.ExternalURLsCount,
			Paginated: models.PaginatedURLs{
				InaccessibleURLs: inaccessible,
				URLs:             pageInfo.URLs[start:end],
			},
		},
	}
}

// This is to build the error response.
func BuildErrorResponse(message string) gin.H {
	return gin.H{"error": message}
}
