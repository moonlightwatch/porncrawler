package siteanalysis

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"porncrawler/data"

	"github.com/yanyiwu/gojieba"
)

func SetJieba() {
	os.MkdirAll("/tmp/dict", 0666)
	resp, _ := http.Get("https://raw.githubusercontent.com/yanyiwu/gojieba/master/dict/hmm_model.utf8")
	b, _ := ioutil.ReadAll(resp.Body)
	ioutil.WriteFile("/tmp/dict/hmm_model.utf8", b, 0666)

	resp, _ = http.Get("https://raw.githubusercontent.com/yanyiwu/gojieba/master/dict/idf.utf8")
	b, _ = ioutil.ReadAll(resp.Body)
	ioutil.WriteFile("/tmp/dict/idf.utf8", b, 0666)

	resp, _ = http.Get("https://raw.githubusercontent.com/yanyiwu/gojieba/master/dict/jieba.dict.utf8")
	b, _ = ioutil.ReadAll(resp.Body)
	ioutil.WriteFile("/tmp/dict/jieba.dict.utf8", b, 0666)

	resp, _ = http.Get("https://raw.githubusercontent.com/yanyiwu/gojieba/master/dict/stop_words.utf8")
	b, _ = ioutil.ReadAll(resp.Body)
	ioutil.WriteFile("/tmp/dict/stop_words.utf8", b, 0666)

	resp, _ = http.Get("https://raw.githubusercontent.com/yanyiwu/gojieba/master/dict/user.dict.utf8")
	b, _ = ioutil.ReadAll(resp.Body)
	ioutil.WriteFile("/tmp/dict/user.dict.utf8", b, 0666)
}

func NewSiteAnalyseTool(d *data.DataInterface) *SiteAnalyseTool {
	s := &SiteAnalyseTool{}
	s.jieba = gojieba.NewJieba("/tmp/dict/jieba.dict.utf8", "/tmp/dict/hmm_model.utf8", "/tmp/dict/user.dict.utf8", "/tmp/dict/idf.utf8", "/tmp/dict/stop_words.utf8")
	s.d = d
	return s
}

type SiteAnalyseTool struct {
	jieba *gojieba.Jieba
	d     *data.DataInterface
}

func (s SiteAnalyseTool) CheckSite(site data.SiteData) bool {
	if site.Title == "" || len(site.Links) == 0 {
		return false
	}
	words := s.jieba.CutForSearch(site.Text, true)
	for _, w := range words {
		for _, sw := range swords {
			if w == sw {
				s.d.AddSite(site)
				passedHost := map[string]bool{}
				for _, u := range site.Links {
					p, err := url.Parse(u)
					if err != nil {
						continue
					}
					if _, ok := passedHost[p.Host]; !ok {
						s.d.AddTarget(p.Host, u)
					}
				}
				log.Printf("%s add %d\n", site.Title, len(site.Links))
				return true
			}
		}
	}
	return false
}
