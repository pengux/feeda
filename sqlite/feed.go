package sqlite

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	feedsTable = "feeds"
)

// Types for feeds
const (
	FeedTypeRSS  feedType = "RSS"
	FeedTypeAtom feedType = "Atom"
)

type (
	feedType string

	// Feed contains the URL to the RSS/Atom feed
	Feed struct {
		ID        int64
		URL       string
		Type      feedType
		CreatedAt time.Time
		SyncedAt  *time.Time
	}

	// FeedFilter is used to filter feeds in lists
	FeedFilter struct {
		FeedID     int64
		ReadStatus itemReadStatus
		Limit      int64
		Offset     int64
	}
)

// CreateIgnoreFeeds persists feeds to DB and if it already exists then skip it
func CreateIgnoreFeeds(db cruderExecer, feeds ...Feed) error {
	var values []string
	var params []interface{}

	if len(feeds) == 0 {
		return errors.New("missing feeds to create")
	}

	for _, feed := range feeds {
		values = append(values, "(?, ?)")
		params = append(params, feed.URL, string(feed.Type))
	}

	_, err := db.Exec(
		fmt.Sprintf(`INSERT OR IGNORE INTO "%s" (url, type) VALUES %s`, feedsTable, strings.Join(values, ",")),
		params...,
	)

	return err
}

// ListFeeds returns a list of feeds from DB
func ListFeeds(db cruderQueryer, ids ...int64) ([]*Feed, error) {
	var feeds []*Feed
	var wheres []string
	var whereSQL string
	var params []interface{}

	for _, id := range ids {
		wheres = append(wheres, "?")
		params = append(params, id)
	}

	if len(wheres) > 0 {
		whereSQL = fmt.Sprintf(" WHERE id IN (%s)", strings.Join(wheres, ","))
	}

	rows, err := db.Query(
		fmt.Sprintf(`SELECT id, url, type, created_at, synced_at FROM "%s"%s ORDER BY id`, feedsTable, whereSQL),
		params...,
	)
	if err != nil {
		return feeds, err
	}
	defer rows.Close()
	for rows.Next() {
		f := &Feed{}
		err = rows.Scan(&f.ID, &f.URL, &f.Type, &f.CreatedAt, &f.SyncedAt)
		if err != nil {
			return feeds, err
		}

		feeds = append(feeds, f)
	}

	return feeds, nil
}

// SetFeedsSyncedAtNow sets the synced_at column of feeds to CURRENT_TIMESTAMP
func SetFeedsSyncedAtNow(db cruderExecer, ids ...int64) error {
	var placeholders []string
	var params []interface{}

	for _, id := range ids {
		placeholders = append(placeholders, "?")
		params = append(params, id)
	}

	if len(placeholders) == 0 {
		return errors.New("missing ids to update")
	}

	_, err := db.Exec(
		fmt.Sprintf(`UPDATE "%s" SET synced_at = CURRENT_TIMESTAMP WHERE id IN (%s)`, feedsTable, strings.Join(placeholders, ",")),
		params...,
	)

	return err
}

// DeleteFeeds removes one or more feeds from DB
func DeleteFeeds(db cruderExecer, ids ...int64) error {
	var placeholders []string
	var params []interface{}

	for _, id := range ids {
		placeholders = append(placeholders, "?")
		params = append(params, id)
	}

	if len(placeholders) == 0 {
		return errors.New("missing ids to update")
	}

	_, err := db.Exec(
		fmt.Sprintf(`DELETE FROM "%s" WHERE id IN (%s)`, feedsTable, strings.Join(placeholders, ",")),
		params...,
	)

	return err
}
