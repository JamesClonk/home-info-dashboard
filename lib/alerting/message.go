package alerting

import (
	"github.com/JamesClonk/home-info-dashboard/lib/alerting/slack"
	"github.com/JamesClonk/home-info-dashboard/lib/alerting/telebot"
	"github.com/JamesClonk/home-info-dashboard/lib/database"
)

var (
	messengers = make([]Messenger, 0)
)

type Messenger interface {
	Send(string) error
}

func Init(hdb database.HomeInfoDB) {
	slack.Init(hdb)
	telebot.Init(hdb)

	messengers = append(messengers, slack.Get())
	messengers = append(messengers, telebot.Get())
}

func Send(message string) error {
	for _, messenger := range messengers {
		if err := messenger.Send(message); err != nil {
			return err
		}
	}
	return nil
}
