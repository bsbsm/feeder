package db

import (
	"database/sql"
	"errors"
	"net/url"
	"sync"

	"github.com/bsbsm/feeder/pkg/feeder"
	_ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = errors.New("Not found")
var ErrIncorrectArgs = errors.New("Incorrect arguments")

type SQLiteDatabase struct {
}

// GetNews returns specific news count from database
func (s *SQLiteDatabase) GetNews(offset int, count int) ([]*News, error) {
	return readNews(getDb(), count, offset)
}

// GetNewsWithTitle returns specific news count from database
func (s *SQLiteDatabase) GetNewsWithTitle(title string, offset int, count int) ([]*News, error) {
	return readNewsWithTitle(getDb(), title, count, offset)
}

// GetNewsDetail returns detail for news
func (s *SQLiteDatabase) GetNewsDetail(id int) (*NewsDetail, error) {
	return readNewsDetail(getDb(), id)
}

// CreateNews insert news to database and return errors if need
func (s *SQLiteDatabase) CreateNews(sourceID int, title string, payloadJSON []byte) error {
	return writeNews(getDb(), sourceID, title, payloadJSON)
}

// GetFeedSources returns all feed sources
func (s *SQLiteDatabase) GetFeedSources() ([]*feeder.FeedSource, error) {
	return readFeedSources(getDb())
}

// CreateFeedSource insert new feed source to database and return errors if need
func (s *SQLiteDatabase) CreateFeedSource(url, rule string) error {
	return writeFeedSource(getDb(), url, rule)
}

type News struct {
	Title  string `json:"Title"`
	Source string `json:"Source"`
	ID     int    `json:"ID"`
}

type NewsDetail struct {
	Title       string `json:"Title"`
	PayloadJSON string `json:"PayloadJSON"`
	Source      string `json:"Source"`
}

// getDb returns database connection pool
func getDb() *sql.DB {
	if databaseInstance != nil {
		return databaseInstance
	}

	dbMut.Lock()
	if databaseInstance == nil {
		databaseInstance = initDB()
	}
	dbMut.Unlock()

	return databaseInstance
}

var (
	dbMut            sync.Mutex
	databaseInstance *sql.DB
)

var databaseFilePath = "./local.db"

// initDB prepare database before using
func initDB() *sql.DB {
	//os.Remove(databaseFilePath)

	db, err := sql.Open("sqlite3", databaseFilePath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db is nil")
	}
	err = createTables(db)

	if err != nil {
		db.Close()
		db = nil
		panic(err)
	}

	return db
}

// createTable creates needed tables if its not exist
func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS news(
		ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Title TEXT UNIQUE,
		PayloadJSON TEXT,
		SourceID INTEGER NOT NULL,
		AddedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS sources(
		ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		URL TEXT NOT NULL UNIQUE,
		Rule TEXT NOT NULL
	);
	`
	_, err := db.Exec(query)

	return err
}

func writeNews(db *sql.DB, sourceID int, title string, payloadJSON []byte) error {
	if title == "" && len(payloadJSON) == 0 {
		return ErrIncorrectArgs
	}

	query := `
	INSERT INTO news(
		Title,
		PayloadJSON,
		SourceID
	) values(?, ?, ?);
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, payloadJSON, sourceID)
	if err != nil {
		return err
	}

	return nil
}

func readNews(db *sql.DB, limit, offset int) ([]*News, error) {
	query := `
	SELECT t1.ID, t1.Title, t2.URL FROM news t1
	LEFT JOIN sources t2 ON t1.SourceID = t2.ID
	ORDER BY t1.AddedAt ASC
	LIMIT ? OFFSET ?
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*News

	for rows.Next() {
		var item News
		err = rows.Scan(&item.ID, &item.Title, &item.Source)
		if err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	return result, nil
}

func readNewsWithTitle(db *sql.DB, title string, limit, offset int) ([]*News, error) {
	query := `
	SELECT t1.ID, t1.Title, t2.URL FROM news t1
	LEFT JOIN sources t2 ON t1.SourceID = t2.ID
	WHERE t1.Title like ?
	ORDER BY t1.AddedAt ASC
	LIMIT ? OFFSET ?
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%"+title+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*News

	for rows.Next() {
		var item News
		err = rows.Scan(&item.ID, &item.Title, &item.Source)
		if err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	return result, nil
}

func readNewsDetail(db *sql.DB, id int) (*NewsDetail, error) {
	query := `
	SELECT t1.Title, t1.PayloadJSON, t2.URL FROM news t1
	LEFT JOIN sources t2 ON t1.SourceID = t2.ID
	WHERE t1.ID = ?;
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item NewsDetail

		err = rows.Scan(&item.Title, &item.PayloadJSON, &item.Source)
		if err != nil {
			return nil, err
		}

		return &item, nil
	}

	return nil, ErrNotFound
}

func readFeedSources(db *sql.DB) ([]*feeder.FeedSource, error) {
	query := `
	SELECT ID, URL, Rule FROM sources
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*feeder.FeedSource

	for rows.Next() {
		item := feeder.FeedSource{}
		var ruleJSON string
		err = rows.Scan(&item.ID, &item.URL, &ruleJSON)
		if err != nil {
			return nil, err
		}

		if err = feeder.ImplementRule(&item, ruleJSON); err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	return result, nil
}

func writeFeedSource(db *sql.DB, u string, rule string) error {
	if _, err := url.ParseRequestURI(u); err != nil || rule == "" {
		return ErrIncorrectArgs
	}

	query := `
	INSERT INTO sources(
		URL,
		Rule
	) values(?, ?);
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u, rule)
	if err != nil {
		return err
	}

	return nil
}
