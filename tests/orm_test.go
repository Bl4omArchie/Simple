package test

import (
    "context"
    "testing"

    "gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"github.com/Bl4omArchie/simple"
)

type User struct {
    ID   int
    Name string
}

func TestBasicOperationOnSqlite(t *testing.T) {
    gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open gorm db: %v", err)
    }
	gormDB.AutoMigrate(&User{})

	alice := User{ID: 1, Name: "Alice"}
	bob := User{ID: 2, Name: "Bob"}

	gormDB.Create(&alice)
	gormDB.Create(&bob)

    ctx := context.Background()
	// ===== Get Alice =====
    user, err := simple.GetRowBy[User](ctx, gormDB, "ID", 1)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("expected Alice, got %s", user.Name)
    }
	// ==========


	// Count rows
	count, err := simple.CountRows[User](ctx, gormDB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2, got : %v", err)
	}
	// ==========


	// Delete Alice
	if err = simple.DeleteRowBy[User](ctx, gormDB, "name", "Alice"); err !=  nil {
		t.Fatalf("couldn't delete user Alice : %v", err)
	}
	// ==========


	// Rename Bob into MegaBob
	if err = simple.UpdateRowBy[User](ctx, gormDB, "name", "Bob", "name", "MegaBob"); err !=  nil {
		t.Fatalf("couldn't rename user Bob : %v", err)
	}

	user, err = simple.GetRowBy[User](ctx, gormDB, "name", "MegaBob")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "MegaBob" {
        t.Errorf("expected MegaBob, got %s", user.Name)
    }
	// ==========
}
