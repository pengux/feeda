package cmd

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/pengux/feeda/sqlite"
	"github.com/spf13/cobra"
)

// addCmd adds one or multiple URLs of RSS feeds to the DB
var addCmd = &cobra.Command{
	Use:   "add [URL of feed] [URL of feed 2]...",
	Short: "Add RSS feeds",
	Long:  `Adds multiple RSS feeds to aggregate.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		c := &http.Client{
			Timeout: 10 * time.Second,
		}
		var wg sync.WaitGroup
		var feeds []sqlite.Feed

		for _, arg := range args {
			_, err = url.Parse(arg)
			if err != nil {
				log.Fatalf("could not parse URL %s: %s", arg, err)
			}

			wg.Add(1)
			go func(url string) {
				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					log.Fatal(err)
				}

				req.Header.Set("User-Agent", userAgent)
				resp, err := c.Do(req)
				if err != nil {
					log.Fatalf("could not fetch URL %s: %v", url, err)
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}

				feed := sqlite.Feed{
					URL:  url,
					Type: sqlite.FeedTypeRSS,
				}

				if isAtom(body) {
					feed.Type = sqlite.FeedTypeAtom
				}

				feeds = append(feeds, feed)
			}(arg)
		}

		err = sqlite.CreateIgnoreFeeds(db, feeds...)
		if err != nil {
			log.Fatal("could not add feeds:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}

func isAtom(b []byte) bool {
	return bytes.Contains(b, []byte("<feed"))
}
