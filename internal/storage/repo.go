package storage

import "go-vocab-bot/pkg/model"

type WordRepository interface {
	GetWords(chatID int64) ([]model.Word, error)
	SaveWords(chatID int64, words []model.Word) error
}
