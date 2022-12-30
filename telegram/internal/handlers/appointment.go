package handlers

import (
	"fmt"
	"strconv"
	"time"

	objs "github.com/SakoDroid/telego/objects"
	"github.com/escalopa/govisa/telegram/core"
	"github.com/escalopa/govisa/telegram/internal/application"
	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
)

const (
	DateLayout    = "02 Jan 06 15:04 MST"
	DateDayLayout = "02 Jan 06"
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
	vType, ok := bh.readType(ch, u)
	if !ok {
		return
	}
	cva.VType = core.VType(vType)

	// Read City
	location, ok := bh.readLocation(ch, u)
	if !ok {
		return
	}
	cva.Location = core.Location(location)

	err = bh.uc.BookVisaAppointment(bh.ctx, u.Message.Chat.Id, cva)
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "An error occurred while booking your appointment, Please try again /book", err, 0)
		return
	}
	bh.simpleSend(u.Message.Chat.Id, "Your appointment has been booked successfully", 0)
}

func (bh *BotHandler) Dates(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	ch, err := bh.b.AdvancedMode().RegisterChannel(chatID, "message")
	if err != nil {
		return
	}

	// Read City
	location, ok := bh.readLocation(ch, u)
	if !ok {
		return
	}

	// Get Dates & Convert dates to md table
	dates, err := bh.uc.GetAvailableVisaAppointmentDates(bh.ctx, location)
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "Error", err, 0)
		return
	}
	table, err := tablify([]string{"Date"}, func() [][]string {
		var rows [][]string
		for _, d := range dates {
			rows = append(rows, []string{d.Format(DateDayLayout)})
		}
		return rows
	})
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "Error", err, 0)
		return
	}

	// Send table
	_, err = bh.b.SendMessage(
		u.Message.Chat.Id, toMarkdown(fmt.Sprintf("Available Dates in %s", location), table), "Markdown",
		u.Message.MessageId, false, false,
	)
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "Failed to make ", err, u.Message.MessageId)
	}
}

func (bh *BotHandler) Status(u *objs.Update) {
	// Get Dates & Convert dates to md table
	currAppt, err := bh.uc.GetCurrentVisaAppointment(bh.ctx, u.Message.Chat.Id)
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "Error", err, 0)
		return
	}

	// Convert status to table
	table, err := tablify([]string{"Date", "Type", "Location"}, func() [][]string {
		var rows [][]string
		if currAppt != nil {
			rows = append(rows, []string{currAppt.Date.Format(DateLayout), string(currAppt.VType), string(currAppt.Post)})
		}
		return rows
	})
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "Error", err, 0)
		return
	}

	// Send table
	_, err = bh.b.SendMessage(u.Message.Chat.Id, toMarkdown("Current Appointment Status", table), "Markdown", u.Message.MessageId, false, false)
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "Failed to send table", err, u.Message.MessageId)
	}
}

func (bh *BotHandler) History(u *objs.Update) {
	// Get previous appointments
	prevAppts, err := bh.uc.GetVisaAppointments(bh.ctx, u.Message.Chat.Id)
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "Error", err, 0)
		return
	}

	// Convert history to table
	table, err := tablify([]string{"Date", "Type", "Location", "Status"}, func() [][]string {
		var rows [][]string
		for _, appt := range prevAppts {
			rows = append(rows, []string{appt.Date.Format(DateDayLayout), string(appt.VType), string(appt.Post), string(appt.Status)})
		}
		return rows
	})
	if err != nil {
		bh.simpleError(u.Message.Chat.Id, "Error", err, 0)
		return
	}

	// Send table
	if _, err := bh.b.SendMessage(u.Message.Chat.Id, toMarkdown("Appointments History", table), "Markdown", u.Message.MessageId, false, false); err != nil {
		bh.simpleError(u.Message.Chat.Id, "Failed to send table", err, u.Message.MessageId)
	}
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
	bh.simpleSend(u.Message.Chat.Id, "Your appointment has been rescheduled successfully", 0)
}

func (bh *BotHandler) readLocation(ch *chan *objs.Update, u *objs.Update) (string, bool) {
	ckb := bh.b.CreateKeyboard(true, true, false, "")
	ckb.AddButton("Abuja", 1)
	ckb.AddButton("Lagos", 1)
	ckb.AddButton("Abort", 2)
	bh.b.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Choose the appointment city", "", 0, false, false, nil, true, true, ckb)
	u = <-*ch
	if bh.checkAbort(u, "Book Appointment") {
		return "", false
	}
	location := u.Message.Text
	return location, true
}

func (bh *BotHandler) readType(ch *chan *objs.Update, u *objs.Update) (string, bool) {
	tkb := bh.b.CreateKeyboard(true, true, false, "")
	tkb.AddButton("F1", 1)
	tkb.AddButton("F2", 1)
	tkb.AddButton("F3", 1)
	tkb.AddButton("Abort", 2)
	bh.b.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Choose type of visa are you applying for", "", 0, false, false, nil, true, true, tkb)
	u = <-*ch
	if bh.checkAbort(u, "Book Appointment") {
		return "", false
	}
	vType := u.Message.Text
	return vType, true
}

func tablify(columns []string, data func() [][]string) (string, error) {
	table, err := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build(columns...).
		Format(data())
	return table, err
}
