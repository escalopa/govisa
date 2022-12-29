package application

import (
	"context"
	"time"

	"github.com/escalopa/govisa/telegram/core"
	validate "github.com/go-playground/validator/v10"
)

var Validate = validate.New()

type UseCase struct {
	uc  UserCache
	srv Server
	enc Encryptor
}

type CreateUser struct {
	ID       int    `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type CreateVisaAppointment struct {
	Date     time.Time     `validate:"required"`
	Type     core.Type     `validate:"required,oneof=F1"`
	Location core.Location `validate:"required,oneof=Abuja Lagos"`
}

func New(uc UserCache, srv Server, enc Encryptor) (*UseCase, error) {
	return &UseCase{uc: uc, srv: srv, enc: enc}, nil
}

func (u *UseCase) Login(ctx context.Context, cu CreateUser) error {
	err := Validate.Struct(cu)
	if err != nil {
		return err
	}
	// Login to server to confirm credentials
	srvUserID, err := u.srv.Login(cu.Email, cu.Password)
	if err != nil {
		return err
	}

	// Save user in cache
	encryptedPassword, err := u.enc.Encrypt(cu.Password)
	if err != nil {
		return err
	}
	user := &core.User{
		ID:           cu.ID,
		ServerUserID: srvUserID,
		Email:        cu.Email,
		Password:     encryptedPassword,
	}
	err = u.uc.SaveUserByID(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) BookVisaAppointment(ctx context.Context, userID int, cva CreateVisaAppointment) error {
	err := Validate.Struct(cva)
	if err != nil {
		return err
	}
	// Get user from cache
	user, err := u.uc.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Book visa appointment in server
	err = u.srv.BookVisaAppointment(user.ServerUserID, cva)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) GetAvailableVisaAppointmentDates(ctx context.Context) ([]time.Time, error) {
	return u.srv.GetAvailableVisaAppointmentDates()
}

func (u *UseCase) CancelVisaAppointment(ctx context.Context, userID int) error {
	// Get user from cache
	user, err := u.uc.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Cancel visa appointment in server
	err = u.srv.CancelVisaAppointment(user.ServerUserID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) GetCurrentVisaAppointment(ctx context.Context, userID int) (*core.VisaAppointment, error) {
	// Get user from cache
	user, err := u.uc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get visa appointments from server
	return u.srv.GetCurrentVisaAppointment(user.ServerUserID)
}

func (u *UseCase) GetVisaAppointments(ctx context.Context, userID int) ([]core.VisaAppointment, error) {
	// Get user from cache
	user, err := u.uc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get visa appointments from server
	return u.srv.GetVisaAppointments(user.ServerUserID)
}
