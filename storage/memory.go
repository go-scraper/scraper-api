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

func StorePageInfo(info *models.PageInfo) string {
	storage.Lock()
	defer storage.Unlock()

	id := generateID()
	storage.data[id] = *info
	return id
}

func RetrievePageInfo(id string) (*models.PageInfo, bool) {
	storage.RLock()
	defer storage.RUnlock()

	info, exists := storage.data[id]
	return &info, exists
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(size int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, size)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range result {
		result[i] = letters[random.Intn(len(letters))]
	}
	return string(result)
}
