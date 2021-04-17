package downloader

import (
	"context"
	"log"
	"porncrawler/data"
	"porncrawler/siteanalysis"
	"time"

	"github.com/chromedp/chromedp"
)

func NewBrowser(d *data.DataInterface, s *siteanalysis.SiteAnalyseTool) *Browser {
	b := &Browser{}
	var ctx context.Context
	ctx, b.cancel1 = chromedp.NewContext(context.Background())
	b.baseContext, b.cancel2 = context.WithTimeout(ctx, 30*time.Second)
	b.d = d
	b.s = s
	b.Stopped = make(chan bool)
	return b
}

type Browser struct {
	baseContext context.Context
	cancel1     context.CancelFunc
	cancel2     context.CancelFunc
	d           *data.DataInterface
	s           *siteanalysis.SiteAnalyseTool
	running     bool
	Stopped     chan bool
}

func (b *Browser) Close() {
	b.running = false
	b.cancel2()
	b.cancel1()

}

func (b *Browser) request(url string) data.SiteData {
	title := ""
	urls := []string{}
	text := ""
	currentURL := ""
	err := chromedp.Run(b.baseContext)
	if err != nil {
		log.Println(err)
	}
	err = chromedp.Run(b.baseContext,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
	if err != nil {
		log.Println(err)
	}
	err = chromedp.Run(b.baseContext,
		chromedp.Evaluate(`document.URL;`, &currentURL),
		chromedp.Title(&title),
		// chromedp.Evaluate(`document.title;`, &title),
		chromedp.Evaluate(`document.body.innerText;`, &text),
		chromedp.Evaluate(`var l=new Array();for(var i=0;i<document.links.length;i++){l.push(document.links[i].href);};l`, &urls),
	)
	if err != nil {
		log.Println(err)
	}

	return data.SiteData{
		URL:   currentURL,
		Title: title,
		Links: urls,
		Text:  text,
	}
}

func (b *Browser) RequestLoop() {
	b.running = true

	for b.running {
		t := b.d.GetTarget()
		if t == "" {
			time.Sleep(time.Second)
			continue
		}
		log.Printf("request: %s\n", t)
		site := b.request(t)
		b.s.CheckSite(site)

	}
	b.Stopped <- true
}
