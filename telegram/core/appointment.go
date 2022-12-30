package core

import (
	"time"
)

var (
	Types     = []string{"F1", "F2", "F3"} // Appointment is for F1 F2 F3
	Locations = []string{"Lagos", "Abuja"} // Appointment is in Lagos or Abuja
)

type Type string
type Location string

type VisaAppointment struct {
	Applicant string    // Name of the applicant
	Post      Location  // Location of the appointment
	Type      string    // Type of visa, Ex: "F1 Student Visa"
	Status    string    // Status of the appointment, Ex: "Scheduled or Canceled"
	Date      time.Time // Date of the appointment
}
