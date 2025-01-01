package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"scraper/models"
	"scraper/services"
	"scraper/storage"
	"scraper/utils"
	"testing"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestScrapeHandler(test_type *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		mockPageInfo   *models.PageInfo
		mockError      error
		mockRequestID  string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Valid URL",
			queryParams: map[string]string{
				"url": "http://example.com",
			},
			mockPageInfo: &models.PageInfo{
				HTMLVersion:       "HTML 5",
				Title:             "Example Title",
				HeadingCounts:     map[string]int{"h1": 2},
				URLs:              []models.URLStatus{},
				InternalURLsCount: 1,
				ExternalURLsCount: 1,
				ContainsLoginForm: false,
			},
			mockError:      nil,
			mockRequestID:  "mockRequestID",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing URL",
			queryParams:    map[string]string{"url": ""},
			mockPageInfo:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "url query parameter is required",
			},
		},
		{
			name: "Error Fetching Page Info",
			queryParams: map[string]string{
				"url": "http://example.com",
			},
			mockPageInfo:   nil,
			mockError:      assert.AnError,
			mockRequestID:  "",
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "failed to fetch page info",
			},
		},
	}

	for _, test_data := range tests {
		test_type.Run(test_data.name, func(test_type *testing.T) {

			patchFetchPageInfo := monkey.Patch(services.FetchPageInfo, func(url string) (*models.PageInfo, error) {
				return test_data.mockPageInfo, test_data.mockError
			})
			defer patchFetchPageInfo.Unpatch()

			patchStorePageInfo := monkey.Patch(storage.StorePageInfo, func(info *models.PageInfo) string {
				return test_data.mockRequestID
			})
			defer patchStorePageInfo.Unpatch()

			router := gin.Default()
			router.GET("/scrape", ScrapeHandler)

			req := httptest.NewRequest(http.MethodGet, "/scrape", nil)
			query := req.URL.Query()
			for k, v := range test_data.queryParams {
				query.Add(k, v)
			}
			req.URL.RawQuery = query.Encode()

			// Perform the request
			resp_recorder := httptest.NewRecorder()
			router.ServeHTTP(resp_recorder, req)

			assert.Equal(test_type, test_data.expectedStatus, resp_recorder.Code)

			if test_data.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(resp_recorder.Body.Bytes(), &response)
				assert.NoError(test_type, err)

				for k, v := range test_data.expectedBody {
					assert.Equal(test_type, v, response[k])
				}
			}
		})
	}
}

func TestPageHandler(test_type *testing.T) {
	tests := []struct {
		name           string
		requestID      string
		pageNum        string
		mockPageInfo   *models.PageInfo
		mockExists     bool
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "Valid Page Request",
			requestID: "mockRequestID",
			pageNum:   "1",
			mockPageInfo: &models.PageInfo{
				URLs: []models.URLStatus{
					{URL: "http://example.com", HTTPStatus: 200},
				},
			},
			mockExists:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Page Info Not Found",
			requestID:      "nonExistentID",
			pageNum:        "1",
			mockPageInfo:   nil,
			mockExists:     false,
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "request ID not found",
			},
		},
		{
			name:           "Invalid Page Number",
			requestID:      "mockRequestID",
			pageNum:        "invalid",
			mockPageInfo:   nil,
			mockExists:     true,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid page number",
			},
		},
	}

	for _, test_data := range tests {
		test_type.Run(test_data.name, func(test_type *testing.T) {

			patchRetrievePageInfo := monkey.Patch(storage.RetrievePageInfo, func(id string) (*models.PageInfo, bool) {
				return test_data.mockPageInfo, test_data.mockExists
			})
			defer patchRetrievePageInfo.Unpatch()

			patchCalculatePageBounds := monkey.Patch(utils.CalculatePageBounds, func(pageNum, totalItems, pageSize int) (int, int) {
				return 0, 1
			})
			defer patchCalculatePageBounds.Unpatch()

			router := gin.Default()
			router.GET("/page/:id/:page", PageHandler)

			url := "/page/" + test_data.requestID + "/" + test_data.pageNum
			req := httptest.NewRequest(http.MethodGet, url, nil)

			// Perform the request
			resp_recorder := httptest.NewRecorder()
			router.ServeHTTP(resp_recorder, req)

			assert.Equal(test_type, test_data.expectedStatus, resp_recorder.Code)

			if test_data.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(resp_recorder.Body.Bytes(), &response)
				assert.NoError(test_type, err)

				for k, v := range test_data.expectedBody {
					assert.Equal(test_type, v, response[k])
				}
			}
		})
	}
}
