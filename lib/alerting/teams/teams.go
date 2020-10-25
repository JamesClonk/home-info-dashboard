package teams

import (
	"sync"

	"github.com/JamesClonk/home-info-dashboard/lib/env"
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
)

var (
	teams *Teams = nil
	once  sync.Once
)

type Teams struct {
	*sync.Mutex
	API        goteamsnotify.API
	WebhookURL string
}

func Get() *Teams {
	return teams
}

func Init() {
	once.Do(func() {
		api := goteamsnotify.NewClient()
		webhookUrl := env.MustGet("TEAMS_WEBHOOK_URL")

		teams = &Teams{&sync.Mutex{}, api, webhookUrl}
	})
}

func (teams *Teams) Send(message string) error {
	teams.Lock()
	defer teams.Unlock()

	msg := goteamsnotify.NewMessageCard()
	//msg.Title = "Home Automation Dashboard"
	//msg.Text = message
	msg.Summary = "Home Automation Dashboard - Alert Notification"
	msg.ThemeColor = "#cc3333"
	if err := msg.AddSection(&goteamsnotify.MessageCardSection{
		ActivityImage:    "https://home-info.jamesclonk.io/images/smart_temperature.png",
		ActivityTitle:    "Home Automation Dashboard",
		ActivitySubtitle: "Alert!",
		ActivityText:     message,
		Markdown:         true,
	}); err != nil {
		return err
	}

	return teams.API.Send(teams.WebhookURL, msg)
}
