package model

import "time"

type Word struct {
	English     string    // английское слово
	Translation string    // перевод
	Level       int       // уровень повторения, начнём с 0
	NextReview  time.Time // время следующего повторения
}

type ReviewSession struct {
	Index int
}
