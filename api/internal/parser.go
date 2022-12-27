package internal

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// TODO: Read the table from the html page
func ParseTable(text string) string {
	var buf strings.Builder

	buf.WriteString("```\n")
	tkn := html.NewTokenizer(strings.NewReader(text))

	var isTd bool
	var n int

	for {

		tt := tkn.Next()

		switch tt {

		case html.ErrorToken:
			return buf.String()

		case html.StartTagToken:
			t := tkn.Token()
			isTd = t.Data == "td"

		case html.TextToken:
			t := tkn.Token()
			if isTd {
				fmt.Printf("%s ", t.Data)
				n++
			}

			if isTd && n%3 == 0 {
				buf.WriteString("\n")
			}

			isTd = false
		}
	}
}

// TODO: From table text to markdown table for telegram bot
func MarkdownTablify(table string) string {
	return ""
}
