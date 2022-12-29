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
}

func (bh *BotHandler) Help(u *objects.Update) {
	bh.simpleSendMD(u.Message.Chat.Id, `
		
		Welcome to the NG USA Visa Bot 🤖, Book your appointment with the USA Embassy in Nigeria 🇺🇸🇳🇬

		Here are the available commands: 👇

		/login - Login to your account 🔑
		/book - Book a new appointment 📅
		
		/dates - Get available dates ⏱
		/status - Get your current appointment status 📝
		/history - Get your appointment history 📚
		
		/reschedule - Reschedule your current appointment 🗓️
		/cancel - Cancel your current appointment ❌

		/help - Show this message 📖
	`, u.Message.MessageId)
}

func (bh *BotHandler) Public(u *objects.Update) {
	bh.l.Println("Public message received")
	bh.simpleSend(u.Message.Chat.Id, "Bot is not avaliable out the scope of private chats", u.Message.MessageId)
}

func (bh *BotHandler) Unknow(u *objects.Update) {
	bh.simpleSend(u.Message.Chat.Id, "Unknow command, please use /help to see the available commands", u.Message.MessageId)
}

// SimpleSend sends a simple message
func (bh *BotHandler) simpleSend(chatID int, text string, replyTo int) {
	bh.send(chatID, text, "", replyTo)
}

// simpleSendMD sends a message with markdown support
func (bh *BotHandler) simpleSendMD(chatID int, text string, replyTo int) {
	bh.send(chatID, text, "Markdown", replyTo)
}

// send sends a message to the chat
func (bh *BotHandler) send(chatID int, text string, parseMode string, replyTo int) {
	_, err := bh.b.SendMessage(chatID, text, parseMode, replyTo, false, false)
	if err != nil {
		bh.l.Println(err)
	}
}

func (bh *BotHandler) checkCancel(u *objects.Update) bool {
	if u.Message.Text == "/abort" {
		bh.simpleSend(u.Message.Chat.Id, "Operation aborted", 0)
		return true
	}
	return false
}
