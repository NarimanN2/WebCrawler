package crawler

import (
	"bytes"
	"golang.org/x/net/html"
)

type HTMLParser struct {
	document []byte
}

func NewHTMLParser(document []byte) HTMLParser {
	return HTMLParser{
		document: document,
	}
}

func (h HTMLParser) FindUrls() []string {
	var tags []string
	tokenizer := html.NewTokenizer(bytes.NewReader(h.document))

	for tokenType := tokenizer.Next(); tokenType != html.ErrorToken; tokenType = tokenizer.Next() {
		tokenType := tokenizer.Next()

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()

			if token.Data == "a" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						tags = append(tags, attribute.Val)
					}
				}
			}
		}
	}

	return tags
}
