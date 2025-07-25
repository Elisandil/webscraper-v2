package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SQLiteDB struct {
	*sql.DB
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	sqliteDB := &SQLiteDB{DB: db}

	if err := sqliteDB.createTables(); err != nil {
		return nil, err
	}
	return sqliteDB, nil
}

func (db *SQLiteDB) createTables() error {
	scrapingQuery := `
	CREATE TABLE IF NOT EXISTS scraping_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		url TEXT NOT NULL,
		title TEXT,
		description TEXT,
		keywords TEXT,
		author TEXT,
		language TEXT,
		favicon TEXT,
		image_url TEXT,
		site_name TEXT,
		links TEXT,
		images TEXT,
		headers TEXT,
		status_code INTEGER,
		content_type TEXT,
		word_count INTEGER DEFAULT 0,
		load_time_ms INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	CREATE INDEX IF NOT EXISTS idx_scraping_results_url ON scraping_results(url);
	CREATE INDEX IF NOT EXISTS idx_scraping_results_created_at ON scraping_results(created_at);
	CREATE INDEX IF NOT EXISTS idx_scraping_results_status_code ON scraping_results(status_code);`

	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		active BOOLEAN NOT NULL DEFAULT true,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_active ON users(active);`

	triggerQuery := `
	CREATE TRIGGER IF NOT EXISTS users_updated_at
	AFTER UPDATE ON users
	FOR EACH ROW
	BEGIN
		UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;`

	if _, err := db.Exec(scrapingQuery); err != nil {
		return err
	}
	if _, err := db.Exec(usersQuery); err != nil {
		return err
	}
	if _, err := db.Exec(triggerQuery); err != nil {
		return err
	}
	return nil
}
