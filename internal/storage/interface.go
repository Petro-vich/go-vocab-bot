package storage

import (
	"go-vocab-bot/internal/word"
)

type Repository interface {
	AddWord(word word.Word) error
	DeleteWord(id int) error
	GetAll() (wd []word.Word, err error)
}
