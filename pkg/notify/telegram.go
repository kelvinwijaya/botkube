// Copyright (c) 2019 InfraCloud Technologies
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package notify

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/infracloudio/botkube/pkg/config"
	"github.com/infracloudio/botkube/pkg/events"
	"github.com/infracloudio/botkube/pkg/log"
)

// Telegram contains URL
type Telegram struct {
	Token     		string
	BotID			string
	ChatID 			string
	NotifType 		config.NotifType
}

// NewTelegram returns new Telegram object
func NewTelegram(c config.Telegram) Notifier {
	return &Telegram{
		Token: c.Token
		BotID: c.BotID,
		ChatID: c.ChatID,
		NotifType: c.NotifType,
	}
}

// SendEvent sends event notification to Telegram Channel
func (t *Telegram) SendEvent(event events.Event) (err error) {
	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		log.Error("error creating Telegram session,", err)
	}

	text := formatTelegramMessage(event, t.NotifType)

	msg := tgbotapi.NewMessage(t.ChatID, text)
	msg.ParseMode = "markdown"
	msg.
	if _, err := bot.Send(msg); err != nil {
		log.Errorf("Error in sending message: %+v", err)
		return err
	}
	log.Debugf("Event successfully sent to channel %s", t.ChannelID)
	return nil
}

func formatTelegramMessage(event events.Event, notifyType config.NotifType) string {

	var text string

	switch notifyType {
	case config.LongNotify:
		// generate Long notification message
		text = telegramLongNotification(event)

	case config.ShortNotify:
		// generate Short notification message
		fallthrough

	default:
		// generate Short notification message
		text = telegramShortNotification(event)
	}

	return text

}

func telegramLongNotification(event events.Event) string {
	text := fmt.Sprintf("*%s*", event.Title) + "\\"
	text += "Kind: " + event.Kind + "\\"
	text += "Name: " + event.Name + "\\"
	if event.Namespace != "" {
		text += "Namespace: " + event.Namespace + "\\"
	}
	if event.Reason != "" {
		text += "Reason: " + event.Reason + "\\"
	}
	if len(event.Messages) > 0 {
		message := ""
		for _, m := range event.Messages {
			message += fmt.Sprintf("%s\\", m)
		}
		text += "Message: " + message + "\\"
	}
	if event.Action != "" {
		text += "Action: " + event.Action + "\\"
	}
	if len(event.Recommendations) > 0 {
		rec := ""
		for _, r := range event.Recommendations {
			rec += fmt.Sprintf("%s\\", r)
		}
		text += "Recommendations: " + rec + "\\"
	}
	if len(event.Warnings) > 0 {
		warn := ""
		for _, w := range event.Warnings {
			warn += fmt.Sprintf("%s\\", w)
		}
		text += "Warnings: " + warn + "\\"
	}
	text += "EOL"
	return text
}

func telegramShortNotification(event events.Event) string {
	text := fmt.Sprintf("*%s*", event.Title) + "\\"
	text := "Description: " + FormatShortMessage(event) + "\\"
	text += "EOL"
	return text
}