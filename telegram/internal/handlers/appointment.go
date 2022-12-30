package handlers

import (
	"strconv"
	"time"

	objs "github.com/SakoDroid/telego/objects"
	"github.com/escalopa/govisa/telegram/core"
	"github.com/escalopa/govisa/telegram/internal/application"
)

func (bh *BotHandler) Book(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	ch, err := bh.b.AdvancedMode().RegisterChannel(chatID, "message")
	if err != nil {
		return
	}

	var cva application.CreateVisaAppointment
	// Read Date
	// TODO: Give option to choose date from calendar
	cva.Date = time.Now()

	// Read Type
	tkb := bh.b.CreateKeyboard(true, true, false, "")
	tkb.AddButton("F1", 1)
	tkb.AddButton("F2", 1)
	tkb.AddButton("F3", 1)
	bh.b.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Choose type of visa are you applying for", "", 0, false, false, nil, true, true, tkb)
	u = <-*ch
	cva.Type = core.Type(u.Message.Text)

	// Read City
	ckb := bh.b.CreateKeyboard(true, true, false, "")
	ckb.AddButton("Abuja", 1)
	ckb.AddButton("Lagos", 1)
	bh.b.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Choose the appointment city", "", 0, false, false, nil, true, true, ckb)
	u = <-*ch
	cva.Location = core.Location(u.Message.Text)

	err = bh.uc.BookVisaAppointment(bh.ctx, u.Message.Chat.Id, cva)
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "An error occurred while booking your appointment, Please try again /book", err, 0)
		return
	}
	bh.simpleSend(u.Message.Chat.Id, "Your appointment has been booked successfully", 0)
}

func (bh *BotHandler) Dates(u *objs.Update) {
	bh.uc.GetAvailableVisaAppointmentDates(bh.ctx)
}

func (bh *BotHandler) Status(u *objs.Update) {
	bh.uc.GetCurrentVisaAppointment(bh.ctx, u.Message.Chat.Id)
}

func (bh *BotHandler) History(u *objs.Update) {
	bh.uc.GetVisaAppointments(bh.ctx, u.Message.Chat.Id)
}

func (bh *BotHandler) Cancel(u *objs.Update) {
	err := bh.uc.CancelVisaAppointment(bh.ctx, u.Message.Chat.Id)
	if err != nil {
		bh.simpleSend(u.Message.Chat.Id, "An error occurred while cancelling your appointment", 0)
		return
	}
	bh.simpleSend(u.Message.Chat.Id, "Your appointment has been cancelled", 0)
}

func (bh *BotHandler) Reschedule(u *objs.Update) {

}
