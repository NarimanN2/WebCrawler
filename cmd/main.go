package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/NarimanN2/WebCrawler/internal/crawler"
	"github.com/NarimanN2/WebCrawler/internal/reader"
	"github.com/NarimanN2/WebCrawler/internal/writer"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	url := flag.String("u", "https://bigplato.tilda.ws", "Website address")
	dir := flag.String("d", "./pages", "Directory")
	flag.Parse()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		sig := <- signals
		log.Info(fmt.Sprintf("Recieved %s signal", sig))
		cancelFunc()
	}()

	fileWriter := writer.NewWriter()
	fileReader := reader.NewReader()
	webCrawler := crawler.NewCrawler(*url, *dir, fileWriter, fileReader)
	webCrawler.Run(ctx)
}
