package services

import (
	"bytes"
	"fmt"
	"net/http"
	"scraper/models"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/assert"
)

func TestFetchPageInfo(test_type *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	client := &http.Client{
		Transport: httpmock.DefaultTransport,
	}

	tests := []struct {
		name       string
		mockURL    string
		mockBody   string
		mockStatus int
		mockError  error
		expected   *models.PageInfo
		expectErr  bool
	}{
		{
			name:       "Valid HTML Page",
			mockURL:    "http://valid-url",
			mockBody:   `<html><head><title>Example Domain</title></head><body></body></html>`,
			mockStatus: http.StatusOK,
			expected: &models.PageInfo{
				Title: "Example Domain",
			},
			expectErr: false,
		},
		{
			name:      "HTTP Get Error",
			mockURL:   "http://invalid-url",
			mockError: fmt.Errorf("mocked error"),
			expectErr: true,
		},
	}

	for _, test_data := range tests {
		test_type.Run(test_data.name, func(test_type *testing.T) {

			if test_data.mockError != nil {
				httpmock.RegisterResponder("GET", test_data.mockURL,
					httpmock.NewErrorResponder(test_data.mockError))
			} else {
				httpmock.RegisterResponder("GET", test_data.mockURL,
					httpmock.NewStringResponder(test_data.mockStatus, test_data.mockBody))
			}

			pageInfo, err := FetchPageInfo(client, test_data.mockURL)

			if test_data.expectErr {
				if err == nil {
					test_type.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					test_type.Errorf("Did not expect error, got %v", err)
				}
				if pageInfo.Title != test_data.expected.Title {
					test_type.Errorf("Expected title %s, got %s", test_data.expected.Title, pageInfo.Title)
				}
			}
		})
	}
}

func TestParseHTML(test_type *testing.T) {
	tests := []struct {
		name           string
		htmlContent    string
		baseURL        string
		expectedResult *models.PageInfo
		expectedError  error
	}{
		{
			name: "Valid HTML",
			htmlContent: `
				<!DOCTYPE html>
				<html>
					<head><title>Sample Page</title></head>
					<body>
						<h1>Heading 1</h1>
						<a href="/internal">Internal Link</a>
						<a href="http://external.com">External Link</a>
						<form>
							<input type="password"/>
						</form>
					</body>
				</html>
			`,
			baseURL: "http://example.com",
			expectedResult: &models.PageInfo{
				HTMLVersion:       "HTML 5",
				Title:             "Sample Page",
				HeadingCounts:     map[string]int{"h1": 1},
				InternalURLsCount: 1,
				ExternalURLsCount: 1,
				ContainsLoginForm: true,
				URLs: []models.URLStatus{
					{URL: "http://example.com/internal"},
					{URL: "http://external.com"},
				},
			},
			expectedError: nil,
		},
	}

	for _, test_data := range tests {
		test_type.Run(test_data.name, func(test_type *testing.T) {

			body := bytes.NewReader([]byte(test_data.htmlContent))
			result, err := ParseHTML(body, test_data.baseURL)

			if test_data.expectedError != nil {
				assert.Error(test_type, err)
				assert.Equal(test_type, test_data.expectedError.Error(), err.Error())
			} else {
				assert.NoError(test_type, err)
				assert.Equal(test_type, test_data.expectedResult.Title, result.Title)
				assert.Equal(test_type, test_data.expectedResult.HTMLVersion, result.HTMLVersion)
				assert.Equal(test_type, test_data.expectedResult.HeadingCounts, result.HeadingCounts)
				assert.Equal(test_type, test_data.expectedResult.InternalURLsCount, result.InternalURLsCount)
				assert.Equal(test_type, test_data.expectedResult.ExternalURLsCount, result.ExternalURLsCount)
				assert.Equal(test_type, test_data.expectedResult.ContainsLoginForm, result.ContainsLoginForm)
				assert.ElementsMatch(test_type, test_data.expectedResult.URLs, result.URLs)
			}
		})
	}
}
