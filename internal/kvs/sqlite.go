package kvs

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite() *SQLite {
	return &SQLite{}
}

func (s *SQLite) Name() string {
	return "sqlite"
}

func (s *SQLite) Open(path string) error {
	dbPath := filepath.Join(path, "sqlite.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS kvs (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return err
	}
	
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_kvs_key ON kvs(key)`)
	if err != nil {
		db.Close()
		return err
	}
	
	s.db = db
	return nil
}

func (s *SQLite) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *SQLite) Set(key string, value *Value) error {
	data, err := value.ToJSON()
	if err != nil {
		return err
	}
	
	_, err = s.db.Exec(
		"INSERT OR REPLACE INTO kvs (key, value) VALUES (?, ?)",
		key, string(data),
	)
	return err
}

func (s *SQLite) Get(key string) (*Value, error) {
	var data string
	err := s.db.QueryRow("SELECT value FROM kvs WHERE key = ?", key).Scan(&data)
	if err != nil {
		return nil, err
	}
	
	return ValueFromJSON([]byte(data))
}