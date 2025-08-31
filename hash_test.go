package simple

import (
	"os"
	"testing"
)

// Test HashData
// Every hash is  from the string "abc"
func TestHashData(t *testing.T) {
    data := []byte("abc")

	expected := map[string]string{
		"md5":      "900150983cd24fb0d6963f7d28e17f72",
		"sha1":     "a9993e364706816aba3e25717850c26c9cd0d89d",
		"sha256":   "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		"sha384":   "cb00753f45a35e8bb5a03d699ac65007272c32ab0eded1631a8b605a43ff5bed8086072ba1e7cc2358baeca134c825a7",
		"sha512":   "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f",
		"sha3-256": "3a985da74fe225b2045c172d6bd390bd855f086e3e9d525b46bfe24511431532",
		"sha3-384": "ec01498288516fc926459f58e2c6ad8df9b473cb0fc08c2596da7cf0e49be4b298d88cea927ac7f539f1edf228376d25",
		"sha3-512": "b751850b1a57168a5693cd924b6b096e08f621827444f70d884f5d0240d2712e10e116e9192af3c91a7ec57647e3934057340b4cf408d5a56592f8274eec53f0",
	}

    for algo, exp := range expected {
        h, err := HashData(algo, data)
        if err != nil {
            t.Errorf(algo, "error:", err)
            continue
        }
        if h != exp {
            t.Errorf("FAILED : %s. Expected : %v, got : %v", algo, exp, h)
        }
    }
}

func TestCompareFiles(t *testing.T) {
	content := "Hello wORLD!"
	if err := os.WriteFile("test1.txt", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test1.txt")

	contentNew := "Hello wORLD!"
	if err := os.WriteFile("test2.txt", []byte(contentNew), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test2.txt")

	got, err := CompareFiles("sha256", "test1.txt", "test2.txt")
	if err != nil {
		t.Errorf("Error while hashing : %v", err)
	}
	want := true

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestCompareFilesKey(t *testing.T) {
	content := "Hello wORLD!"
	if err := os.WriteFile("test1.txt", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test1.txt")

	contentNew := "Hello wORLD!"
	if err := os.WriteFile("test2.txt", []byte(contentNew), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test2.txt")

	got, err := CompareFilesKey("blake2b-256", []byte("key1"), "test1.txt", []byte("key1"), "test2.txt")
	if err != nil {
		t.Errorf("Error while hashing : %v", err)
	}
	want := true

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
