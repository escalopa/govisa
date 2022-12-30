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

func (s *Server) Login(email, password string) (int64, error) {
	return 1, nil
}

func (s *Server) BookVisaAppointment(userID int64, cva application.CreateVisaAppointment) error {
	return nil
}

func (s *Server) CancelVisaAppointment(userID int64) error {
	return nil
}

func (s *Server) GetAvailableVisaAppointmentDates(city string) ([]time.Time, error) {
	return []time.Time{time.Now(), time.Now().Add(time.Hour * 24)}, nil
}

func (s *Server) GetCurrentVisaAppointment(userID int64) (*core.VisaAppointment, error) {
	return &core.VisaAppointment{
		Applicant: "John Doe",
		Post:      "Lagos",
		VType:     "F1",
		Status:    "Scheduled",
		Date:      time.Now().Add(time.Hour * 24 * 7), // 7 days from now
	}, nil
}

func (s *Server) GetVisaAppointments(userID int64) ([]core.VisaAppointment, error) {
	return []core.VisaAppointment{
		{
			Applicant: "John Doe",
			Post:      "Lagos",
			VType:     "F1",
			Status:    "Scheduled",
			Date:      time.Now().Add(time.Hour * 24 * 7), // 7 days from now
		},
		{
			Applicant: "John Doe",
			Post:      "Abuja",
			VType:     "F2",
			Status:    "Canceled",
			Date:      time.Now().Add(-time.Hour * 24 * 14), // 14 days from now
		},
	}, nil
}
