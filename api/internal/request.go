package internal

import (
	"log"
	"net/http"
	"net/url"
)

const (
	LoginPageURL          = "https://cgifederal.secure.force.com/"
	AppointmentHistoryURL = "https://cgifederal.secure.force.com/appointmenthistory"
	SchedulePageURL       = "https://cgifederal.secure.force.com/scheduleappointment"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(cred Credentials) []*http.Cookie {

	data := url.Values{}
	data.Set("email", cred.Email)
	data.Set("password", cred.Password)

	res, err := http.PostForm(LoginPageURL, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.StatusCode)
	body, _ := ReadResBody(res)
	err = HTMLToFile("response.html", string(body))
	if err != nil {
		log.Println(err)
		return nil
	}
	return res.Cookies()
}

func GetAppointmentHistory(cookies []*http.Cookie) (string, error) {
	req, err := http.NewRequest("GET", SchedulePageURL, nil)
	if err != nil {
		return "", err
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return doGetRequestBody(req)
}

func doGetRequestBody(req *http.Request) (string, error) {
	// Do the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	// Read the body
	html, err := ReadResBody(res)
	if err != nil {
		return "", err
	}

	log.Println(res.StatusCode)
	return string(html), nil
}
