package simple


import (
	"fmt"
	"gorm.io/gorm"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/postgres"
)

type DialectorFactory func(dsn string) gorm.Dialector

var dialectorRegistry = map[string]DialectorFactory{
	"sqlite":   sqlite.Open,
	"mysql":    mysql.Open,
	"postgres": postgres.Open,
}


func OpenSqlite(name string) (*gorm.DB) {
	db, err := OpenDatabase("sqlite", name)
	if err != nil {
		fmt.Println("Database opended")
	}
	return db
}

func OpenSqliteMemory() *gorm.DB {
	db, err := OpenDatabase("sqlite", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("failed to open in-memory sqlite database: %v", err))
	}
	return db
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
