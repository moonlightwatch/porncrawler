package siteanalysis

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/url"
	"porncrawler/data"

	"github.com/yanyiwu/gojieba"
)

func SetJieba() {
	gojieba.DICT_PATH = "/tmp/jieba.dict.utf8"
	buf := bytes.Buffer{}
	for _, line := range jiebawords {
		buf.WriteString(line)
	}
	ioutil.WriteFile(gojieba.DICT_PATH, buf.Bytes(), 0777)
}

func NewSiteAnalyseTool(d *data.DataInterface) *SiteAnalyseTool {
	s := &SiteAnalyseTool{}
	s.jieba = gojieba.NewJieba()
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
