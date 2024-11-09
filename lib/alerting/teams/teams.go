package teams

import (
	"sync"

	"github.com/JamesClonk/home-info-dashboard/lib/env"
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

var (
	teams *Teams = nil
	once  sync.Once
)

type Teams struct {
	*sync.Mutex
	API        *goteamsnotify.TeamsClient
	WebhookURL string
}

func Get() *Teams {
	return teams
}

func Init() {
	once.Do(func() {
		api := goteamsnotify.NewTeamsClient()
		webhookUrl := env.MustGet("TEAMS_WEBHOOK_URL")

		teams = &Teams{&sync.Mutex{}, api, webhookUrl}
	})
}

func (teams *Teams) Send(message string) error {
	teams.Lock()
	defer teams.Unlock()

	// msg, err := adaptivecard.NewSimpleMessage(
	// 	message,
	// 	"Home Automation Dashboard - Alert Notification",
	// 	true,
	// )
	// if err != nil {
	// 	return err
	// }

	card, err := adaptivecard.NewTextBlockCard(message, "Home Automation Dashboard - Alert Notification", true)
	if err != nil {
		return err
	}

	targetURL := "https://home-info.jamesclonk.io/dashboard"
	targetURLDesc := "Home-Info Dashboard"
	urlAction, err := adaptivecard.NewActionOpenURL(targetURL, targetURLDesc)
	if err != nil {
		return err
	}
	if err := card.AddAction(true, urlAction); err != nil {
		return err
	}

	// Create Message from Card
	msg, err := adaptivecard.NewMessageFromCard(card)
	if err != nil {
		return err
	}

	// msg := goteamsnotify.NewMessageCard()
	// //msg.Title = "Home Automation Dashboard"
	// //msg.Text = message
	// msg.Summary = "Home Automation Dashboard - Alert Notification"
	// msg.ThemeColor = "#cc3333"
	// if err := msg.AddSection(&goteamsnotify.MessageCardSection{
	// 	ActivityImage:    "https://home-info.jamesclonk.io/images/smart_temperature.png",
	// 	ActivityTitle:    "Home Automation Dashboard",
	// 	ActivitySubtitle: "Alert Notification!",
	// 	//ActivityText:     message,
	// 	Text:     message,
	// 	Markdown: true,
	// }); err != nil {
	// 	return err
	// }

	return teams.API.Send(teams.WebhookURL, msg)
}
