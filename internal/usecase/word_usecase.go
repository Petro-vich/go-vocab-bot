package usecase

import (
	"fmt"
	"go-vocab-bot/internal/domain"
	"go-vocab-bot/internal/storage"
	"go-vocab-bot/internal/word"
	"time"
)

const (
	lvl1 = 1
	lvl2 = 9
	lvl3 = 24
	lvl4 = 48
	lvl5 = 96
	lvl6 = 184
	lvl7 = 736
)

type WordUseCase struct {
	repo       storage.Repository
	translator domain.Translator
}

func NewWordUseCase(repo storage.Repository, translator domain.Translator) *WordUseCase {
	return &WordUseCase{
		repo:       repo,
		translator: translator,
	}
}

func (w *WordUseCase) Add(forTranslate string) error {
	listWords, err := w.repo.GetAll()
	if err != nil {
		fmt.Errorf("error get all words: %w", err)
	}

	for _, wd := range listWords {
		if wd.English == forTranslate {
			return fmt.Errorf("слово: '%s', уже существует", forTranslate)
		}
	}

	translated, err := w.translator.Translate(forTranslate, "ru")
	if err != nil {
		return fmt.Errorf("error translate: %w", err)
	}
	newWord := word.Word{
		Russian:   translated,
		English:   forTranslate,
		CreatedAt: time.Now(),
		Status:    "learned",
		Lvl:       1,
	}
	return w.repo.AddWord(newWord)
}

func (w *WordUseCase) Delete(text int) error {
	return w.repo.DeleteWord(text)
}

func (w *WordUseCase) TrainList() (listWords []word.Word, err error) {
	words, err := w.repo.GetAll()
	if err != nil {
		fmt.Errorf("error get all words: %w", err)
	}

	var lw []word.Word
	for _, wd := range words {
		if needTrain(wd) {
			lw = append(lw, wd)
		}
	}
	return lw, nil
}

func needTrain(wd word.Word) bool {
	now := time.Now()
	housePassed := now.Sub(wd.CreatedAt)

	switch wd.Lvl {
	case lvl1:
		return lvl1 <= housePassed
	case lvl2:
		return lvl2 <= housePassed
	case lvl3:
		return lvl3 <= housePassed
	case lvl4:
		return lvl4 <= housePassed
	case lvl5:
		return lvl5 <= housePassed
	case lvl6:
		return lvl6 <= housePassed
	case lvl7:
		return lvl7 <= housePassed
	default:
		return false
	}
}
