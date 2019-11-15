package sqlite_test

import (
	"database/sql"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"feeda/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

const (
	testFeedURL    = "https://www.example.com"
	testFeedURL2   = "https://www.example2.com"
	testItemGUID   = "guid"
	testItemGUID2  = "guid2"
	testItemGUID3  = "guid3"
	testItemURL    = testFeedURL + "/item.html"
	testItemURL2   = testFeedURL + "/item2.html"
	testItemURL3   = testFeedURL + "/item3.html"
	testItemTitle  = "title"
	testItemTitle2 = "title2"
	testItemTitle3 = "title3"
	testItemDesc   = "desc"
	testItemDesc2  = "desc2"
	testItemDesc3  = "desc3"
)

var (
	db    *sql.DB
	err   error
	tmpDB = path.Join(os.TempDir(), "feeda_test.db")
)

func init() {
	log.Printf("tmp db at %s\n", tmpDB)
	os.Remove(tmpDB)

	db, err = sql.Open("sqlite3", tmpDB)
	if err != nil {
		log.Fatal(err)
	}

	sqlite.EnsureTables(db)
}

func TestCRUD(t *testing.T) {
	start := time.Now()

	feed1 := sqlite.Feed{URL: testFeedURL, Type: sqlite.FeedTypeRSS}
	feed2 := sqlite.Feed{URL: testFeedURL2, Type: sqlite.FeedTypeAtom}

	// Add feed
	err = sqlite.CreateIgnoreFeeds(db, feed1)
	if err != nil {
		t.Fatal(err)
	}

	// Add feed 2 and should ignore feed 1
	err = sqlite.CreateIgnoreFeeds(db, feed1, feed2)
	if err != nil {
		t.Fatal(err)
	}

	// List feeds
	var feeds []*sqlite.Feed
	feeds, err = sqlite.ListFeeds(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(feeds) != 2 {
		t.Fatalf("expecting length of feeds to be 2, got %d", len(feeds))
	}
	if feeds[0].ID != 1 {
		t.Fatalf("expecting id of feed to be 1, got %d", feeds[0].ID)
	}
	if feeds[1].ID != 3 {
		t.Fatalf("expecting id of feed to be 3, got %d", feeds[1].ID)
	}
	if feeds[0].URL != testFeedURL {
		t.Fatalf("expecting url of feed to be %s, got %s", testFeedURL, feeds[0].URL)
	}
	if feeds[1].URL != testFeedURL2 {
		t.Fatalf("expecting url of feed to be %s, got %s", testFeedURL2, feeds[1].URL)
	}
	if feeds[0].Type != sqlite.FeedTypeRSS {
		t.Fatalf("expecting type of feed to be RSS, got %s", feeds[0].Type)
	}
	if feeds[1].Type != sqlite.FeedTypeAtom {
		t.Fatalf("expecting url of feed to be Atom, got %s", feeds[1].Type)
	}
	if !feeds[0].CreatedAt.After(start.Add(-1 * time.Minute)) {
		t.Fatalf("expected created_at to be after %s, got %s", start, feeds[0].CreatedAt)
	}
	if !feeds[1].CreatedAt.Equal(feeds[0].CreatedAt) {
		t.Fatalf("expected created_at to be after %s, got %s", feeds[0].CreatedAt, feeds[1].CreatedAt)
	}
	if feeds[0].SyncedAt != nil {
		t.Fatalf("expecting synced_at to be nil, got %s", feeds[0].SyncedAt)
	}
	if feeds[1].SyncedAt != nil {
		t.Fatalf("expecting synced_at to be nil, got %s", feeds[1].SyncedAt)
	}

	// Set feed 1 and 2 as synced
	err = sqlite.SetFeedsSyncedAtNow(db, feeds[0].ID, feeds[1].ID)
	if err != nil {
		t.Fatal(err)
	}

	// List feeds
	feeds, err = sqlite.ListFeeds(db, feeds[0].ID, feeds[1].ID)
	if err != nil {
		t.Fatal(err)
	}

	// Check that feeds 1 and 2 has synced_at set
	if feeds[0].SyncedAt == nil || !feeds[0].SyncedAt.After(start.Add(-1*time.Minute)) {
		t.Fatalf("expecting read_at to be %s, got %s", start, feeds[0].SyncedAt)
	}
	if feeds[1].SyncedAt == nil || !feeds[1].SyncedAt.After(start.Add(-1*time.Minute)) {
		t.Fatalf("expecting read_at to be %s, got %s", start, feeds[1].SyncedAt)
	}

	// Add item 1
	affected, err := sqlite.CreateIgnoreItems(db, sqlite.Item{
		FeedID:      feeds[0].ID,
		GUID:        testItemGUID,
		URL:         testItemURL,
		Title:       testItemTitle,
		Desc:        testItemDesc,
		PublishedAt: start,
	})
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatalf("expecting affected to be 1, got %d", affected)
	}

	// Add item 2 and 3 and should ignore item 1
	affected, err = sqlite.CreateIgnoreItems(db,
		sqlite.Item{
			FeedID:      feeds[0].ID,
			GUID:        testItemGUID,
			URL:         testItemURL,
			Title:       testItemTitle,
			Desc:        testItemDesc,
			PublishedAt: start,
		},
		sqlite.Item{
			FeedID:      feeds[0].ID,
			GUID:        testItemGUID2,
			URL:         testItemURL2,
			Title:       testItemTitle2,
			Desc:        testItemDesc2,
			PublishedAt: start.Add(1 * time.Minute),
		},
		sqlite.Item{
			FeedID:      feeds[1].ID,
			GUID:        testItemGUID3,
			URL:         testItemURL3,
			Title:       testItemTitle3,
			Desc:        testItemDesc3,
			PublishedAt: start.Add(2 * time.Minute),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if affected != 2 {
		t.Fatalf("expecting affected to be 2, got %d", affected)
	}

	// List items with empty filter
	var items []*sqlite.Item
	items, err = sqlite.ListItems(db, sqlite.ItemFilter{})
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 3 {
		t.Fatalf("expecting length of items to be 3, got %d", len(items))
	}
	if items[0].ID != 1 {
		t.Fatalf("expecting id of feed to be 1, got %d", items[0].ID)
	}
	if items[1].ID != 3 {
		t.Fatalf("expecting id of feed to be 3, got %d", items[1].ID)
	}
	if items[2].ID != 4 {
		t.Fatalf("expecting id of feed to be 4, got %d", items[2].ID)
	}
	if items[0].FeedID != feeds[0].ID {
		t.Fatalf("expecting feed_id of item to be %d, got %d", feeds[0].ID, items[0].FeedID)
	}
	if items[1].FeedID != feeds[0].ID {
		t.Fatalf("expecting feed_id of item to be %d, got %d", feeds[0].ID, items[1].FeedID)
	}
	if items[2].FeedID != feeds[1].ID {
		t.Fatalf("expecting feed_id of item to be %d, got %d", feeds[1].ID, items[2].FeedID)
	}
	if items[0].GUID != testItemGUID {
		t.Fatalf("expecting url of item to be %s, got %s", testItemGUID, items[0].GUID)
	}
	if items[1].GUID != testItemGUID2 {
		t.Fatalf("expecting url of item to be %s, got %s", testItemGUID2, items[1].GUID)
	}
	if items[2].GUID != testItemGUID3 {
		t.Fatalf("expecting url of item to be %s, got %s", testItemGUID3, items[2].GUID)
	}
	if items[0].URL != testItemURL {
		t.Fatalf("expecting url of item to be %s, got %s", testItemURL, items[0].URL)
	}
	if items[1].URL != testItemURL2 {
		t.Fatalf("expecting url of item to be %s, got %s", testItemURL2, items[1].URL)
	}
	if items[2].URL != testItemURL3 {
		t.Fatalf("expecting url of item to be %s, got %s", testItemURL3, items[2].URL)
	}
	if items[0].Title != testItemTitle {
		t.Fatalf("expecting desc of item to be %s, got %s", testItemTitle, items[0].Title)
	}
	if items[1].Title != testItemTitle2 {
		t.Fatalf("expecting desc of item to be %s, got %s", testItemTitle2, items[1].Title)
	}
	if items[2].Title != testItemTitle3 {
		t.Fatalf("expecting desc of item to be %s, got %s", testItemTitle3, items[2].Title)
	}
	if items[0].Desc != testItemDesc {
		t.Fatalf("expecting desc of item to be %s, got %s", testItemDesc, items[0].Desc)
	}
	if items[1].Desc != testItemDesc2 {
		t.Fatalf("expecting desc of item to be %s, got %s", testItemDesc2, items[1].Desc)
	}
	if items[2].Desc != testItemDesc3 {
		t.Fatalf("expecting desc of item to be %s, got %s", testItemDesc3, items[2].Desc)
	}
	if !items[0].PublishedAt.Equal(start) {
		t.Fatalf("expected published_at to equal %s, got %s", start, items[0].PublishedAt)
	}
	if !items[1].PublishedAt.Equal(start.Add(1 * time.Minute)) {
		t.Fatalf("expected published_at to equal %s, got %s", start.Add(1*time.Minute), items[1].PublishedAt)
	}
	if !items[2].PublishedAt.Equal(start.Add(2 * time.Minute)) {
		t.Fatalf("expected published_at to equal %s, got %s", start.Add(2*time.Minute), items[2].PublishedAt)
	}
	if items[0].ReadAt != nil {
		t.Fatalf("expecting read_at to be nil, got %s", items[0].ReadAt)
	}
	if items[1].ReadAt != nil {
		t.Fatalf("expecting read_at to be nil, got %s", items[1].ReadAt)
	}
	if items[2].ReadAt != nil {
		t.Fatalf("expecting read_at to be nil, got %s", items[2].ReadAt)
	}

	// Set item 1 and 3 as read
	err = sqlite.SetItemsAsReadNow(db, items[0].ID, items[2].ID)
	if err != nil {
		t.Fatal(err)
	}

	// List items with empty filter
	items, err = sqlite.ListItems(db, sqlite.ItemFilter{})
	if err != nil {
		t.Fatal(err)
	}

	// Check that item 1 and 3 is set as read
	if items[0].ReadAt == nil || !items[0].ReadAt.After(start.Add(-1*time.Minute)) {
		t.Fatalf("expecting read_at to be %s, got %s", start, items[0].ReadAt)
	}
	if items[1].ReadAt != nil {
		t.Fatalf("expecting synced_at to be nil, got %s", items[1].ReadAt)
	}
	if items[2].ReadAt == nil || !items[2].ReadAt.After(start.Add(-1*time.Minute)) {
		t.Fatalf("expecting read_at to be %s, got %s", start, items[2].ReadAt)
	}

	total, err := sqlite.CountTotalByFeed(db, feeds[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 {
		t.Fatalf("expecting total of feed %d to be 2, got %d", feeds[0].ID, total)
	}

	unread, err := sqlite.CountUnreadByFeed(db, feeds[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if unread != 1 {
		t.Fatalf("expecting unread of feed %d to be 1, got %d", feeds[0].ID, unread)
	}

	// Delete item
	err = sqlite.DeleteItems(db, items[0].ID, items[1].ID, items[2].ID)
	if err != nil {
		t.Fatal(err)
	}

	// List items again should return empty list
	items, err = sqlite.ListItems(db, sqlite.ItemFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("expecting length of items to be 0, got %d", len(items))
	}

	// Delete feed
	err = sqlite.DeleteFeeds(db, feeds[0].ID, feeds[1].ID)
	if err != nil {
		t.Fatal(err)
	}

	// List feeds again should return empty list
	feeds, err = sqlite.ListFeeds(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(feeds) != 0 {
		t.Fatalf("expecting length of feeds to be 0, got %d", len(feeds))
	}
}

func TestCleanup(t *testing.T) {
	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	os.Remove(tmpDB)
}
