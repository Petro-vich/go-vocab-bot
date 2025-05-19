package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go-vocab-bot/internal/storage"
	"log"
	"os"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не найдет в .env")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	log.Print("Запущен бот:", bot.Self.UserName)

	rep, err := storage.NewJSONrep("internal/storage/words.json")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(rep)

}
