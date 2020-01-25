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
	*sync.Mutex
	hdb    database.HomeInfoDB
	API    *tgbotapi.BotAPI
	ChatID int64
}

func Get() *Telebot {
	return telebot
}

func Init(hdb database.HomeInfoDB) {
	once.Do(func() {
		telebot = &Telebot{&sync.Mutex{}, hdb, nil, 0}
		telebot.UpdateAPI()      // setup api
		telebot.WatchForChatID() // background checking for chatID
	})
}

func (bot *Telebot) WatchForChatID() {
	go func() {
		for {
			if bot.ChatID <= 0 {
				if bot.API == nil {
					bot.UpdateAPI()
				}
				bot.UpdateChatID()
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}

func (bot *Telebot) UpdateChatID() {
	bot.Lock()
	defer bot.Unlock()

	if bot.API == nil {
		log.Println("Telebot API is nil, cannot update chat ID")
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.API.GetUpdates(u)
	if err != nil {
		log.Fatalf("Telebot API error: %v\n", err)
	}

	username := env.MustGet("TELEBOT_USERNAME")
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

func (bot *Telebot) UpdateAPI() {
	bot.Lock()
	defer bot.Unlock()

	token := env.MustGet("TELEBOT_TOKEN")
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("Telebot API error: %v\n", err)
		return
	}

	api.Debug = env.Get("TELEBOT_DEBUG", "false") == "true"
	log.Printf("Telebot authorized on account [%s]", api.Self.UserName)
	bot.API = api
}

func (bot *Telebot) Send(message string) error {
	bot.Lock()
	defer bot.Unlock()

	// TODO: if chatID <= 0, then put into wait-loop/queue in background until not <= 0 anymore
	msg := tgbotapi.NewMessage(bot.ChatID, message)
	_, err := bot.API.Send(msg)
	return err
}
