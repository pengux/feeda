package cmd

import (
	"log"
	"strconv"

	"github.com/pengux/feeda/sqlite"
	"github.com/spf13/cobra"
)

// deleteFeedCmd deletes one or more feeds and their items from DB
var deleteFeedCmd = &cobra.Command{
	Use:   "deleteFeed",
	Short: "Delete feeds",
	Long:  `Deletes one or more feeds from DB. All there items will also be deleted`,
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int64

		for _, arg := range args {
			id, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			ids = append(ids, id)
		}

		err := sqlite.DeleteFeeds(db, ids...)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deleteFeedCmd)
}
