package storage

import (
	"encoding/json"
	"go-vocab-bot/pkg/model"
	"os"
	"path/filepath"
	"sync"
)

type JSONRepo struct {
	mu       sync.RWMutex
	filepath string
	data     map[int64][]model.Word
}

func NewJSONrep(path string) (*JSONRepo, error) {
	repo := &JSONRepo{
		filepath: path,
		data:     make(map[int64][]model.Word),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *JSONRepo) load() error {
	file, err := os.ReadFile(r.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(file, &r.data)
}

func (r *JSONRepo) GetWords(chatID int64) ([]model.Word, error) {
	r.mu.RLock()
	defer r.mu.Unlock()

	return r.data[chatID], nil
}

func (r *JSONRepo) SaveWords(chatID int64, words []model.Word) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[chatID] = words
	return r.save()
}

func (r *JSONRepo) save() error {
	// Создаём директорию, если её нет
	dir := filepath.Dir(r.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(r.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filepath, bytes, 0644)
}
