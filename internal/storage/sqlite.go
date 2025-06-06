package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go-vocab-bot/internal/word"
	"time"
)

type WordStorage struct {
	db *sql.DB
}

func InitRepository() (Repository, error) {
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return nil, fmt.Errorf("error open db: %w", err)
	}
	expression := `
		CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    russian TEXT NOT NULL,
    english TEXT NOT NULL,
		createdAt TEXT NOT NULL,
		status TEXT NOT NULL,
		lvl INTEGER NOT NULL                
)`
	_, err = db.Exec(expression)
	if err != nil {
		return nil, fmt.Errorf("error create table: %w", err)
	}

	return &WordStorage{db}, nil
}

func (s *WordStorage) AddWord(word word.Word) error {
	expression := `
	INSERT INTO words (
	russian,
	english,
	createdAt,
  status,
	lvl)
	VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(expression, word.Russian, word.English, word.CreatedAt.Format("2006-01-02 15:04:05"), word.Status, word.Lvl)
	if err != nil {
		return fmt.Errorf("error insert: %w", err)
	}

	return nil
}

func (s *WordStorage) DeleteWord(id int) error {
	expression := `DELETE FROM words WHERE id = ?`
	res, err := s.db.Exec(expression, id)
	if err != nil {
		return fmt.Errorf("error delete: %w", err)
	}

	if num, _ := res.RowsAffected(); num == 0 {
		return fmt.Errorf("word not found")
	}
	return nil
}

func (s *WordStorage) GetAll() ([]word.Word, error) {
	var result []word.Word
	var createdAtStr string
	rows, err := s.db.Query(`SELECT
    russian,
    english,
    createdAt, 
    status, 
    lvl 
	FROM words`)

	if err != nil {
		return result, fmt.Errorf("error query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var wd word.Word
		err := rows.Scan(
			&wd.Russian,
			&wd.English,
			&createdAtStr,
			&wd.Status,
			&wd.Lvl)

		if err != nil {
			return result, fmt.Errorf("error scan: %w", err)
		}
		wd.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return result, fmt.Errorf("error parse createdAt: %w", err)
		}

		result = append(result, wd)
	}
	return result, nil
}
