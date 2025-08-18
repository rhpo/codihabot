package markdown

import (
	"strings"

	"github.com/russross/blackfriday"
)

func MarkdownToHTML(markdownContent string) string {
	markdownContent = strings.ReplaceAll(markdownContent, "<br>", "<br>\n\n")

	html := blackfriday.MarkdownCommon([]byte(markdownContent))
	// Remove surrounding <p> tags if present to make it inline
	result := string(html)
	if len(result) >= 7 && result[:3] == "<p>" && result[len(result)-5:] == "</p>\n" {
		result = result[3 : len(result)-5]
	}
	return result
}
