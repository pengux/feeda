package sqlite

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	itemsTable = "items"
)

// Statuses for whether an item is read or unread
const (
	ItemRead itemReadStatus = iota + 1
	ItemUnread
)

type (
	itemReadStatus int

	// Item is an entry in a feed
	Item struct {
		ID          int64
		FeedID      int64
		GUID        string
		URL         string
		Title       string
		Desc        string
		PublishedAt time.Time
		ReadAt      *time.Time
	}

	// ItemFilter is used to filter feed items in lists
	ItemFilter struct {
		FeedID     int64
		ReadStatus itemReadStatus
		Limit      int64
		Offset     int64
	}
)

// CreateIgnoreItems persists items to DB and if it already exists then skip it
// returns number of items inserted and error if any
func CreateIgnoreItems(db cruderExecer, items ...Item) (int64, error) {
	var values []string
	var params []interface{}

	if len(items) == 0 {
		return 0, errors.New("missing items to create")
	}

	for _, item := range items {
		values = append(values, "(?, ?, ?, ?, ?, ?)")
		params = append(params, item.FeedID, item.GUID, item.URL, item.Title, item.Desc, item.PublishedAt)
	}

	r, err := db.Exec(
		fmt.Sprintf(`INSERT OR IGNORE INTO "%s" (feed_id, guid, url, title, desc, published_at) VALUES %s`, itemsTable, strings.Join(values, ",")),
		params...,
	)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

// CountTotalByFeed returns the total number of items for a feed
func CountTotalByFeed(db cruderQueryRower, feedID int64) (int64, error) {
	var total int64

	err := db.QueryRow(
		fmt.Sprintf(`SELECT COUNT(id) FROM "%s" WHERE feed_id = ?`, itemsTable),
		feedID,
	).Scan(&total)

	return total, err
}

// CountUnreadByFeed returns the total number of unread items for a feed
func CountUnreadByFeed(db cruderQueryRower, feedID int64) (int64, error) {
	var unread int64

	err := db.QueryRow(
		fmt.Sprintf(`SELECT COUNT(id) FROM "%s" WHERE feed_id = ? AND read_at IS NULL`, itemsTable),
		feedID,
	).Scan(&unread)

	return unread, err
}

// ListItems returns a list of items from DB
func ListItems(db cruderQueryer, filter ItemFilter) ([]*Item, error) {
	var items []*Item
	var wheres []string
	var whereSQL, limitSQL string
	var params []interface{}

	if filter.FeedID > 0 {
		wheres = append(wheres, "feed_id = ?")
		params = append(params, filter.FeedID)
	}

	if filter.ReadStatus == ItemRead {
		wheres = append(wheres, "read_at IS NOT NULL")
	} else if filter.ReadStatus == ItemUnread {
		wheres = append(wheres, "read_at IS NULL")
	}

	if len(wheres) > 0 {
		whereSQL = " WHERE " + strings.Join(wheres, " AND ")
	}

	if filter.Limit > 0 {
		limitSQL = fmt.Sprintf(" LIMIT %d", filter.Limit)

		if filter.Offset > 0 {
			limitSQL = fmt.Sprintf("%s %d", limitSQL, filter.Offset)
		}
	}

	rows, err := db.Query(
		fmt.Sprintf(`SELECT * FROM "%s"%s ORDER BY published_at%s`, itemsTable, whereSQL, limitSQL),
		params...,
	)
	if err != nil {
		return items, err
	}
	defer rows.Close()
	for rows.Next() {
		i := &Item{}
		err = rows.Scan(&i.ID, &i.FeedID, &i.GUID, &i.URL, &i.Title, &i.Desc, &i.PublishedAt, &i.ReadAt)
		if err != nil {
			return items, err
		}

		items = append(items, i)
	}

	return items, nil
}

// SetItemsAsReadNow updates the read_at column for all items to CURRENT_TIMESTAMP
func SetItemsAsReadNow(db cruderExecer, ids ...int64) error {
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
		fmt.Sprintf(`UPDATE "%s" SET read_at = CURRENT_TIMESTAMP WHERE id IN (%s)`, itemsTable, strings.Join(placeholders, ",")),
		params...,
	)

	return err
}

// DeleteItems removes one or more items from DB
func DeleteItems(db cruderExecer, ids ...int64) error {
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
		fmt.Sprintf(`DELETE FROM "%s" WHERE id IN (%s)`, itemsTable, strings.Join(placeholders, ",")),
		params...,
	)

	return err
}
