package cmd

import (
	"encoding/xml"
)

type (
	rss struct {
		XMLName xml.Name `xml:"rss"`
		Channel *channel `xml:"channel"`
	}

	channel struct {
		Items []item `xml:"channel>item"`
	}

	item struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		GUID        string `xml:"guid"`
		PubDate     string `xml:"pubDate"`
	}
)
