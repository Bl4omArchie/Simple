package simple

import (
    "fmt"

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


func OpenDatabase(driver string, dsn string) (*gorm.DB, error) {
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

func Migrate(odb *gorm.DB, models ...any) error {
    err := odb.AutoMigrate(models...)
    return err
}

func GetBy[T any](odb *gorm.DB, key string, value string) (*T, error) {
	var model T
	if err := odb.First(&model, fmt.Sprintf("%s = ?", key), value).Error; err != nil {
		return nil, fmt.Errorf("GetBy, invalid inputs: %w", err)
	}
	return &model, nil
}

func GetTable[T any](odb *gorm.DB) ([]T, error) {
	var model []T
	if err := odb.Find(&model).Error; err != nil {
		return nil, fmt.Errorf("GetTable, invalid inputs: %w", err)
	}
	return model, nil
}

func GetColumn[T any, C any](odb *gorm.DB, columnName string) ([]C, error) {
	var model T
	var values []C
	if err := odb.Model(&model).Pluck(columnName, &values).Error; err != nil {
		return nil, fmt.Errorf("GetColumn, invalid inputs: %w", err)
	}
	return values, nil
}

func UpdateColumnWhereValue[T any](odb *gorm.DB, key string, value string, newColumn string, newValue string) (*T, error) {
	var model T
	if err := odb.Model(&model).
		Where(fmt.Sprintf("%s = ?", key), value).
		Update(newColumn, newValue).Error; err != nil {
		return nil, fmt.Errorf("couldn't update row %s = %s: %w", key, value, err)
	}
	return &model, nil
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
