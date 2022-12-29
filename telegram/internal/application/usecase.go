package application

import (
	"context"
	"time"

	"github.com/escalopa/govisa/telegram/core"
)

type UseCase struct {
	uc  UserCache
	srv Server
	enc Encryptor
}

type CreateUser struct {
	ID       string
	Email    string
	Password string
}

type CreateVisaAppointment struct {
	Date     time.Time
	Type     core.Type
	Location core.Location
}

func New(uc UserCache, srv Server, enc Encryptor) (*UseCase, error) {
	return &UseCase{uc: uc, srv: srv, enc: enc}, nil
}

func (u *UseCase) Login(ctx context.Context, cu CreateUser) error {
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

func (u *UseCase) BookVisaAppointment(ctx context.Context, userID string, cva CreateVisaAppointment) error {
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

func (u *UseCase) CancelVisaAppointment(ctx context.Context, userID string) error {
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

func (u *UseCase) GetCurrentVisaAppointment(ctx context.Context, userID string) (*core.VisaAppointment, error) {
	// Get user from cache
	user, err := u.uc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get visa appointments from server
	return u.srv.GetCurrentVisaAppointment(user.ServerUserID)
}

func (u *UseCase) GetVisaAppointments(ctx context.Context, userID string) ([]core.VisaAppointment, error) {
	// Get user from cache
	user, err := u.uc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get visa appointments from server
	return u.srv.GetVisaAppointments(user.ServerUserID)
}
