package postgres

import "fmt"

// table name for torrents
const tableTorrents = "torrents"

// table name for users
const tableUsers = "users"

// table for new comments
const tableComments = "comments"

// table for old comments
const tableOldComments = "comments_old"

// table for torrent reports
const tableTorrentReports = "torrent_reports"

// table for user follows
const tableUserFollows = "user_follows"

// table for old user uploads
const tableOldUserUploads = "user_uploads_old"

// all tables that we have in current database schema in the order they are created
var tables = []createTable{
	// users table
	createTable{
		name: tableUsers,
		columns: tableColumns{
			"user_id SERIAL PRIMARY KEY",
			"username TEXT NOT NULL",
			"password TEXT NOT NULL",
			"email TEXT",
			"status INTEGER NOT NULL",
			"created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL",
			"updated_at TIMESTAMP WITHOUT TIME ZONE",
			"last_login_at TIMESTAMP WITHOUT TIME ZONE",
			"last_login_ip TEXT",
			"api_token TEXT",
			"api_token_expiry TIMESTAMP WITHOUT TIME ZONE NOT NULL",
			"language TEXT",
			"MD5 TEXT",
		},
		postCreate: []sqlQuery{
			createIndex(tableUsers, "username"),
		},
	},
	// torrents table
	createTable{
		name: tableTorrents,
		columns: tableColumns{
			"torrent_id SERIAL PRIMARY KEY",
			"torrent_name TEXT",
			"torrent_hash TEXT NOT NULL",
			"category INTEGER NOT NULL",
			"sub_category INTEGER NOT NULL",
			"status INTEGER NOT NULL",
			"date TIMESTAMP WITHOUT TIME ZONE",
			fmt.Sprintf("uploader INTEGER NOT NULL REFERENCES %s (user_id)", tableUsers),
			"downloads INTEGER",
			"stardom INTEGER NOT NULL",
			"filesize BIGINT",
			"description TEXT NOT NULL",
			"website_link TEXT",
			"deleted_at TIMESTAMP WITHOUT TIME ZONE",
			"seeders INTEGER",
			"leechers INTEGER",
			"completed INTEGER",
			"last_scrape TIMESTAMP WITHOUT TIME ZONE",
		},
		postCreate: []sqlQuery{
			createIndex(tableTorrents, "torrent_id"),
			createIndex(tableTorrents, "deleted_at"),
			createIndex(tableTorrents, "uploader"),
			createTrigraph(tableTorrents, "category", "torrent_name"),
			createTrigraph(tableTorrents, "sub_category", "torrent_name"),
			createTrigraph(tableTorrents, "status", "torrent_name"),
			createTrigraph(tableTorrents, "torrent_name"),
		},
	},
	// new comments table
	createTable{
		name: tableComments,
		columns: tableColumns{
			"comment_id SERIAL PRIMARY KEY",
			fmt.Sprintf("torrent_id INTEGER REFERENCES %s (torrent_id)", tableTorrents),
			fmt.Sprintf("user_id INTEGER REFERENCES %s (user_id)", tableUsers),
			"content TEXT NOT NULL",
			"created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL",
			"updated_at TIMESTAMP WITHOUT TIME ZONE",
			"deleted_at TIMESTAMP WITH TIME ZONE",
		},
		postCreate: []sqlQuery{
			createIndex(tableComments, "torrent_id"),
		},
	},
	// old comments table
	createTable{
		name: tableOldComments,
		columns: tableColumns{
			fmt.Sprintf("torrent_id INTEGER NOT NULL REFERENCES %s (torrent_id)", tableTorrents),
			"username TEXT NOT NULL",
			"content TEXT NOT NULL",
			"date TIMESTAMP WITHOUT TIME ZONE NOT NULL",
		},
	},
	// torrent reports table
	createTable{
		name: tableTorrentReports,
		columns: tableColumns{
			"torrent_report_id SERIAL PRIMARY KEY",
			"type TEXT",
			fmt.Sprintf("torrent_id INTEGER REFERENCES %s (torrent_id)", tableTorrents),
			"user_id INTEGER",
			"created_at TIMESTAMP WITH TIME ZONE",
		},
		postCreate: []sqlQuery{
			createIndex(tableTorrentReports, "torrent_report_id"),
		},
	},
	// user follows table
	createTable{
		name: tableUserFollows,
		columns: tableColumns{
			"user_id INTEGER NOT NULL",
			"following INTEGER NOT NULL",
			"PRIMARY KEY(user_id, following)",
		},
	},
	// old uploads table
	createTable{
		name: tableOldUserUploads,
		columns: tableColumns{
			"username TEXT IS NOT NULL",
			fmt.Sprintf("torrent_id INTEGER IS NOT NULL REFERENCES %s (torrent_id)", tableTorrents),
			"PRIMARY KEY (torrent_id)",
		},
	},
}
