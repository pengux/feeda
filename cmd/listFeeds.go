package cmd

import (
	"fmt"
	"log"

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
			fmt.Printf("%d. %s\n", feed.ID, feed.URL)
		}
	},
}

func init() {
	RootCmd.AddCommand(listFeedsCmd)
}
