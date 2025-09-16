package simple

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
)


// Implementation of hash.Hash
type shakeAdapter struct {
	shake    sha3.ShakeHash
	length int
}
var _ hash.Hash = (*shakeAdapter)(nil)


func (x *shakeAdapter) Write(p []byte) (int, error) { return x.shake.Write(p) }
func (x *shakeAdapter) Sum(b []byte) []byte {
    out := make([]byte, x.length)
    shakeCopy := x.shake.Clone()
    io.ReadFull(shakeCopy, out)
    return append(b, out...)
}
func (x *shakeAdapter) Reset()         { x.shake.Reset() }
func (x *shakeAdapter) Size() int      { return x.length }
func (x *shakeAdapter) BlockSize() int { return x.shake.BlockSize() }



var Registry = map[string]func() hash.Hash{
	"sha256": sha256.New,
	"sha384": sha512.New384,
	"sha512": sha512.New,
	"sha3-224": sha3.New224,
	"sha3-256": sha3.New256,
	"sha3-384": sha3.New384,
	"sha3-512": sha3.New512,
	"shake-128": func() hash.Hash {
		return &shakeAdapter{shake: sha3.NewShake128(), length: 32}
	},
	"shake-256": func() hash.Hash {
		return &shakeAdapter{shake: sha3.NewShake256(), length: 64}
	},
}

// Legacy hash not supported yet.
var RegistryLegacy = map[string]func() hash.Hash {
	"md5": md5.New,	      	  // insecure, legacy only
	"sha1": sha1.New,		  // insecure, legacy only
}

type HashKeyFactory func(key []byte) (hash.Hash, error)

var RegistryKey = map[string]HashKeyFactory {
	"blake2b-256": func(key []byte) (hash.Hash, error) {
		return blake2b.New256(key)
	},
	"blake2b-384": func(key []byte) (hash.Hash, error) {
		return blake2b.New384(key)
	},
	"blake2b-512": func(key []byte) (hash.Hash, error) {
		return blake2b.New512(key)
	},
	"blake2s-128": func(key []byte) (hash.Hash, error) {
		return blake2s.New128(key)
	},
	"blake2s-256": func(key []byte) (hash.Hash, error) {
		return blake2s.New256(key)
	},
}


const (
	KB = 1024
	MB = KB * 1024

	buf_32_kb  int = 32 * KB
	buf_64_kb  int = 64 * KB
	buf_1_mb   int = MB
	buf_5_mb   int = 5 * MB
	buf_10_mb  int = 10 * MB
)


func bufferSize(fileSize int64) int {
    switch {
    case fileSize < 10*MB:
        return 32 * KB
    case fileSize < 1024*MB:
        return 1 * MB
    default:
        return 5 * MB
    }
}

// hash a file with moist effective buffer size
func HashFile(hash, filePath string) (string, error) {
    fi, err := os.Stat(filePath)
    if err != nil {
        return "", fmt.Errorf("stat file: %w", err)
    }

    bufSize := bufferSize(fi.Size())

    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("open file: %w", err)
    }
    defer file.Close()

    factory, ok := Registry[strings.ToLower(hash)]
    if !ok {
        return "", fmt.Errorf("unsupported hash: %s", hash)
    }

    hasher := factory()
    buf := make([]byte, bufSize)

    for {
        n, err := file.Read(buf)
        if n > 0 {
            if _, werr := hasher.Write(buf[:n]); werr != nil {
                return "", fmt.Errorf("hash write: %w", werr)
            }
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return "", fmt.Errorf("read file: %w", err)
        }
    }

    return hex.EncodeToString(hasher.Sum(nil)), nil
}

func HashFileKey(hash string, key []byte, filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("couldn't open file: %w", err)
    }
    defer file.Close()

	factory, ok := RegistryKey[strings.ToLower(hash)]
	if !ok {
		return "", fmt.Errorf("unsupported hash: %s", hash)
	}

	hasher, err := factory(key)
	if err !=  nil {
		return "", fmt.Errorf("couldn't hash file: %w", err)
	}

	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("couldn't hash file: %w", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// hash a given byte string
func HashData(hash string, data []byte) (string, error) {
	factory, ok := Registry[strings.ToLower(hash)]
	if !ok {
		return "", fmt.Errorf("unsupported hash: %s", hash)
	}

	hasher := factory()
    if _, err := hasher.Write(data); err != nil {
        return "", fmt.Errorf("couldn't hash data: %w", err)
    }
    return hex.EncodeToString(hasher.Sum(nil)), nil
}

// compare two files
func CompareFiles(hash string, fileA string, fileB string) (bool, error) {
	hashA, errorA := HashFile(hash, fileA)
    if errorA != nil {
        return false, fmt.Errorf("hashing %s failed: %w", fileA, errorA)
    }

    hashB, errorB := HashFile(hash, fileB)
    if errorB != nil {
        return false, fmt.Errorf("hashing %s failed: %w", fileB, errorB)
    }

	return hashA == hashB, nil
}

func CompareFilesKey(hash string, keyA []byte, fileA string, keyB []byte, fileB string) (bool, error) {
	hashA, errorA := HashFileKey(hash, keyA, fileA)
    if errorA != nil {
        return false, fmt.Errorf("hashing %s failed: %w", fileA, errorA)
    }

    hashB, errorB := HashFileKey(hash, keyB, fileB)
    if errorB != nil {
        return false, fmt.Errorf("hashing %s failed: %w", fileB, errorB)
    }

	return hashA == hashB, nil
}
