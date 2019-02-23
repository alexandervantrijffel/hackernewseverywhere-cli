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
	errorcheck.CheckLogFatalf(err, "goquery failed to parse text. %s", ssml)

	speak := doc.Find("speak")
	if len(speak.Nodes) > 0 {
		return processSsml(speak, maxChunkChars)
	}

	paragraphs := doc.Find("p")
	logging.Infof("Found %d paragraps", len(paragraphs.Nodes))
	if len(paragraphs.Nodes) > 0 {
		return processParagraphs(paragraphs, maxChunkChars)
	}
	return nil, errorcheck.LogAndWrapAsError("No <speak> or <p> elements found. Processing of plain text is not supported")
}
func processParagraphs(paragraphs *goquery.Selection, maxChunkChars int) ([]string, error) {
	logging.Info("Processing paragraphs")
	var chunks []string
	var chunkHtml string
	var lastErr error
	paragraphs.Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if len(strings.TrimSpace(text)) == 0 {
			ohtml, _ := goquery.OuterHtml(s)
			logging.Debugf("Skipping paragraph without text. %s", ohtml)
			return
		}
		html := fmt.Sprintf("<p>%s</p>%s", addBreaks(s.Text()), br(800))
		if len(html) > maxChunkChars {
			lastErr = errorcheck.LogAndWrapAsError("A single paragraph has more chars than the maximum per chunk of %d chars. Cannot process: %s",
				maxChunkChars, html)
			return
		}
		if len(chunkHtml)+len(html) > maxChunkChars {
			chunks = append(chunks, addSpeak(chunkHtml))
			chunkHtml = ""
		}
		chunkHtml += html
	})
	chunks = append(chunks, addSpeak(chunkHtml))
	if lastErr != nil {
		return nil, lastErr
	}
	logging.Infof("Have %d chunks", len(chunks))
	return chunks, nil
}

func processSsml(speak *goquery.Selection, maxChunkChars int) ([]string, error) {
	logging.Info("Processing SSML text")
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
				if len(html) > maxChunkChars {
					lastErr = errorcheck.LogAndWrapAsError("A single html element has more chars than the maximum per chunk of %d chars. Cannot process: %s",
						maxChunkChars, html)
					return
				}
				if len(chunkHtml)+len(html) > maxChunkChars {
					chunks = append(chunks, addSpeak(chunkHtml))
					chunkHtml = ""
				}
				chunkHtml += html
			}
		})
		chunks = append(chunks, addSpeak(chunkHtml))
		if lastErr != nil {
			return nil, lastErr
		}
		logging.Infof("Have %d chunks", len(chunks))
		return chunks, nil
	}
	err := errors.New("No html children found in ssml")
	return nil, err
}

func br(ms int) string {
	return fmt.Sprintf(`<break time="%dms"></break>`, ms)
}

func addBreaks(text string) string {
	a := strings.Replace(text, ",", ","+br(200), -1)
	return strings.Replace(a, ";", ","+br(200), -1)
}

func addSpeak(text string) string {
	return fmt.Sprintf("<speak>%s</speak>", text)
}
