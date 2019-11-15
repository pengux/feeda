package cmd

import (
	"log"
	"strconv"

	"feeda/sqlite"

	"github.com/spf13/cobra"
)

// deleteCmd deletes one or more items from DB
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete items",
	Long:  `Deletes one or more items from DB`,
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int64

		for _, arg := range args {
			id, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			ids = append(ids, id)
		}

		err := sqlite.DeleteItems(db, ids...)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
