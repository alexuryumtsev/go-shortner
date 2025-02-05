package fileutils

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alexuryumtsev/go-shortener/internal/app/models"
)

// EnsureDirExists проверяет, существует ли директория, и создаёт её, если не существует.
func EnsureDirExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// FileStorage представляет собой хранилище, использующее файлы для персистентности.
type FileStorage struct {
	filePath string
}

// NewFileStorage создает новое файловое хранилище.
func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{filePath: filePath}
}

// SaveRecord сохраняет запись в файл.
func (fs *FileStorage) SaveRecord(w io.WriteCloser, counter int, urlModel models.URLModel) error {
	defer w.Close()

	bufferedWriter := bufio.NewWriter(w)
	defer bufferedWriter.Flush()

	record := struct {
		UUID        string `json:"uuid"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}{
		UUID:        strconv.Itoa(counter),
		ShortURL:    urlModel.ID,
		OriginalURL: urlModel.URL,
	}

	encoder := json.NewEncoder(bufferedWriter)
	if err := encoder.Encode(record); err != nil {
		return err
	}

	return nil
}

// LoadRecords загружает записи из файла.
func (fs *FileStorage) LoadRecords(r io.Reader) (map[string]string, error) {
	data := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var record struct {
			UUID        string `json:"uuid"`
			ShortURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, err
		}
		data[record.ShortURL] = record.OriginalURL
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
