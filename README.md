# Simple package

Simple is a Go package designed to avoid repetitive operations such as opening a database connection or making web requests.
For now, it provides utility functions that help with the development of small scripts.

# Features

- Database wrappers :
    - sqlite, postgresql, clickhouse, gaussDB and mysql are supported
    - easily open, close, and modify the DSN of your database
    - insert data and run migrations

- Web client management: :
    - pick your http client or tor socks proxy easily
    - fetch web content and parse HTML documents
    - download documents with hash computation

- Hash functions :
    - functions supported : sha224, sha256, sha384, sha512, sha3 family and md5
    - buffer feature that allow hashing a file, block by block for better efficiency
    - constant variable for easy buffer size picking


# Development

The package is currently at version v0.1.0 as I continue to write tests and improve the codebase.
Once the tests are stable, I'll release the v1 of the package.

- 🚩 **v0.1.0** : first commit
- 🪜 **v0.2.0** : more feature for database and major rework of requests with a clean client management
- 🪜 **v0.2.1** : add Hash feature
- ➡️ **v0.2.2** : new feature OpenEnv() + improve registry for Hash featuure + add Test for hash and env + add Mock for requests
- 👷‍♂️ **v0.2.3** : add optional parameters for more flexible DSN, add support for advanced driver configuration and existing database connection.
- 🪜 **v0.2.4** : add Context for database and requests
- 🏁 **v1.0.0** : first release, full revision of the code and comments above the functions


# Dependencies

- Database features are based on the [GORM](https://gorm.io) library  
- Requests features are based on Go’s standard [net/http](https://pkg.go.dev/net/http) package
- Hash features are based on Go’s standard [crypto](golang.org/x/crypto) package
