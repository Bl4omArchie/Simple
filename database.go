package simple


import (
	"fmt"
	"sync"
	"database/sql"
)


type DatabaseManager interface {
	Open(dbName string, dsn string) error
	Close() error
	Query(query string, args ...string) (*sql.Rows, error)
}

// Default database manager
type Database struct {
	Dsn string
	Conn *sql.DB
	Driver string
	PingSuccess bool
	mu *sync.Mutex
}

func OpenDatabase(driver string, dsn string) (*Database, error) {
	db := &Database{}
	if err := db.Open(driver, dsn); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return db, nil
}

func (db *Database) Open(driver string, dsn string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("Couldn't connect to the database : %w", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		db.PingSuccess = false
		return fmt.Errorf("ping to database failed: %w", err)
	}

	db.Dsn = dsn
	db.Conn = conn
	db.Driver = driver
	db.PingSuccess = true
	return nil
}

func (db *Database) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.Conn != nil {
		db.Conn = nil
		return db.Conn.Close()
	}
	return fmt.Errorf("Error while closing database : database not initialized")
}

// Issue : leaking connection with rows.
func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.Conn.Query(query, args)
	if err != nil {
		return nil, fmt.Errorf("query has failed : %w", err)
	}
	return rows, nil
}

func GetMysql(user string, password string, host string, port string, dbName string) (string, string) {
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
}

func GetPostgres(host string, user string, password string, dbName string, port int, sslMode string, timezone string) (string, string) {
	return "postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s, TimeZone=%s", host, user, password, dbName, port, sslMode, timezone)
}

func GetSqlite(filePath string) (string, string) {
	return "sqlite", filePath
}
