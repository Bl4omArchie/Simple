package simple


import (
	"fmt"
	"gorm.io/gorm"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/gaussdb"
	"gorm.io/driver/postgres"
	"gorm.io/driver/clickhouse"
	"github.com/joho/godotenv"
)

type DialectorFactory func(dsn string) gorm.Dialector

var dialectorRegistry = map[string]DialectorFactory{
	"sqlite":   sqlite.Open,
	"mysql":    mysql.Open,
	"postgres": postgres.Open,
	"gaussdb": gaussdb.Open,
	"clickhouse": clickhouse.Open,
}


// TODO
func OpenEnv() (*gorm.DB) {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Some error occured. Err: %s", err)
	}
	return nil
}

func OpenGaussdb(host string, user string, password string, dbName string, port int, sslMode string) (*gorm.DB, error) {
	db, err := OpenDatabase("gaussdb", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", host, user, password, dbName, port, sslMode))
	return db, err
}

func OpenClickhouse( user string, password string, host string, port int, dbName string) (*gorm.DB, error) {
	db, err := OpenDatabase("clickhouse", fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s?read_timeout=10s&write_timeout=20s", user, password, host, port, dbName))
	return db, err
}

func OpenSqlite(name string) (*gorm.DB, error) {
	db, err := OpenDatabase("sqlite", name)
	return db, err
}

func OpenSqliteMemory() (*gorm.DB, error) {
	db, err := OpenDatabase("sqlite", ":memory:")
	return db, err
}

func OpenMysql(user string, password string, host string, port int, dbName string) (*gorm.DB, error) {
	db, err := OpenDatabase("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName))
	return db, err
}

func OpenMysqlUnixSocket(user, password, socket, dbName string) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", user, password, socket, dbName)
    return OpenDatabase("mysql", dsn)
}

func OpenPostgres(host string, user string, password string, dbName string, port int, sslMode string) (*gorm.DB, error) {
	db, err := OpenDatabase("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", host, user, password, dbName, port, sslMode))
	return db, err
}

func OpenDatabase(driver string, dsn string) (*gorm.DB, error) {
	factory, ok := dialectorRegistry[driver]
	if !ok {
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
	dialector := factory(dsn)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

func MigrateTables(db *gorm.DB, tables ...any) error {
	if err := db.AutoMigrate(tables...); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}

func InsertRows(db *gorm.DB, entries ...any) error {
	for _, entry := range entries {
		if err := db.Create(entry).Error; err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}
	}
	return nil
}
