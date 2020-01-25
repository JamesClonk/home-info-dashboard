package telebot

import (
	"log"
	"sync"
	"time"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	telebot *Telebot = nil
	once    sync.Once
)

type Telebot struct {
	hdb    database.HomeInfoDB
	api    *tgbotapi.BotAPI
	ChatID int64
}

func Get(hdb database.HomeInfoDB) *Telebot {
	once.Do(func() {
		token := env.MustGet("TELEBOT_TOKEN")
		api, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Fatalf("Telebot API error: %v\n", err)
		}

		api.Debug = true
		log.Printf("Telebot authorized on account [%s]", api.Self.UserName)
		telebot = &Telebot{hdb, api, 0}

		telebot.WatchForChatID() // background checking for chatID
	})
	return telebot
}

func (bot *Telebot) WatchForChatID() {
	username := env.MustGet("TELEBOT_USERNAME")

	go func() {
		for {
			if bot.ChatID <= 0 {
				u := tgbotapi.NewUpdate(0)
				u.Timeout = 60
				updates, err := bot.api.GetUpdates(u)
				if err != nil {
					log.Fatalf("Telebot API error: %v\n", err)
				}

				var chatID int64
				for _, update := range updates {
					if update.Message == nil { // ignore any non-Message Updates
						continue
					}

					if update.Message.From.UserName == username {
						log.Printf("Telebot chat initialized via message: [%s] %s", update.Message.From.UserName, update.Message.Text)
						chatID = update.Message.Chat.ID
						break
					}
				}
				if chatID > 0 {
					bot.ChatID = chatID
				}
			}

			time.Sleep(5 * time.Minute)
		}
	}()
}

func (bot *Telebot) Send(message string) {
	// TODO: if chatID <= 0, then put into wait-loop/queue in background until not <= 0 anymore
	msg := tgbotapi.NewMessage(bot.ChatID, message)
	bot.api.Send(msg)
}
