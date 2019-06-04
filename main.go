package main

import (
	"fmt"
	"github.com/jinzhu/configor"
	"log"
	"regexp"
	"time"
	tb "gopkg.in/tucnak/telebot.v2"
)

var Config = struct {
	Bot   struct {
		Token string
		Name  string
		Id    int
	}
}{}
var arabic = regexp.MustCompile(`[\x{600}-\x{6FF}]+`)
var rtl = regexp.MustCompile(`[\x{0590}-\x{05ff}\x{0600}-\x{06ff}]+`)

func main()  {
	configor.Load(&Config, "config.json")
	b, err := tb.NewBot(tb.Settings{
		Token:  Config.Bot.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Authorized on account %s", b.Me.Username)

	b.Handle("/start", func(m *tb.Message) {
		if (m.Private() == true) {
			options := new(tb.SendOptions)
			options.ParseMode = tb.ModeHTML
			options.ReplyTo = m

			message := fmt.Sprintf("Hello %s, I am @GoClearBot\n", m.Sender.FirstName)
			message += "Just add me to chat and I'll start delete stickers, \"joined\" and Arabic messages\n"
			message += "I'm written in <i>Go</i> lang\n"
			message += "Our channel @MagicBots\n"
			message += "Feel free @temamagic\n"
			message += "Source code: <a href=\"https://github.com/magicbotsru/GoClearBot\">Github</a>\n"

			b.Send(m.Sender, message, options)
		}
	})
	b.Handle("/ping", func(m *tb.Message) {
		options := new(tb.SendOptions)
		options.ParseMode = tb.ModeHTML
		options.ReplyTo = m
		message := "<b>Pong</b>"
		b.Send(m.Chat, message, options)
	})
	b.Handle(tb.OnSticker, func(m *tb.Message) {
		if (m.Private() != true) {
			b.Delete(m)
		}
	})
	b.Handle(tb.OnUserJoined, func(m *tb.Message) {
		b.Delete(m)
		go filterArabic(m,b)
	})
	b.Handle(tb.OnText, func(m *tb.Message) {
		if (m.FromGroup() == true) {
			go filterArabic(m,b)
		}
	})
	b.Start()
}

func filterArabic(m *tb.Message,b *tb.Bot)  {
	isArabic := arabic.MatchString(fmt.Sprintf("%s %s %s", m.Sender.FirstName, m.Sender.LastName, m.Text))
	isRTL := rtl.MatchString(fmt.Sprintf("%s %s %s", m.Sender.FirstName, m.Sender.LastName, m.Text))
	if isArabic || isRTL {
		log.Printf("Arabic found")
		b.Delete(m)
		b.Ban(m.Chat,&tb.ChatMember{
			Rights:          tb.Rights{},
			User:            m.Sender,
			Role:            "",
			RestrictedUntil: 0,
		})
	}
}