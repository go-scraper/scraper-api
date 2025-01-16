package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
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
		mockSessionID  string
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
			mockSessionID:  "mockSessionID",
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
			name: "Unexpected Error Fetching Page Info",
			queryParams: map[string]string{
				"url": "http://example.com",
			},
			mockPageInfo:   nil,
			mockError:      assert.AnError,
			mockRequestID:  "",
			mockSessionID:  "",
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "An unexpected error occurred",
			},
		},
		{
			name: "Timeout Error Fetching Page Info",
			queryParams: map[string]string{
				"url": "http://example.com",
			},
			mockPageInfo: nil,
			mockError: &net.DNSError{
				IsTimeout: true,
			},
			mockRequestID:  "",
			mockSessionID:  "",
			expectedStatus: http.StatusGatewayTimeout,
			expectedBody: map[string]interface{}{
				"error": "Request timeout during the page fetch",
			},
		},
		{
			name: "Failed to Reach The Request URL",
			queryParams: map[string]string{
				"url": "http://example.com",
			},
			mockPageInfo: nil,
			mockError: &net.DNSError{
				IsTimeout: false,
			},
			mockRequestID:  "",
			mockSessionID:  "",
			expectedStatus: http.StatusBadGateway,
			expectedBody: map[string]interface{}{
				"error": "Failed to reach the requested URL",
			},
		},
		{
			name: "Failed to Reach The Request URL",
			queryParams: map[string]string{
				"url": "http://example",
			},
			mockPageInfo:   nil,
			mockError:      nil,
			mockRequestID:  "",
			mockSessionID:  "",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid URL format, please provide a valid URL.",
			},
		},
	}

	for _, test_data := range tests {
		test_type.Run(test_data.name, func(test_type *testing.T) {

			patchFetchPageInfo := monkey.Patch(services.FetchPageInfo,
				func(client *http.Client, url string) (*models.PageInfo, error) {
					return test_data.mockPageInfo, test_data.mockError
				})
			defer patchFetchPageInfo.Unpatch()

			patchRetriveDatabase := monkey.Patch(storage.RetriveDatabase,
				func(sessionID string) storage.Database {
					return &storage.InMemoryDatabase{Data: make(map[string]models.PageInfo, 10_000)}
				})
			defer patchRetriveDatabase.Unpatch()

			patchStorePageInfo := monkey.PatchInstanceMethod(
				reflect.TypeOf(&storage.InMemoryDatabase{}), // Type of the struct
				"StorePageInfo", // Method name to patch
				func(db *storage.InMemoryDatabase, info *models.PageInfo) string {
					// Mocked implementation
					db.Data[test_data.mockRequestID] = *info
					return test_data.mockRequestID
				},
			)
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

			patchRetriveDatabase := monkey.Patch(storage.RetriveDatabase,
				func(sessionID string) storage.Database {
					return &storage.InMemoryDatabase{Data: make(map[string]models.PageInfo, 10_000)}
				})
			defer patchRetriveDatabase.Unpatch()

			patchRetrievePageInfo := monkey.PatchInstanceMethod(
				reflect.TypeOf(&storage.InMemoryDatabase{}), // Type of the struct
				"RetrievePageInfo",                          // Method name to patch
				func(db *storage.InMemoryDatabase, id string) (*models.PageInfo, bool) {
					return test_data.mockPageInfo, test_data.mockExists
				},
			)
			defer patchRetrievePageInfo.Unpatch()

			patchCalculatePageBounds := monkey.Patch(utils.CalculatePageBounds,
				func(pageNum, totalItems, pageSize int) (int, int) {
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
