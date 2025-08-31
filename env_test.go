package simple


import (
	"os"
	"testing"
)


// Test OpenEnv function
// write a temporary .env
func TestOpenEnv(t *testing.T) {
	content := "DB_HOST=test\nDB_PASS=password\n"
	if err := os.WriteFile(".env", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".env")

	got := OpenEnv("DB_HOST", "DB_PASS")
	want := []string{"test", "password"}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got %q, wanted %q", got[i], want[i])
		}
	}
}

// Test OpenEnvFilenames function
func TestOpenEnvFilenames(t *testing.T) {
	content := "DB_HOST=test\nDB_PASS=password\n"
	if err := os.WriteFile("test.env", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.env")

	contentNew := "DB_HOST2=test2\nDB_PASS2=password2\n"
	if err := os.WriteFile("db.env", []byte(contentNew), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("db.env")


	got := OpenEnvFilenames([]string{"test.env", "db.env"}, "DB_HOST", "DB_PASS", "DB_HOST2", "DB_PASS2")
	want := []string{"test", "password", "test2", "password2"}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got %q, wanted %q", got[i], want[i])
		}
	}
}
