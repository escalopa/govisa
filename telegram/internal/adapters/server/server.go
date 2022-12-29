package server

import (
	"time"

	"github.com/escalopa/govisa/telegram/core"
	"github.com/escalopa/govisa/telegram/internal/application"
)

type Server struct {
	endpoint string
}

func NewServer(endpoint string) (*Server, error) {
	return &Server{endpoint: endpoint}, nil
}

func (s *Server) Connect() error {
	return nil
}

func (s *Server) Login(email, password string) (int64, error) {
	return 1, nil
}

func (s *Server) BookVisaAppointment(userID int64, cva application.CreateVisaAppointment) error {
	return nil
}

func (s *Server) CancelVisaAppointment(userID int64) error {
	return nil
}

func (s *Server) GetAvailableVisaAppointmentDates() ([]time.Time, error) {
	return nil, nil
}

func (s *Server) GetCurrentVisaAppointment(userID int64) (*core.VisaAppointment, error) {
	return nil, nil
}

func (s *Server) GetVisaAppointments(userID int64) ([]core.VisaAppointment, error) {
	return nil, nil
}
