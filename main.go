package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var respectExpressions = []string{"+", "+1"}

func main() {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalln(err)
	}

	bot, err := tgbotapi.NewBotAPI("<TOKEN>")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.ReplyToMessage != nil {
			switch {
			case update.Message.ReplyToMessage.From.IsBot:
				continue
			case isSameUser(update.Message.From.ID, update.Message.ReplyToMessage.From.ID):
				continue
			case containsRespect(update.Message.Text, respectExpressions):
				log.Printf("adding respect to @%v", update.Message.ReplyToMessage.From.UserName)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("+1 respect to @%v", update.Message.ReplyToMessage.From.UserName))
				msg.ReplyToMessageID = update.Message.ReplyToMessage.MessageID

				u := getUser(db, int(update.Message.ReplyToMessage.From.ID))
				u.incrementRespect()
				bot.Send(msg)
			}
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "myrespect":
				u := getUser(db, int(update.Message.From.ID))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("@%v, your respect count is: %v", update.Message.From.UserName, u.respectCount))
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}
}

func containsRespect(message string, expressions []string) bool {
	for _, expression := range expressions {
		if expression == message {
			return true
		}
	}
	return false
}
func isSameUser(senderId, respectedUserId int64) bool {
	if senderId == respectedUserId {
		return true
	}
	return false
}
