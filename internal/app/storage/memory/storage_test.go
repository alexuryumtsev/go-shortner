package memory

import (
	"context"
	"testing"

	"github.com/alexuryumtsev/go-shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStorage_SaveAndGet(t *testing.T) {
	storage := NewInMemoryStorage()
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

func TestInMemoryStorage_GetNonExistent(t *testing.T) {
	storage := NewInMemoryStorage()
	ctx := context.Background()

	// Test Get for non-existent URL
	_, exists := storage.Get(ctx, "nonExistentID")
	assert.False(t, exists, "URL should not exist in storage")
}

func TestInMemoryStorage_LoadFromFile(t *testing.T) {
	storage := NewInMemoryStorage()

	// Test LoadFromFile (should do nothing and return nil)
	err := storage.LoadFromFile()
	assert.NoError(t, err, "LoadFromFile should not return an error")
}
