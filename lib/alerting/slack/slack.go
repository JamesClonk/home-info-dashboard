package slack

import (
	"strings"
	"sync"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/env"
	slack_api "github.com/nlopes/slack"
)

var (
	slack *Slack = nil
	once  sync.Once
)

type Slack struct {
	*sync.Mutex
	hdb     database.HomeInfoDB
	API     *slack_api.Client
	Channel string
}

func Get() *Slack {
	return slack
}

func Init(hdb database.HomeInfoDB) {
	once.Do(func() {
		api := slack_api.New(env.MustGet("SLACK_TOKEN"), slack_api.OptionDebug(env.Get("SLACK_DEBUG", "false") == "true"))

		channel := env.MustGet("SLACK_CHANNEL")
		channel = strings.ReplaceAll(channel, "USER_", "@")
		channel = strings.ReplaceAll(channel, "CHANNEL_", "#")

		slack = &Slack{&sync.Mutex{}, hdb, api, channel}
	})
}

func (slack *Slack) Send(message string) error {
	slack.Lock()
	defer slack.Unlock()

	attachment := slack_api.Attachment{
		//Pretext:    "Notification:",
		Text:       message,
		AuthorName: "Home Automation Dashboard",
		AuthorIcon: "https://home-info.scapp.io/images/smart_temperature.png",
		Color:      "#cc3333",
	}
	_, _, err := slack.API.PostMessage(
		slack.Channel,
		slack_api.MsgOptionAttachments(attachment),
		//slack_api.MsgOptionIconURL("https://home-info.scapp.io/images/smart_temperature.png"),
		slack_api.MsgOptionIconEmoji("house_with_garden"),
		slack_api.MsgOptionUsername("home-info"),
	)

	// retry?
	if err != nil {
		_, _, err = slack.API.PostMessage(
			slack.Channel,
			slack_api.MsgOptionAttachments(attachment),
			//slack_api.MsgOptionIconURL("https://home-info.scapp.io/images/smart_temperature.png"),
			slack_api.MsgOptionIconEmoji("house_with_garden"),
			slack_api.MsgOptionUsername("home-info"),
		)
	}
	return err
}
