package simple

import (
    "fmt"
	"context"

    "gorm.io/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/postgres"
)

type FactoryDB func(dsn string) gorm.Dialector

var registryDB = map[string]FactoryDB {
	"sqlite":   sqlite.Open,
	"mysql":    mysql.Open,
	"postgres": postgres.Open,
}


// Open a database manually or with function GetMysql(), GetPostgres() anb GetSqlite()
func OpenDatabase(driver, dsn string) (*gorm.DB, error) {
	factory, ok := registryDB[driver]
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

// Database migration
func Migrate(ctx context.Context, odb *gorm.DB, models ...any) error {
    err := odb.WithContext(ctx).AutoMigrate(models...)
    return err
}

// GetRowBy returns the first record of type T where column `key` equals `value`.
func GetRowBy[T any](ctx context.Context, db *gorm.DB, key string, value any) (*T, error) {
	var model T
	if err := db.WithContext(ctx).
		Where("? = ?", key, value).
		First(&model).Error; err != nil {
			return nil, fmt.Errorf("GetRowBy for key `%s` and value `%s`. %w", key, value, err)
	}
	return &model, nil
}

// GetRowsBy returns every record of type T where column `key` equals `value`.
func GetRowsBy[T any](ctx context.Context, db *gorm.DB, key string, value any) ([]T, error) {
	var model []T
	if err := db.WithContext(ctx).
		Where("? = ?", key, value).
		Find(&model).Error; err != nil {
			return nil, fmt.Errorf("GetRowsBy for key `%s` and value `%s`. %w", key, value, err)
	}
	return model, nil
}

// GetTable return every rows of the given table `T`
// Set the limit of rows with `limit`. Use -1 to get every rows
func GetTable[T any](ctx context.Context, db *gorm.DB, limit int) ([]T, error) {
	var models []T
	if err := db.WithContext(ctx).
		Limit(limit).
		Find(&models).
		Error; err != nil {
			return nil, fmt.Errorf("GetTable: %w", err)
	}
	return models, nil
}

// GetColumn returns a specific column of your table `T`of type `C`.
func GetColumn[T any, C any](ctx context.Context, odb *gorm.DB, columnName string) ([]C, error) {
	var model T
	var values []C
	if err := odb.WithContext(ctx).
		Model(&model).
		Pluck(columnName, &values).
		Error; err != nil {
			return nil, fmt.Errorf("GetColumn for columnName : `%s`. %w", columnName, err)
	}
	return values, nil
}

// Updates rows where column `key` equals `value`
func UpdateRowBy[T any](ctx context.Context, db *gorm.DB, key string, value any, field string, newValue any) error {
	if err := db.WithContext(ctx).
		Model(new(T)).
		Where(fmt.Sprintf("%s = ?", key), value).
		Update(field, newValue).Error; err != nil {
			return fmt.Errorf("UpdateRowBy: %w", err)
	}
    return nil
}

// Deletes rows where column `key` equals `value`
func DeleteRowBy[T any](ctx context.Context, db *gorm.DB, key string, value any) error {
	if err := db.WithContext(ctx).
		Where(fmt.Sprintf("%s = ?", key), value).
		Delete(new(T)).Error; err != nil {
		return fmt.Errorf("DeleteRowBy: %w", err)
	}
    return nil
}

// Counts how much rows there is in your table
func CountRows[T any](ctx context.Context, db *gorm.DB) (int64, error) {
    var count int64

    if err := db.WithContext(ctx).
		Model(new(T)).
		Count(&count).Error; err != nil {
        	return 0, fmt.Errorf("CountRows: %w", err)
    	}
    return count, nil
}


func GetMysql(user string, password string, host string, port string, dbName string) (string, string) {
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
}

func GetPostgres(host string, user string, password string, dbName string, port int, sslMode string, timezone string) (string, string) {
	return "postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, dbName, port, sslMode, timezone)
}

func GetSqlite(filePath string) (string, string) {
	return "sqlite", filePath
}
