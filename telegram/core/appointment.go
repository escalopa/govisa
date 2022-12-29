package core

import (
	"time"
)

const (
	TypeF1 = "F-1" // F1 Student Visa

	StatusScheduled = "Scheduled" // Appointment is scheduled
	StatusCanceled  = "Canceled"  // Appointment is canceled

	LocationLagos = "Lagos" // Appointment is in Lagos
	LocationAbuja = "Abuja" // Appointment is in Abuja
)

type Type string
type Status string
type Location string

type VisaAppointment struct {
	Applicant string    // Name of the applicant
	Post      Location  // Location of the appointment
	Type      string    // Type of visa, Ex: "F1 Student Visa"
	Status    Status    // Status of the appointment, Ex: "Scheduled or Canceled"
	Date      time.Time // Date of the appointment
}

func (v Type) IsValidType() bool {
	return v == TypeF1
}

func (v Status) IsValidStatus() bool {
	return v == StatusScheduled || v == StatusCanceled
}

func (v Location) IsValidLocation() bool {
	return v == LocationLagos || v == LocationAbuja
}
