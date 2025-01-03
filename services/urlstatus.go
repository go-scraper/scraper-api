package services

import (
	"net/http"
	"scraper/logger"
	"scraper/models"
	"sync"
)

// This is to check the URL status and decide wether it is accessible or not.
// It marks the status of each collected URL.
// Since the URL collection can be huge we check status based on given start and end positions.
func CheckURLStatus(client *http.Client, urls []models.URLStatus, start, end int) int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var inaccessibleCount int

	for i := start; i < end; i++ {
		wg.Add(1)

		go func(idx int) {
			defer wg.Done()

			resp, err := client.Get(urls[idx].URL)
			if err != nil {
				logger.Error(err)

				mu.Lock()
				inaccessibleCount++
				mu.Unlock()

				urls[idx].Error = err.Error()
				return
			}

			defer resp.Body.Close()
			urls[idx].HTTPStatus = resp.StatusCode

			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				mu.Lock()
				inaccessibleCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return inaccessibleCount
}
