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

func TestMakeChunksFromSSML(t *testing.T) {
	ssml := `<speak> <p>Welcome to this new episode of hn dot audio. <break time="500ms"></break> This episode was recorded on Tuesday  <say-as interpret-as="date" format="yyyymmdd" detail="1">2019-02-12</say-as> and today we talk about <emphasis level="strong">Why China is obsessed with numbers.</emphasis></p>
<audio src="https://www.thesoundarchive.com/ringtones/Woke-up-This-Morning-Chosen-One-Mix.mp3" clipEnd="9s">could not load mp3 file</audio>
<p>The Chinese fascination with numbers – and how much they are part of both online and offline lives – is a societal quirk that baffles long-term tourists and expats alike.</p>
<p>As a newly minted Beijinger, there were certain things my brain quickly scrambled to make room for: the exact time I needed to leave home in the mornings to avoid being squashed into human dumpling filling on the rush-hour subway ride; the location of the best spots for mala xiang guo (a stir-fried version of hot pot); to never flush toilet paper; and to never, ever attempt eating a soup dumpling by putting it straight into your mouth.</p> <p>One task, though, seemed impossible: remembering my QQ number, a string of randomly assigned digits that served as the user identification for the QQ messaging service our office – and many others in China – used.</p>
</speak>`
	chunks, err := MakeChunks(ssml, 540)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(chunks))
}

func TestNoSsml(t *testing.T) {
	_, err := MakeChunks("This is just plain text", 10)
	assert.NotNil(t, err)
	assert.Equal(t, "No <speak> or <p> elements found. Processing of plain text is not supported", err.Error())
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

func TestHtml(t *testing.T) {
	html := `<div><div class="rad-cover stacked-cover" style="background-color:#FFFFFF"><header class="rad-header header-black">
    <div class="rad-header-wrapper">
      <div class="rad-cover-kicker"><i>The</i> FUTURE <i>of</i> WORK</div>
      <div class="rad-headline-svg">
        <img src="https://static01.nyt.com/newsgraphics/2019/02/18/work/e56b77224653a7cddb016a0b7e4164f04f6a6865/opener-svgs/24mag-24happiness-t.svg" class="cover-svg cover-svg-nonmobile" alt="America’s Professional Elite: Wealthy, Successful and Miserable">
        
      </div>
      
      <p class="rad-summary" style="white-space: normal; float: none; display: block; position: static;">The upper echelon is hoarding money and privilege to a degree<br data-owner="balance-text">not seen in decades. But that doesn’t make them happy at work.</p>
      <p class="rad-byline-pubdate"><span class="rad-byline">By <a class="byline-author" href="https://www.nytimes.com/by/charles-duhigg">CHARLES DUHIGG</a></span></p>
      <div class="rad-second-byline">
        Illustration by TRACY MA
      </div>
      
      
    </div>
  </header>
  
</div>




<div class="rad-story-body">
  
    
    <p class="paragraph"><strong>My first, charmed</strong> week as a 
student at Harvard Business School, late in the summer of 2001, felt 
like a halcyon time for capitalism. AOL Time Warner, Yahoo and Napster 
were benevolently connecting the world. Enron and WorldCom were bringing
 innovation to hidebound industries. President George W. Bush — an 
H.B.S. graduate himself — had promised to deliver progress and 
prosperity with businesslike efficiency.</p>

    
    
    
    
    
    
    
    
  
    
    <p class="paragraph">The next few years would prove how little we 
(and Washington and much of corporate America) really understood about 
the economy and the world. But at the time, for the 895 first-years 
preparing ourselves for business moguldom, what really excited us was 
our good luck. A Harvard M.B.A. seemed like a winning lottery ticket, a 
gilded highway to world-changing influence, fantastic wealth and — if 
those self-satisfied portraits that lined the hallways were any 
indication — a lifetime of deeply meaningful work.</p>

    
    

    
    
    <div class="rad-interactive full_bleed" id="interactive-24mag-happiness-pq1" data-id="100000006370079" data-slug="24mag-happiness-pq1">
  <div class="rad-interactive-wrapper">
    
      <h2 class="interactive-headline"></h2>
      
    
    <p class="pullquote">‘When I die, is anyone going to care that I earned an extra percentage point of return? My work feels totally meaningless.’</p>

    <p class="credit"></p>
    <p class="notes"></p>
  </div>
</div
    
    <p class="paragraph">After our reunion, I wondered if my Harvard 
class — or even just my own friends there — were an anomaly. So I began 
looking for data about the nation’s professional psyche. What I found 
was that my classmates were hardly unique in their dissatisfaction; even
 in a boom economy, a surprising portion of Americans are professionally
 miserable right now. In the mid-1980s, roughly 61 percent of workers 
told pollsters they were satisfied with their jobs. Since then, that 
number has declined substantially, hovering around half; the low point 
was in 2010, when only 43 percent of workers were satisfied, <a href="https://www.conference-board.org/blog/postdetail.cfm?post=6391">according to data collected by the Conference Board</a>,
 a nonprofit research organization. The rest said they were unhappy, or 
at best neutral, about how they spent the bulk of their days. Even among
 professionals given to lofty self-images, like those in medicine and 
law, other studies have noted a rise in discontent. Why? Based on my own
 conversations with classmates and the research I began reviewing, the 
answer comes down to oppressive hours, political infighting, increased 
competition sparked by globalization, an “always-on culture” bred by the
 internet — but also something that’s hard for these professionals to 
put their finger on, an underlying sense that their work isn’t worth the
 grueling effort they’re putting into it.</p>

    
    
    
    
    
    
    
    
  
    
    <p class="paragraph">This wave of dissatisfaction is especially 
perverse because corporations now have access to decades of scientific 
research about how to make jobs better. “We have so much evidence about 
what people need,” says Adam Grant, a professor of management and 
psychology at the University of Pennsylvania (and a contributing opinion
 writer at The Times). Basic financial security, of course, is critical —
 as is a sense that your job won’t disappear unexpectedly. What’s 
interesting, however, is that once you can provide financially for 
yourself and your family, according to studies, additional salary and 
benefits don’t reliably contribute to worker satisfaction. Much more 
important are things like whether a job provides a sense of autonomy — 
the ability to control your time and the authority to act on your unique
 expertise. People want to work alongside others whom they respect (and,
 optimally, enjoy spending time with) and who seem to respect them in 
return.</p>

    
    
    
    
    
    
    
    
  
    
    <p class="paragraph">And finally, workers want to feel that their 
labors are meaningful. “You don’t have to be curing cancer,” says Barry 
Schwartz, a visiting professor of management at the University of 
California, Berkeley. We want to feel that we’re making the world 
better, even if it’s as small a matter as helping a shopper find the 
right product at the grocery store. “You can be a salesperson, or a toll
to figure out why particular janitors at a large hospital were so much 
more enthusiastic than others. So they began conducting interviews and 
found that, by design and habit, some members of the janitorial staff 
saw their jobs not as just tidying up but as a form of healing. One 
woman, for instance, mopped rooms inside a brain-injury unit where many 
residents were comatose. The woman’s duties were basic: change bedpans, 
pick up trash. But she also sometimes took the initiative to swap around
 the pictures on the walls, because she believed a subtle stimulation 
change in the unconscious patients’ environment might speed their 
recovery. She talked to other convalescents about their lives. “I enjoy 
entertaining the patients,” she told the researchers. “That is not 
really part of my job description, but I like putting on a show for 
them.” She would dance around, tell jokes to families sitting vigil at 
bedsides, try to cheer up or distract everyone from the pain and 
uncertainty that otherwise surrounded them. <a href="http://webuser.bus.umich.edu/janedut/High%20Quality%20Connections/Interpersonal%20Sensemaking.pdf">In a 2003 study</a>
 led by the researchers, another custodian described cleaning the same 
room two times in order to ease the mind of a stressed-out father.</p>

    
    
    
    
    
    
    
    
  
    
    <p class="paragraph">To some, the moral might seem obvious: If you 
see your job as healing the sick, rather than just swabbing up messes, 
you’re likely to have a deeper sense of purpose whenever you grab the 
mop. But what’s remarkable is how few workplaces seem to have 
internalized this simple lesson. “There are so many jobs where people 
feel like what they do is relatively meaningless,” Wrzesniewski says. 
“Even for well-paid positions, or jobs where you assume workers feel a 
sense of meaning, people feel like what they’re doing doesn’t matter.” 
That’s certainly true for my miserable classmate earning $1.2 million a 
year. Even though, in theory, the investments he makes each day help 
fund pensions — and thus the lives of retirees — it’s pretty hard to see
 that altruism from his window office in a Manhattan skyscraper. “It’s 
just numbers on a screen to me,” he told me. “I’ve never met a retiree 
who enjoyed a vacation because of what I do. It’s so theoretical it 
hardly seems real.”</p>

    
    
    
    
    
    
    
    
  
    
    <p class="paragraph"><strong>There is a</strong> raging debate — on 
newspaper pages, inside Silicon Valley, among presidential hopefuls — as
 to what constitutes a “good job.” I’m an investigative business 
reporter, and so I have a strange perspective on this question. When I 
speak to employees at a company, it’s usually because something has gone
 wrong. My stock-in-trade are sources who feel their employers are 
acting unethically or ignoring sound advice. The workers who speak to me
 are willing to describe both the good and the bad in the places where 
they work, in the hope that we will all benefit from their insights.</p>

    
    
    <div class="rad-interactive full_bleed" id="interactive-24mag-happiness-pq2" data-id="100000006370329" data-slug="24mag-happiness-pq2">
  <div class="rad-interactive-wrapper">
    
      <h2 class="interactive-headline"></h2>
      
    
    <p class="pullquote">The smoothest life paths sometimes fail to teach us about what really brings us satisfaction day to day.</p>

    <p class="credit"></p>
    <p class="notes"></p>
  </div>
</div>
    
  
</div></div>`

	_, _ = MakeChunks(html, 1800)
	t.Fatal("no")
}
