package crawler

import (
	"context"
	"fmt"
	"github.com/NarimanN2/WebCrawler/internal/reader"
	"github.com/NarimanN2/WebCrawler/internal/writer"
	"github.com/hashicorp/go-uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

type Crawler struct {
	baseURL     string
	directory   string
	seedUrls    chan string
	visitedUrls map[string]bool
	waitGroup   sync.WaitGroup
	writer      writer.Writer
	reader      reader.Reader
}

func NewCrawler(url string, directory string, writer writer.Writer, reader reader.Reader) *Crawler {
	err := writer.Mkdir(directory)
	if err != nil {
		log.Error("Failed to create directory")
		panic(err)
	}
	return &Crawler{
		baseURL:     url,
		directory:   directory,
		seedUrls:    make(chan string),
		visitedUrls: make(map[string]bool),
		waitGroup:   sync.WaitGroup{},
		writer:      writer,
		reader:      reader,
	}
}

func (c *Crawler) Run(ctx context.Context) {
	previousSeedUrls := c.reader.ReadLines(path.Join(c.directory, "seeds.txt"))

	if len(previousSeedUrls) > 0 {
		for _, seedUrl := range previousSeedUrls {
			c.waitGroup.Add(1)
			go c.crawl(seedUrl)
		}
	} else {
		c.waitGroup.Add(1)
		go c.crawl(c.baseURL)
	}

	go func() {
		done := false

		for {
			select {
			case seedUrl := <-c.seedUrls:
				_, visited := c.visitedUrls[seedUrl]
				if !visited {
					if done {
						_ = c.writer.Write(path.Join(c.directory, "seeds.txt"), []byte(seedUrl+"\n"), true) // Save for later use
						continue
					}

					c.waitGroup.Add(1)
					go c.crawl(seedUrl)
					c.visitedUrls[seedUrl] = true
				}
			case <-ctx.Done():
				done = true
			}
		}
	}()

	c.waitGroup.Wait()
	close(c.seedUrls)
}

func (c *Crawler) fetchURL(url string) ([]byte, error) {
	maxRetries := 3
	var err error
	var response *http.Response

	for i := 0; i < maxRetries; i++ {
		response, err = http.Get(url)
		if err == nil {
			defer response.Body.Close()
			bytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return nil, err
			}

			return bytes, nil
		}

		time.Sleep(3 * time.Second)
	}

	return nil, err
}

func (c *Crawler) crawl(seedUrl string) {
	log.Info(fmt.Sprintf("Crawling %s...", seedUrl))
	defer c.waitGroup.Done()
	seed, err := url.Parse(seedUrl)
	if err != nil {
		log.Warn(fmt.Sprintf("Invalid url: %s", seedUrl))
		return
	}

	document, err := c.fetchURL(seed.String())
	if err != nil {
		log.Error(fmt.Sprintf("Failed to fetch url: %s %v", seedUrl, err))
		return
	}

	c.save(document)
	htmlParser := NewHTMLParser(document)
	urls := htmlParser.FindUrls()

	for _, newUrl := range urls {
		u, err := url.Parse(newUrl)
		if err != nil {
			log.Warn(fmt.Sprintf("Invalid url: %s", newUrl))
		} else {
			if !u.IsAbs() {
				newUrl = seed.ResolveReference(u).String()
			}

			if strings.HasPrefix(newUrl, c.baseURL) {
				c.seedUrls <- newUrl
			}
		}
	}
}

func (c *Crawler) save(document []byte) {
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		log.Error("failed to generate uuid")
		return
	}

	fileName := uuid + ".html"
	err = c.writer.Write(path.Join(c.directory, fileName), document, false)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to write to file"))
	}
}
