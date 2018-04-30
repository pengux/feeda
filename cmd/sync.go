package cmd

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pengux/feeda/sqlite"
	"github.com/spf13/cobra"
)

// syncCmd fetches one or multiple feeds and persists their items
// to DB
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Download latest items of one or multiple feeds",
	Long: `Fetch latest content from one or multiple feeds, parse it and
store the items to DB. Calling this command without any arguments will
sync all feeds. If feed IDs are provided as arguments, only those feeds will
be synced. Example:

# Sync only feeds with ID = 1 and ID = 3
sync 1 3

# Sync all feeds
sync`,
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int64

		for _, arg := range args {
			id, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			ids = append(ids, id)
		}

		feeds, err := sqlite.ListFeeds(db, ids...)
		if err != nil {
			log.Fatal(err)
		}

		var c = &http.Client{
			Timeout: 10 * time.Second,
		}
		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)

			go func(feed sqlite.Feed) {
				resp, err := c.Get(feed.URL)
				if err != nil {
					log.Fatalf("could not fetch URL %s: %v", feed.URL, err)
				}
				defer resp.Body.Close()

				var rss channel
				decoded := xml.NewDecoder(resp.Body)
				err = decoded.Decode(&rss)
				if err != nil {
					log.Fatalf("could not read data for URL %s: %v", feed.URL, err)
				}

				var items []sqlite.Item
				for _, item := range rss.Items {
					pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
					if err != nil {
						log.Fatalf("could not parse PubDate for URL %s: %v", feed.URL, err)
					}

					if strings.TrimSpace(item.GUID) == "" {
						item.GUID = item.Link
					}

					items = append(items, sqlite.Item{
						FeedID:      feed.ID,
						GUID:        item.GUID,
						URL:         item.Link,
						Title:       item.Title,
						Desc:        item.Description,
						PublishedAt: pubDate,
					})
				}

				inserted, err := sqlite.CreateIgnoreItems(db, items...)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("%d. %d items added\n", feed.ID, inserted)

				wg.Done()
			}(*feed)
		}

		wg.Wait()

		err = sqlite.SetFeedsSyncedAtNow(db, ids...)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
