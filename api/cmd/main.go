package main

import (
	"log"

	"github.com/escalopa/govisa/api/internal"
)

func main() {

	cookies := internal.Login(internal.Credentials{
		Email:    "mokpara5@gmail.com",
		Password: "Barrister2nd",
	})

	body, err := internal.GetAppointmentHistory(cookies)
	if err != nil {
		log.Fatal(err)
	}

	if err = internal.HTMLToFile("Login.html", string(body)); err != nil {
		log.Fatal(err)
	}
}
