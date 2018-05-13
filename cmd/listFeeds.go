package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/pengux/feeda/sqlite"
	"github.com/spf13/cobra"
)

// listFeedsCmd represents the listFeeds command
var listFeedsCmd = &cobra.Command{
	Use:   "listFeeds",
	Short: "List all feeds",
	Long:  `List all feeds that has been added`,
	Run: func(cmd *cobra.Command, args []string) {
		feeds, err := sqlite.ListFeeds(db)
		if err != nil {
			log.Fatal(err)
		}

		for _, feed := range feeds {
			var attrs []string

			if feed.SyncedAt != nil {
				attrs = append(attrs, fmt.Sprintf("Synced: %s", feed.SyncedAt.Format("2006-01-02 15:04:05")))
			}

			total, err := sqlite.CountTotalByFeed(db, feed.ID)
			if err != nil {
				log.Fatal(err)
			}
			attrs = append(attrs, fmt.Sprintf("Total: %d", total))

			unread, err := sqlite.CountUnreadByFeed(db, feed.ID)
			if err != nil {
				log.Fatal(err)
			}
			attrs = append(attrs, fmt.Sprintf("Unread: %d", unread))

			fmt.Printf("%d. %s (%s)\n", feed.ID, feed.URL, strings.Join(attrs, ", "))
		}
	},
}

func init() {
	RootCmd.AddCommand(listFeedsCmd)
}
