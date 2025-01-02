package services

import (
	"fmt"
	"net/http"
	"scraper/models"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCheckURLStatus(test_type *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	client := &http.Client{
		Transport: httpmock.DefaultTransport,
	}

	// Mock URLs
	urls := []models.URLStatus{
		{URL: "http://example.com/valid"},
		{URL: "http://example.com/invalid"},
		{URL: "http://example.com/error"},
	}

	// Mock response for valid URL
	httpmock.RegisterResponder("GET", "http://example.com/valid", httpmock.NewStringResponder(200, "OK"))

	// Mock response for invalid URL (non-2xx status)
	httpmock.RegisterResponder("GET", "http://example.com/invalid", httpmock.NewStringResponder(404, "Not Found"))

	// Mock response for error (e.g., network failure)
	httpmock.RegisterResponder("GET", "http://example.com/error", httpmock.NewErrorResponder(fmt.Errorf("network error")))

	inaccessibleCount := CheckURLStatus(client, urls, 0, len(urls))

	assert.Equal(test_type, 2, inaccessibleCount, "The count of inaccessible URLs should be 2")
	assert.NotNil(test_type, urls[2].Error, "Expected an error for the network failure URL")
	assert.Equal(test_type, 404, urls[1].HTTPStatus, "Expected 404 status for the invalid URL")
}
