package cmd

import (
	"fmt"
	"log"

	"feeda/sqlite"

	"github.com/spf13/cobra"
)

var (
	unread, setAsRead, onlyURL *bool
	limit                      *int64
	feedID                     *int64
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List items from feeds",
	Long:  `List items synced from one or more feeds`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		filter := sqlite.ItemFilter{}

		if *unread {
			filter.ReadStatus = sqlite.ItemUnread
		}

		if *limit > 0 {
			filter.Limit = *limit
		}

		if *feedID > 0 {
			filter.FeedID = *feedID
		}

		items, err := sqlite.ListItems(db, filter)
		if err != nil {
			log.Fatal(err)
		}

		var ids []int64
		for _, item := range items {
			if *onlyURL {
				fmt.Printf("%s\n", item.URL)
			} else {
				fmt.Printf("%d. %s\n", item.ID, item.Title)
				fmt.Println(item.URL)
				fmt.Printf("Published: %s\n", item.PublishedAt.Format("2006-01-02 15:04:05"))
				if item.ReadAt != nil {
					fmt.Printf("Read: %s\n", item.ReadAt.Format("2006-01-02 15:04:05"))
				} else {
					fmt.Println("Unread")
				}
				fmt.Println(item.Desc)
				fmt.Println("")
			}

			ids = append(ids, item.ID)
		}

		if *setAsRead {
			err = sqlite.SetItemsAsReadNow(db, ids...)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	unread = listCmd.Flags().BoolP("unread", "u", false, "List only unread items")
	setAsRead = listCmd.Flags().BoolP("setAsRead", "r", false, "Set the listed items as read")
	limit = listCmd.Flags().Int64P("limit", "l", 10, "Limit number of items to listed")
	feedID = listCmd.Flags().Int64P("feed", "f", 0, "Feed ID of items to be listed")
	onlyURL = listCmd.Flags().BoolP("onlyURL", "o", false, "List only the item's URL")
}
