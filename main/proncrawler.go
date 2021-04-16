package main

import (
	"log"
	"os"
	"os/signal"
	"porncrawler/data"
	"porncrawler/downloader"
	"porncrawler/siteanalysis"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	log.Printf("%s 'addr' 'password'\n", os.Args[0])
	log.Println("启动")
	d := data.NewDataInterface(&redis.Options{Addr: os.Args[1], Password: os.Args[2]})
	s := siteanalysis.NewSiteAnalyseTool(d)
	browserList := []*downloader.Browser{}
	browserList = append(browserList, downloader.NewBrowser(d, s))
	browserList = append(browserList, downloader.NewBrowser(d, s))
	browserList = append(browserList, downloader.NewBrowser(d, s))
	for _, b := range browserList {
		go b.RequestLoop()
		time.Sleep(time.Second)
		log.Println("start.")
	}
	quit := make(chan os.Signal, 5)
	signal.Notify(quit, os.Interrupt)
	<-quit
	for _, b := range browserList {
		b.Close()
	}
	for _, b := range browserList {
		<-b.Stopped
		log.Println("stoped.")
	}
	d.Close()
}
