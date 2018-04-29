package sqlite

import "fmt"

// EnsureTables will creates the DB tables if not already exists
func EnsureTables(db cruderExecer) error {
	_, err := db.Exec(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"url" TEXT NOT NULL UNIQUE,
			"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"synced_at" TIMESTAMP
		);`, feedsTable),
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"feed_id" INTEGER NOT NULL,
			"guid" TEXT UNIQUE,
			"url" TEXT,
			"title" TEXT,
			"desc" TEXT,
			"published_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"read_at" TIMESTAMP,
			FOREIGN KEY("feed_id") REFERENCES "%s"("id") ON DELETE CASCADE
		);`, itemsTable, feedsTable),
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)

	return err
}
