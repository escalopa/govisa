package handlers

import (
	"context"

	bt "github.com/SakoDroid/telego"
	"github.com/SakoDroid/telego/objects"
	"github.com/escalopa/govisa/telegram/internal/application"
)

type BotHandler struct {
	uc *application.UseCase
	b  *bt.Bot
}

func NewBotHandler(bot *bt.Bot, uc *application.UseCase, ctx context.Context) *BotHandler {
	return &BotHandler{b: bot, uc: uc}
}

func (bh *BotHandler) Register() {
	bh.b.AddHandler("/login", bh.Login)
	bh.b.AddHandler("/book", bh.Book)
	bh.b.AddHandler("/dates", bh.Dates)
	bh.b.AddHandler("/status", bh.Status)
	bh.b.AddHandler("/history", bh.History)
	bh.b.AddHandler("/cancel", bh.Cancel)
	bh.b.AddHandler("/reschedule", bh.Reschedule)
	bh.b.AddHandler("/help", bh.Help)
}

func (bh *BotHandler) Help(u *objects.Update) {

}
