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
	browserList = append(browserList, downloader.NewBrowser("1", d, s))
	browserList = append(browserList, downloader.NewBrowser("2", d, s))
	browserList = append(browserList, downloader.NewBrowser("3", d, s))
	for _, b := range browserList {
		go b.RequestLoop()
		time.Sleep(time.Second)
		log.Printf("(%s) start.\n", b.Name)
	}
	quit := make(chan os.Signal, 5)
	signal.Notify(quit, os.Interrupt)
	<-quit
	for _, b := range browserList {
		go b.Close()
	}
	for _, b := range browserList {
		<-b.Stopped
		log.Printf("(%s) stoped.\n", b.Name)
	}
	d.Close()
}
