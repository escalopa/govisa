package handlers

import (
	"strconv"

	objs "github.com/SakoDroid/telego/objects"
	"github.com/escalopa/govisa/telegram/internal/application"
)

func (bh *BotHandler) Login(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	ch, err := bh.b.AdvancedMode().RegisterChannel(chatID, "message")
	if err != nil {
		bh.l.Println(err)
		return
	}

	// Request email
	bh.simpleSend(u.Message.Chat.Id, "Please enter your email", 0)
	u = <-*ch
	if bh.checkCancel(u) {
		return
	}
	email := u.Message.Text

	// Read password
	bh.simpleSend(u.Message.Chat.Id, "Please enter your password", 0)
	u = <-*ch
	if bh.checkCancel(u) {
		return
	}
	password := u.Message.Text

	// Call login usecase
	err = bh.uc.Login(bh.ctx, application.CreateUser{
		ID:       u.Message.From.Id,
		Email:    email,
		Password: password,
	})

	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "Wrong credentials, Pleese try again /loign", err, 0)
		return
	}
	bh.simpleSend(u.Message.Chat.Id, "You have logged in successfully", 0)
}
