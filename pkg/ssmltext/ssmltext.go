package ssmltext

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexandervantrijffel/goutil/errorcheck"
	"github.com/alexandervantrijffel/goutil/logging"
)

func MakeChunks(ssml string, maxChunkChars int) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(ssml))
	errorcheck.CheckLogFatalf(err, "goquery failed to parse ssml. %s", ssml)
	speak := doc.Find("speak")
	if len(speak.Nodes) == 0 {
		return nil, errors.New("Invalid ssml; No <speak> element found.")
	}
	children := speak.Children()
	var chunks []string
	if len(children.Nodes) > 0 {
		var chunkHtml string
		var lastErr error
		children.Each(func(i int, s *goquery.Selection) {
			if html, err := goquery.OuterHtml(s); err != nil {
				lastErr = errorcheck.CheckLogf(err, "Failed to retrieve html of %s", s.Text())
			} else {
				if strings.ToLower(goquery.NodeName(s)) == "p" {
					html = fmt.Sprintf("<p>%s</p>%s", addBreaks(s.Text()), br(800))
					logging.Info("HTML", html)
				}
				if len(chunkHtml)+len(html) > maxChunkChars {
					logging.Infof("Adding chunk %s", chunkHtml)
					chunks = append(chunks, chunkHtml)
					chunkHtml = ""
				}
				if len(html) > maxChunkChars {
					lastErr = errorcheck.LogAndWrapAsError("A single html element has more chars than the maximum per chunk of %d chars. Cannot process: %s",
						maxChunkChars, html)
					return
				}
				chunkHtml += html
			}
		})
		chunks = append(chunks, chunkHtml)
		if lastErr != nil {
			return nil, lastErr
		}
	} else {
		err := errors.New("No html children found in ssml")
		return nil, err
	}
	return chunks, nil
}

func br(ms int) string {
	return fmt.Sprintf(`<break time="%dms"></break>`, ms)
}

func addBreaks(text string) string {
	return strings.Replace(text, ",", ","+br(400), -1)
}
