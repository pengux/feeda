package cmd

import "encoding/xml"

type (
	rss2 struct {
		XMLName xml.Name   `xml:"rss"`
		Items   []rss2Item `xml:"channel>item"`
	}

	rss2Item struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		GUID        string `xml:"guid"`
		PubDate     string `xml:"pubDate"`
	}

	atom struct {
		XMLName xml.Name   `xml:"feed"`
		Items   []atomItem `xml:"entry"`
	}

	atomItem struct {
		Title   string `xml:"title"`
		Link    string `xml:"link"`
		ID      string `xml:"id"`
		Content string `xml:"content"`
		Summary string `xml:"summary"`
		Updated string `xml:"updated"`
	}
)
