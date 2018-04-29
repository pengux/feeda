package cmd

import (
	"log"
	"net/url"

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
		for _, arg := range args {
			_, err = url.Parse(arg)
			if err != nil {
				log.Fatalf("could not parse URL %s: %s", arg, err)
			}

		}

		err = sqlite.CreateIgnoreFeeds(db, args...)
		if err != nil {
			log.Fatal("could not add feeds:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
