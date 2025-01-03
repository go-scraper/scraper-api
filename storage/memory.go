// This is a simple in-memory storage to store page info to support for pagination.
// Each store page info is mapped to a random unique ID which generated upon storing data.
// To retrieve stored page info need to provide the ID generator upon storing data.
// This simple storage supports only store and retrieve operations as of now.
package storage

import (
	"math/rand"
	"scraper/models"
	"sync"
	"time"
)

var storage = struct {
	sync.RWMutex
	data map[string]models.PageInfo
}{data: make(map[string]models.PageInfo)}

// This is to store page info.
func StorePageInfo(info *models.PageInfo) string {
	storage.Lock()
	defer storage.Unlock()

	id := generateID()
	storage.data[id] = *info
	return id
}

// This is to retrieve page info by unique ID.
func RetrievePageInfo(id string) (*models.PageInfo, bool) {
	storage.RLock()
	defer storage.RUnlock()

	info, exists := storage.data[id]
	return &info, exists
}

// This is to generate the random unique ID.
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// This is to generate a random string to append into the random unique ID.
func randomString(size int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, size)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range result {
		result[i] = letters[random.Intn(len(letters))]
	}
	return string(result)
}
