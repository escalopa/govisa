package handlers

import (
	"context"
	"fmt"
	"log"

	bt "github.com/SakoDroid/telego"
	"github.com/SakoDroid/telego/objects"
	"github.com/escalopa/govisa/pkg/govisa"
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
	govisa.CheckError(bh.b.AddHandler("/start", bh.Start, "all"))
	govisa.CheckError(bh.b.AddHandler("/login", bh.Login, "private"))
	govisa.CheckError(bh.b.AddHandler("/book", bh.Book, "private"))
	govisa.CheckError(bh.b.AddHandler("/dates", bh.Dates, "private"))
	govisa.CheckError(bh.b.AddHandler("/status", bh.Status, "private"))
	govisa.CheckError(bh.b.AddHandler("/history", bh.History, "private"))
	govisa.CheckError(bh.b.AddHandler("/cancel", bh.Cancel, "private"))
	govisa.CheckError(bh.b.AddHandler("/reschedule", bh.Reschedule, "private"))
	govisa.CheckError(bh.b.AddHandler("/help", bh.Help, "private"))
}

func (bh *BotHandler) Start(u *objects.Update) {
	bh.simpleSend(u.Message.Chat.Id, `
		Welcome to the NG USA Visa Bot 🤖, Book your appointment with the USA Embassy in Nigeria 🇺🇸🇳🇬 in minutes
	`, 0)
	bh.Help(u)
}

func (bh *BotHandler) Help(u *objects.Update) {
	bh.simpleSend(u.Message.Chat.Id, `
		You can use the following commands to interact with the bot, List of available commands are below: 👇

		/login - Login to your account 🔑
		/book - Book a new appointment 📅
		
		/dates - Get available dates ⏱
		/status - Get your current appointment status 📝
		/history - Get your appointment history 📚
		
		/reschedule - Reschedule your current appointment 🗓️
		/cancel - Cancel your current appointment ❌

		/help - Show this message 📖 
	`, 0)
}

func (bh *BotHandler) Public(u *objects.Update) {
	bh.l.Println("Public message received")
	bh.simpleSend(u.Message.Chat.Id, "Bot is not avaliable out the scope of private chats", u.Message.MessageId)
}

// SimpleSend sends a simple message
func (bh *BotHandler) simpleSend(chatID int, text string, replyTo int) {
	_, err := bh.b.SendMessage(chatID, text, "", replyTo, false, false)
	if err != nil {
		go bh.l.Println(err)
	}
}

func (bh *BotHandler) simpleError(chatID int, msg string, err error, replyTo int) {
	go bh.l.Printf("chatID: %d, Error: %s, Msg: %s", chatID, err, msg)
	bh.simpleSend(chatID, msg, replyTo)
}

func (bh *BotHandler) checkAbort(u *objects.Update, operation string) bool {
	if u.Message.Text == "Abort" {
		_, err := bh.b.SendMessage(u.Message.Chat.Id, fmt.Sprintf("Operation: <b>%s</b> has been aborted", operation), "HTML", 0, false, false)
		if err != nil {
			go bh.l.Println("Failed to send message", err)
		}
		return true
	}
	return false
}
