// This is a simple in-memory storage to store page info to support for pagination.
// Each store page info is mapped to a random unique ID which generated upon storing data.
// To retrieve stored page info need to provide the ID generator upon storing data.
// This simple storage supports only store and retrieve operations as of now.
package storage

import (
	"fmt"
	"math/rand"
	"scraper/logger"
	"scraper/models"
	"sync"
	"time"
)

// Database registry to store dbs per user session
// Here we set the initial capacity to 10000 to reduce resizing overhead.
var dbRegistry = struct {
	sync.RWMutex
	dbs map[string]Database
}{dbs: make(map[string]Database, 10_000)}

type Database interface {
	StorePageInfo(info *models.PageInfo) string
	RetrievePageInfo(id string) (*models.PageInfo, bool)
}

type InMemoryDatabase struct {
	sync.RWMutex
	Data map[string]models.PageInfo
}

// This is to retrieve the database by session ID.
// If the database does not exist, it will create a new one for the current session.
// Here we set the initial capacity to 10000 to reduce resizing overhead.
func RetriveDatabase(sessionID string) Database {
	dbRegistry.RLock()
	defer dbRegistry.RUnlock()

	db, exists := dbRegistry.dbs[sessionID]
	if !exists {
		db = &InMemoryDatabase{Data: make(map[string]models.PageInfo, 10_000)}
		dbRegistry.dbs[sessionID] = db
	}
	return db
}

// This is to store page info in the given database.
func (db *InMemoryDatabase) StorePageInfo(info *models.PageInfo) string {
	db.Lock()
	defer db.Unlock()

	id := GenerateID()
	db.Data[id] = *info
	logger.Debug(fmt.Sprintf("Updated database with page info:\n %v", db.Data))
	return id
}

// This is to retrieve page info by unique ID from the given database.
func (db *InMemoryDatabase) RetrievePageInfo(id string) (*models.PageInfo, bool) {
	db.RLock()
	defer db.RUnlock()

	info, exists := db.Data[id]
	logger.Debug(fmt.Sprintf("Retrieved page info from:\n %v", db.Data))
	return &info, exists
}

// var storage = struct {
// 	sync.RWMutex
// 	data map[string]models.PageInfo
// }{data: make(map[string]models.PageInfo)}

// // This is to store page info.
// func StorePageInfo(info *models.PageInfo) string {
// 	storage.Lock()
// 	defer storage.Unlock()

// 	id := generateID()
// 	storage.data[id] = *info
// 	return id
// }

// // This is to retrieve page info by unique ID.
// func RetrievePageInfo(id string) (*models.PageInfo, bool) {
// 	storage.RLock()
// 	defer storage.RUnlock()

// 	info, exists := storage.data[id]
// 	return &info, exists
// }

// This is to generate the random unique ID.
func GenerateID() string {
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
