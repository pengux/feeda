# feeda
Feeds (RSS2/Atom) aggregator as a CLI tool.

[![Build Status](https://travis-ci.org/pengux/feeda.svg?branch=master)](https://travis-ci.org/pengux/feeda)

## Installation

```sh
go get -u github.com/pengux/feeda
```

## Usage

```sh
Feeds are stored in a SQLite database which default to ~/.feeda/db.sqlite
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

Usage:
  feeda [command]

Available Commands:
  add         Add RSS feeds
  delete      Delete items
  deleteFeed  Delete feeds
  help        Help about any command
  list        List items from feeds
  listFeeds   List all feeds
  sync        Download latest items of one or multiple feeds

Flags:
      --db string   Location of DB, defaults to ~/.feeda/db.sqlite
  -h, --help        help for feeda

Use "feeda [command] --help" for more information about a command.
```

Use [cron](https://en.wikipedia.org/wiki/Cron) to sync your feeds regularly, for example:

```
*/10 * * * * feeda sync
```
