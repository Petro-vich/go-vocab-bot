package main

import (
	"github.com/joho/godotenv"
	"go-vocab-bot/internal/infrastructure/translation"
	"go-vocab-bot/internal/storage"
	"go-vocab-bot/internal/telegram"
	"go-vocab-bot/internal/usecase"
	"log"
	"os"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("error load .env file:", err)
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не найдет в .env")
	}

	translator := translation.NewTranslator()

	rep, err := storage.InitRepository()
	if err != nil {
		log.Fatalf("error init repository: %v", err)
	}

	wordUC := usecase.NewWordUseCase(rep, translator)

	tgbot, err := telegram.NewBot(token, wordUC)
	if err != nil {
		log.Fatal("error init bot:", err)
	}

	if err := tgbot.Start(); err != nil {
		log.Fatal("error start bot:", err)
	}

}
