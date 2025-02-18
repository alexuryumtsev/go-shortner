package pg

import (
	"context"
	"testing"

	"github.com/alexuryumtsev/go-shortener/internal/app/db"
	"github.com/alexuryumtsev/go-shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *db.Database {
	ctx := context.Background()
	dsn := "postgres://postgres:test@localhost:5432/testdb?sslmode=disable"

	// Подключаемся к тестовой базе данных
	database, err := db.NewDatabaseConnection(ctx, dsn)
	require.NoError(t, err, "Failed to connect to test database")

	// Очищаем таблицу перед каждым тестом
	_, err = database.Pool.Exec(ctx, "TRUNCATE TABLE urls")
	require.NoError(t, err, "Failed to truncate table")

	return database
}

func TestDatabaseStorage_SaveAndGet(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	storage := NewDatabaseStorage(database)

	ctx := context.Background()

	urlModel := models.URLModel{
		ID:  "testID",
		URL: "https://example.com",
	}

	// Test Save
	err := storage.Save(ctx, urlModel)
	assert.NoError(t, err, "Save should not return an error")

	// Test Get
	retrievedURLModel, exists := storage.Get(ctx, "testID")
	assert.True(t, exists, "URL should exist in storage")
	assert.Equal(t, urlModel, retrievedURLModel, "Retrieved URL should match the saved URL")
}

func TestDatabaseStorage_SaveBatch(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	storage := NewDatabaseStorage(database)

	ctx := context.Background()

	urlModels := []models.URLModel{
		{ID: "testID1", URL: "https://example1.com"},
		{ID: "testID2", URL: "https://example2.com"},
	}

	// Test SaveBatch
	err := storage.SaveBatch(ctx, urlModels)
	assert.NoError(t, err, "SaveBatch should not return an error")

	// Test Get for each URL
	for _, urlModel := range urlModels {
		retrievedURLModel, exists := storage.Get(ctx, urlModel.ID)
		assert.True(t, exists, "URL should exist in storage")
		assert.Equal(t, urlModel, retrievedURLModel, "Retrieved URL should match the saved URL")
	}
}

func TestDatabaseStorage_GetNonExistent(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	storage := NewDatabaseStorage(database)

	ctx := context.Background()

	// Test Get for non-existent URL
	_, exists := storage.Get(ctx, "nonExistentID")
	assert.False(t, exists, "URL should not exist in storage")
}

func TestDatabaseStorage_Ping(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	storage := NewDatabaseStorage(database)

	ctx := context.Background()

	// Test Ping
	err := storage.Ping(ctx)
	assert.NoError(t, err, "Ping should not return an error")
}
