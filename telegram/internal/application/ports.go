package application

import (
	"context"
	"time"

	"github.com/escalopa/govisa/telegram/core"
)

type Server interface {
	Login(email, password string) (int64, error)
	GetCurrentVisaAppointment(userID int64) (*core.VisaAppointment, error)
	GetVisaAppointments(userID int64) ([]core.VisaAppointment, error)
	GetAvailableVisaAppointmentDates(userID int64, city string) ([]time.Time, error)
	BookVisaAppointment(userID int64, cva CreateVisaAppointment) error
	RescheduleVisaAppointment(userID int64, cva CreateVisaAppointment) error
	CancelVisaAppointment(userID int64) error
}

type UserCache interface {
	GetUserByID(ctx context.Context, id int) (*core.User, error)
	SaveUserByID(ctx context.Context, user *core.User) error
}

type Encryptor interface {
	Encrypt(text string) (string, error)
	Decrypt(text string) (string, error)
}
