package application

import (
	"context"
	"errors"
	"time"

	"github.com/escalopa/govisa/pkg/govisa"
	"github.com/escalopa/govisa/telegram/core"
	validate "github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
)

var (
	Validate = validate.New()

	ErrUserNotFound = errors.New("ops!, you are not logged in, Use /login to login")
)

func init() {
	var err error
	err = Validate.RegisterValidation("vlocation", isValidLocation)
	govisa.CheckError(err)
	err = Validate.RegisterValidation("vcity", isValidType)
	govisa.CheckError(err)
}

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
	Date     time.Time  `validate:"required"`
	VType    core.VType `validate:"required,vtype"`
	Location core.City  `validate:"required,vcity"`
}

func New(uc UserCache, srv Server, enc Encryptor) *UseCase {
	return &UseCase{uc: uc, srv: srv, enc: enc}
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
	user, err := u.getUser(ctx, userID)
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

func (u *UseCase) RescheduleVisaAppointment(ctx context.Context, userID int, cva CreateVisaAppointment) error {
	user, err := u.getUser(ctx, userID)
	if err != nil {
		return err
	}
	// TODO: Validate struct
	// err = Validate.Struct(cva)
	// if err != nil {
	// 	return err
	// }
	err = u.srv.RescheduleVisaAppointment(user.ServerUserID, cva)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) CancelVisaAppointment(ctx context.Context, userID int) error {
	// Get user from cache
	user, err := u.getUser(ctx, userID)
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

func (u *UseCase) GetAvailableVisaAppointmentDates(ctx context.Context, userID int, city string) ([]time.Time, error) {
	user, err := u.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return u.srv.GetAvailableVisaAppointmentDates(user.ServerUserID, city)
}

func (u *UseCase) GetCurrentVisaAppointment(ctx context.Context, userID int) (*core.VisaAppointment, error) {
	// Get user from cache
	user, err := u.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get visa appointments from server
	return u.srv.GetCurrentVisaAppointment(user.ServerUserID)
}

// GetVisaAppointments returns all the previous visa appointments for a user
func (u *UseCase) GetVisaAppointments(ctx context.Context, userID int) ([]core.VisaAppointment, error) {
	// Get user from cache
	user, err := u.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get visa appointments from server
	return u.srv.GetVisaAppointments(user.ServerUserID)
}

func (u *UseCase) getUser(ctx context.Context, userID int) (*core.User, error) {
	user, err := u.uc.GetUserByID(ctx, userID)
	if err != nil {
		if err == redis.Nil {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func isValidLocation(fl validate.FieldLevel) bool {
	l1 := fl.Field().String()
	for _, l2 := range core.Cities {
		if l2 == l1 {
			return true
		}
	}
	return false
}

func isValidType(fl validate.FieldLevel) bool {
	t1 := fl.Field().String()
	for _, t2 := range core.Types {
		if t2 == t1 {
			return true
		}
	}
	return false
}
