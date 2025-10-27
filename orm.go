package simple

import (
    "fmt"
    "path/filepath"

    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)


func OpenSqliteDatabase(dbPath string) (*gorm.DB, error) {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Incorrect path : %s", dbPath)
	}
	db, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{})
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("invalid inputs: %w", err)
	}
	return &model, nil
}

func GetTable[T any](odb *gorm.DB) ([]*T, error) {
	var model []*T
	if err := odb.Find(&model).Error; err != nil {
		return nil, fmt.Errorf("invalid inputs: %w", err)
	}
	return model, nil
}

func UpdateTable[T any](odb *gorm.DB, key string, value string, newColumn string, newValue string) (*T, error) {
	var model T
	if err := odb.Model(&model).
		Where(fmt.Sprintf("%s = ?", key), value).
		Update(newColumn, newValue).Error; err != nil {
		return nil, fmt.Errorf("couldn't update row %s = %s: %w", key, value, err)
	}
	return &model, nil
}


// Todo
func GetMysql(user string, password string, host string, port string, dbName string) (string, string) {
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
}

func GetPostgres(host string, user string, password string, dbName string, port int, sslMode string, timezone string) (string, string) {
	return "postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s, TimeZone=%s", host, user, password, dbName, port, sslMode, timezone)
}

func GetSqlite(filePath string) (string, string) {
	return "sqlite", filePath
}
