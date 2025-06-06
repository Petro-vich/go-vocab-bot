package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-vocab-bot/internal/usecase"
	"go-vocab-bot/internal/word"
	"strconv"
	"strings"
)

const (
	AddCommand      = "add"
	DeleteCommand   = "delete"
	ListCommand     = "list"
	TrainingCommand = "train"
)

type Bot struct {
	tgBot        *tgbotapi.BotAPI
	wordUC       *usecase.WordUseCase
	trainSession *trainingSession
}

type trainingSession struct {
	words        []word.Word
	isTraining   bool
	currentIndex int
}

func NewBot(token string, useCase *usecase.WordUseCase) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	tr := trainingSession{
		words:        []word.Word{},
		isTraining:   false,
		currentIndex: 0,
	}

	if err != nil {
		return nil, err
	}
	return &Bot{bot,
		useCase,
		&tr}, nil
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	updates := b.tgBot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
	return nil
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := b.tgBot.Send(msg)
	if err != nil {
		panic(err)
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	var parsText []string
	parsText = strings.Split(msg.Text, " ")
	command := msg.Command()
	bodyText := msg.CommandArguments()

	if b.trainSession.isTraining {
		if b.checkWords(parsText[0]) {
			b.sendMessage(msg.Chat.ID, "Верно")
		} else {
			b.sendMessage(msg.Chat.ID, "Не верно")
		}
		b.sendNextTrainingWords(msg)
	}

	switch command {
	case AddCommand:
		if 1 >= len(parsText) {
			b.sendMessage(msg.Chat.ID, "Используйте команду /add <english>")
			return
		}
		if err := b.wordUC.Add(strings.ToLower(bodyText)); err != nil {
			b.sendMessage(msg.Chat.ID, "Ошибка добавления слова")
		}
		b.sendMessage(msg.Chat.ID, "Слово успешно добавлено")
	case TrainingCommand:
		trainsList, err := b.wordUC.TrainList()
		if err != nil {
			b.sendMessage(msg.Chat.ID, "Ошибка получения тренировочного листа")
			return
		}
		b.startTraining(trainsList, msg)

	case DeleteCommand:
		idToDelete, err := strconv.Atoi(bodyText)
		if err != nil {
			b.sendMessage(msg.Chat.ID, "Ошибка удаления. Испольуейте комманду: /del <id>")
			return
		}

		if err := b.wordUC.Delete(idToDelete); err != nil {
			b.sendMessage(msg.Chat.ID, "Пока не знаю как обработать")
			return
		}
		b.sendMessage(msg.Chat.ID, "Успешно удалено")
	}
}

func (b *Bot) startTraining(word []word.Word, msg *tgbotapi.Message) {
	b.trainSession.words = word
	b.trainSession.isTraining = true
	b.trainSession.currentIndex = 0
	b.sendNextTrainingWords(msg)
}

func (b *Bot) sendNextTrainingWords(msg *tgbotapi.Message) {
	cIndex := b.trainSession.currentIndex
	if cIndex >= len(b.trainSession.words) {
		b.sendMessage(msg.Chat.ID, "Тренировка завершена")
		b.trainSession.isTraining = false
		return
	}
	wdForTr := b.trainSession.words[cIndex].English
	b.sendMessage(msg.Chat.ID, fmt.Sprintf("Переведите слово: %s", wdForTr))
	b.trainSession.currentIndex++
}

func (b *Bot) checkWords(bodyText string) bool {
	if bodyText == b.trainSession.words[b.trainSession.currentIndex-1].Russian {
		return true
	} else {
		return false
	}
}
