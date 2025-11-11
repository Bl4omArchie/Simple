# Simple package

Simple is a Go package providing utilities tasks in order to reduces boilerplate code.
Each file represent a feature 


# Features

## ORM for database management

- Choose your preferred driver: MySQL, PostgreSQL, or SQLite
- Simple methods: Migrate(), GetRowBy(), GetRows(), UpdateRowBy(), DeleteRowBy() and CountRows()

Example :
```go
var ctx context.Context

db, err := simple.OpenDatabase(GetSqlite("path/to/my/database.db"))
if err != nil {
    fmt.Println(err)
}

book, err := simple.GetRowBy[&Book](ctx, db, "title", "The Go Programming Language Phrasebook")
if err != nil {
    fmt.Println(err)
} else {
    fmt.Println(book.Title)
}
```

## Web client management: :
- Choose the client of your choice : http, Tor sock, your own proxy...
- Fetch raw or pased HTML content
- Download files with DownloadDocument(), and with automatic sha-256 hash computation via DownloadDocumentReturnHash().

Example :
```go
var ctx context.Context
body, err := simple.GetContent(ctx, "https://golangdocs.com", simple.HttpClient(), nil)
if err != nil {
    fmt.Println(err)
}

hash, err := simple.DownloadDocumentReturnHash(ctx, "https://golangdocs.com/how-to-install-go-on-a-vps-server", "storage/file.html", simple.HttpClient(), nil)
if err != nil {
    fmt.Println(err)
} else {
    fmt.Println(hash)
}

//Custom request
body, err := simple.GetContent(ctx, "https://example.com/api", simple.HttpClient(), func(req *http.Request) error {
	req.Header.Set("User-Agent", "Go-http-client/1.1")
	req.Header.Set("Accept", "*/*")
    return nil
})
if err != nil {
    fmt.Println(err)
}
```

## Hash functions
- Supported algorithms:
    - sha-224, sha-256, sha-384, sha-512 and sha3 family
    - shake (128, 256)
    - blake2b (256, 384, 512) and blake2s (128, 256)
    - legacy : md5 and sha1
- File hashing
- Compare files or hash arbitrary data.

Example :
```go
sha_hash, err := simple.HashFile("sha3-512", "storage/file.html")
if err != nil {
    fmt.Println(err)
} else {
    fmt.Println(sha_hash)
}

blake_hash, err := simple.HashFileKey("blake2b-384", "mykey", "storage/file.html")
if err != nil {
    fmt.Println(err)
} else {
    fmt.Println(blake_hash)
}
```

## File

- Deserialization : 
    - one function for multiple file format support (json, yaml, toml and xml)
    - Use limit parameter to deserialize only a specific amount of elements
    - Set validation to true in order to apply tag validation from validator package

Example :
```go
data_json, err := simple.LoadFile[DataJson]("test.json", 0, false)
if err != nil {
    fmt.Println(err)
} else {
    fmt.Println(data_json)
}

data_yaml, err = simple.LoadFile[DataYaml]("test2.yaml", 0, true)
if err != nil {
    fmt.Println(err)
} else {
    fmt.Println(data_json)
}

err = simple.Unzip()

```

# Development

## v0 :
- add first features : Orm, Hash, Requests, File, Env 
- add **tests/** folder for mock and unit tests
- add Readme, MIT License, Changelog
- Code correction + comments + cleaning

## v1 :
- add context for Orm and Requests
- New feature on File : Unzip()
- Improved ORM functions

# Dependencies

- Orm features are based on [gorm](https://pkg.go.dev/gorm.io/gorm@v1.31.0) package
- Requests features are based on Go’s standard [net/http](https://pkg.go.dev/net/http) package
- Hash features are based on Go’s standard [crypto](golang.org/x/crypto) package
- Env features are based on [gotenv](https://github.com/subosito/gotenv) package
- Deserialize features are based on 
    - Go’s standard [encoding](https://pkg.go.dev/encoding/xml) package
    - [go-toml](github.com/pelletier/go-toml) package
    - [yaml.v3](gopkg.in/yaml.v3) package
    - [validator](github.com/go-playground/validator/v10) package
