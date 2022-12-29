package handlers

import (
	"context"
	"log"

	bt "github.com/SakoDroid/telego"
	"github.com/SakoDroid/telego/objects"
	"github.com/escalopa/govisa/telegram/internal/application"
)

type BotHandler struct {
	uc  *application.UseCase
	b   *bt.Bot
	l   *log.Logger
	ctx context.Context
}

func NewBotHandler(bot *bt.Bot, uc *application.UseCase, l *log.Logger, ctx context.Context) *BotHandler {
	return &BotHandler{b: bot, uc: uc, l: l, ctx: ctx}
}

func (bh *BotHandler) Register() {
	bh.b.AddHandler("/login", bh.Login, "private")
	bh.b.AddHandler("/book", bh.Book, "private")
	bh.b.AddHandler("/dates", bh.Dates, "private")
	bh.b.AddHandler("/status", bh.Status, "private")
	bh.b.AddHandler("/history", bh.History, "private")
	bh.b.AddHandler("/cancel", bh.Cancel, "private")
	bh.b.AddHandler("/reschedule", bh.Reschedule, "private")
	bh.b.AddHandler("/help", bh.Help, "private")
	bh.b.AddHandler("*", bh.Public, "supergroup", "group")
}

func (bh *BotHandler) Public(u *objects.Update) {

}

func (bh *BotHandler) Help(u *objects.Update) {

}

func (bh *BotHandler) SimpleSend(chatID int, text string, replyTo int) {
	_, err := bh.b.SendMessage(chatID, text, "", replyTo, false, false)
	if err != nil {
		bh.l.Println(err)
	}
}
