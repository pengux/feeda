package cmd

import (
	"encoding/xml"
	"fmt"
	"io"
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
		var ids, syncedAtIds []int64

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

		c := &http.Client{
			Timeout: 10 * time.Second,
		}
		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)
			syncedAtIds = append(syncedAtIds, feed.ID)

			go func(feed sqlite.Feed) {
				req, err := http.NewRequest(http.MethodGet, feed.URL, nil)
				if err != nil {
					log.Fatal(err)
				}

				req.Header.Set("User-Agent", userAgent)
				resp, err := c.Do(req)
				if err != nil {
					log.Fatalf("could not fetch URL %s: %v", feed.URL, err)
				}
				defer resp.Body.Close()

				var items []sqlite.Item
				if feed.Type == sqlite.FeedTypeRSS {
					items, err = createItemsFromRSS(resp.Body, feed)
					if err != nil {
						log.Fatal(err)
					}
				} else if feed.Type == sqlite.FeedTypeAtom {
					items, err = createItemsFromAtom(resp.Body, feed)
					if err != nil {
						log.Fatal(err)
					}
				}

				var inserted int64
				if len(items) > 0 {
					inserted, err = sqlite.CreateIgnoreItems(db, items...)
					if err != nil {
						log.Fatal(err)
					}

				}

				fmt.Printf("%d. %d items added\n", feed.ID, inserted)

				wg.Done()
			}(*feed)
		}

		wg.Wait()

		err = sqlite.SetFeedsSyncedAtNow(db, syncedAtIds...)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}

func createItemsFromRSS(body io.Reader, feed sqlite.Feed) ([]sqlite.Item, error) {
	var err error
	var content rss2
	var items []sqlite.Item

	decoded := xml.NewDecoder(body)
	err = decoded.Decode(&content)
	if err != nil {
		return items, fmt.Errorf("could not read data for URL %s: %s", feed.URL, err)
	}

	for _, item := range content.Items {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			// Try to parse as "Mon, 2 Jan 2006 15:04:05 -0700" (without leading zero on the day of month)
			pubDate, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
			if err != nil {
				return items, fmt.Errorf("could not parse PubDate for URL %s: %s", feed.URL, err)
			}

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

	return items, nil
}

func createItemsFromAtom(body io.Reader, feed sqlite.Feed) ([]sqlite.Item, error) {
	var err error
	var content atom
	var items []sqlite.Item

	decoded := xml.NewDecoder(body)
	err = decoded.Decode(&content)
	if err != nil {
		return items, fmt.Errorf("could not read data for URL %s: %s", feed.URL, err)
	}

	for _, item := range content.Items {
		pubDate, err := time.Parse(time.RFC3339, item.Updated)
		if err != nil {
			return items, fmt.Errorf("could not parse Updated for URL %s: %s", feed.URL, err)
		}

		if strings.TrimSpace(item.ID) == "" {
			item.ID = item.Link.Href
		}

		desc := strings.TrimSpace(item.Content)
		if desc == "" {
			desc = strings.TrimSpace(item.Summary)
		}

		items = append(items, sqlite.Item{
			FeedID:      feed.ID,
			GUID:        item.ID,
			URL:         item.Link.Href,
			Title:       item.Title,
			Desc:        desc,
			PublishedAt: pubDate,
		})
	}

	return items, nil
}
