package ssmltext

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexandervantrijffel/goutil/logging"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitWith("hackernewseverywhere-cli unit tests", false)
}

func TestMakeChunks(t *testing.T) {
	ssml := `<speak> <p>Welcome to this new episode of hn dot audio. <break time="500ms"></break> This episode was recorded on Tuesday  <say-as interpret-as="date" format="yyyymmdd" detail="1">2019-02-12</say-as> and today we talk about <emphasis level="strong">Why China is obsessed with numbers.</emphasis></p>
<audio src="https://www.thesoundarchive.com/ringtones/Woke-up-This-Morning-Chosen-One-Mix.mp3" clipEnd="9s">could not load mp3 file</audio>
<p>The Chinese fascination with numbers – and how much they are part of both online and offline lives – is a societal quirk that baffles long-term tourists and expats alike.</p>
<p>As a newly minted Beijinger, there were certain things my brain quickly scrambled to make room for: the exact time I needed to leave home in the mornings to avoid being squashed into human dumpling filling on the rush-hour subway ride; the location of the best spots for mala xiang guo (a stir-fried version of hot pot); to never flush toilet paper; and to never, ever attempt eating a soup dumpling by putting it straight into your mouth.</p> <p>One task, though, seemed impossible: remembering my QQ number, a string of randomly assigned digits that served as the user identification for the QQ messaging service our office – and many others in China – used.</p>
</speak>`
	chunks, err := MakeChunks(ssml, 560)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(chunks))
	t.Fatal("no")
}

func TestNoSsml(t *testing.T) {
	_, err := MakeChunks("This is just plain text", 10)
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid ssml; No <speak> element found.", err.Error())
}

func TestGetFirstElementHtml(t *testing.T) {
	test := `<speak><p>My paragraph</p></speak>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(test))
	var childrenHtml []string
	doc.Find("speak").Children().Each(func(i int, s *goquery.Selection) {
		html, _ := goquery.OuterHtml(s)
		childrenHtml = append(childrenHtml, html)
	})
	if childrenHtml[0] != "<p>My paragraph</p>" {
		t.Fatalf("First element html is not valid: '%s'", childrenHtml[0])
	}
}

func TestGetTwoElements(t *testing.T) {
	test := `<speak><p>My paragraph</p><break time="800ms"/></speak>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(test))
	var childrenHtml []string
	doc.Find("speak").Children().Each(func(i int, s *goquery.Selection) {
		html, _ := goquery.OuterHtml(s)
		childrenHtml = append(childrenHtml, html)
	})
	if len(childrenHtml) != 2 {
		t.Fatalf("Found %d elements instead of 2", len(childrenHtml))
	}
}
