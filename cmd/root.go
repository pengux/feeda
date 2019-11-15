package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	// SQLite3 driver
	"feeda/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

const (
	userAgent = "Feeda_feed_aggregator/1.0"
)

var (
	dbPath = "~/.feeda/db.sqlite"
	db     *sql.DB
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "feeda",
	Short: "Command line tool to manage feeds (RSS)",
	Long: `Feeds are stored in a SQLite database which default to ~/.feeda/db.sqlite
Usual work flow is:

# Add a feed
feeda add [URL to feed]

# List feeds
feeda listFeeds

# Sync feeds
feeda sync

# Sync feed with ID=1
feeda sync 1

# List 10 unread entries from feed with ID=1 and set them as read order by oldest first
feeda list --unread --setAsRead --limit=10 --feed=1

# List 50 unread entries from all feeds and show only their URLs and set them as read
# order by oldest first. Pipe it to "open" to open in the default browser
feeda list -l=50 -u -r -o | xargs open
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initDB)

	RootCmd.Flags().StringVar(&dbPath, "db", "", "Location of DB, defaults to ~/.feeda/db.sqlite")
}

// initDB initializes the SQLite DB
func initDB() {
	var err error

	if dbPath == "" {
		dbPath, err = homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		dbPath = path.Join(dbPath, ".feeda")
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			os.Mkdir(dbPath, os.ModePerm)
		}

		dbPath = path.Join(dbPath, "db.sqlite")
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	err = sqlite.EnsureTables(db)
	if err != nil {
		log.Fatal(err)
	}

}
