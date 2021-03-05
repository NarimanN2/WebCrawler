package reader

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
)

type Reader interface {
	ReadLines(path string) []string
}

type FileReader struct {
}

func NewReader() Reader {
	return FileReader{
	}
}

func (fr FileReader) ReadLines(path string) []string {
	var urls []string
	file, err := os.Open(path)
	if err != nil {
		log.Info("Directory is empty")
		return urls
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	return urls
}
