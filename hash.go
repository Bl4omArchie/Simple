package simple

import (
	"io"
	"os"
	"fmt"
	"hash"
	"strings"
	"encoding/hex"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/blake2b"
)

// Warning : working on a generic interface for every hash

type ShaFactory func() hash.Hash
type BlakeFactory func(key []byte) (hash.Hash, error)
type ShakeFactory func() sha3.ShakeHash


var registrySha = map[string]ShaFactory{
	"md5":       	 md5.New,	      // insecure, legacy only
	"sha1":      	 sha1.New,		  // insecure, legacy only
	"sha224":    	 sha256.New224,
	"sha256":    	 sha256.New,
	"sha384":    	 sha512.New384,
	"sha512":    	 sha512.New,
	"sha3-224":  	 sha3.New224,
	"sha3-256":  	 sha3.New256,
	"sha3-384":  	 sha3.New384,
	"sha3-512":  	 sha3.New512,
}

var registryBlake = map[string]BlakeFactory {
	"blake2b-256":   blake2b.New256,
	"blake2b-384":   blake2b.New384,
	"blake2b-512":   blake2b.New512,
}

var registryShake = map[string]ShakeFactory {
	"shake-128":   sha3.NewShake128,
	"shake-256":   sha3.NewShake256,
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

// Warning : currently, it only works for shaRegistry

func HashFile(hash string, filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("couldn't open file: %w", err)
    }
    defer file.Close()

	factory, ok := registrySha[strings.ToLower(hash)]
	if !ok {
		return "", fmt.Errorf("unsupported hash: %s", hash)
	}

	hasher := factory()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("couldn't hash file: %w", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func HashFileBuffer(hash string, filePath string, bufSize int) (string, error) {
	file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("couldn't open file: %w", err)
    }
    defer file.Close()

	factory, ok := registrySha[strings.ToLower(hash)]
	if !ok {
		return "", fmt.Errorf("unsupported hash: %s", hash)
	}

	hasher := factory()
	buf := make([]byte, bufSize)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			if _, werr := hasher.Write(buf[:n]); werr != nil {
				return "", fmt.Errorf("error while hashing : %w", werr)
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

func HashData(hash string, data []byte) (string, error) {
	factory, ok := registrySha[strings.ToLower(hash)]
	if !ok {
		return "", fmt.Errorf("unsupported hash: %s", hash)
	}

	hasher := factory()
    if _, err := hasher.Write(data); err != nil {
        return "", fmt.Errorf("couldn't hash data: %w", err)
    }
    return hex.EncodeToString(hasher.Sum(nil)), nil
}

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
