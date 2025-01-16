package storage

import (
	"scraper/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorePageInfo(test_type *testing.T) {
	pageInfo := &models.PageInfo{
		Title: "Test Page",
	}

	database := RetriveDatabase("test-session")
	id := database.StorePageInfo(pageInfo)

	// Ensure the ID is not empty
	assert.NotEmpty(test_type, id, "Generated ID should not be empty")

	retrievedInfo, exists := database.RetrievePageInfo(id)

	// Assert that the PageInfo exists
	assert.True(test_type, exists, "Stored PageInfo should be retrievable")
	// Assert that the retrieved info matches the original
	assert.Equal(test_type, pageInfo.Title, retrievedInfo.Title, "Title should match")
}

func TestRetrievePageInfo_NotFound(test_type *testing.T) {
	database := RetriveDatabase("test-session")
	// Try retrieving a non-existent PageInfo
	retrievedInfo, exists := database.RetrievePageInfo("nonexistent-id")

	// Assert that the info does not exist
	assert.False(test_type, exists, "Non-existent ID should not be found")
	// Assert that the returned PageInfo is empty
	assert.Equal(test_type, &models.PageInfo{}, retrievedInfo,
		"Retrieved info should be an empty PageInfo for non-existent ID")
}

func TestGenerateID(test_type *testing.T) {
	// Call the private function indirectly by calling StorePageInfo
	pageInfo := &models.PageInfo{Title: "Test Page"}
	database := RetriveDatabase("test-session")
	id := database.StorePageInfo(pageInfo)

	// Assert that the ID follows the expected format
	assert.Regexp(test_type, `^\d{14}-[a-zA-Z0-9]{8}$`, id,
		"Generated ID should follow the correct format")
}

func TestRandomString(test_type *testing.T) {
	randomStr := randomString(10)

	// Assert that the random string is of the correct length
	assert.Equal(test_type, 10, len(randomStr), "Random string should have the correct length")

	// Check that the string only contains valid characters
	for _, char := range randomStr {
		assert.Contains(test_type,
			"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", string(char),
			"Random string should only contain valid characters")
	}
}
