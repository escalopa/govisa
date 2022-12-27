package internal

import (
	"io"
	"net/http"
	"os"
)

const (
	HTMLDir = "./html/"
)

func HTMLToFile(fileName string, html string) error {
	return os.WriteFile((HTMLDir + fileName), []byte(html), 0644)
}

func ReadResBody(res *http.Response) (string, error) {
	html, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(html), nil
}
