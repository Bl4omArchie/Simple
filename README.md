# Simple package

Simple is a lightweight Go package designed for utilities tasks. Avoid boilerplate and insecure code and make it simple!

The package is currently in its early v1. The purpose of this package is to get simple functions with strong features like context management, easy configuration, algorithm efficiency or default security.
While I am developping, I am improving day by day the package. Finding new ideas and ways to make the features even more simpler. Don't mind to share you thought and advices.

# Features

- Database management :   
    - Choose your preferred driver: MySQL, PostgreSQL, or SQLite
    - Simple methods: OpenDatabase(), CloseDatabase(), and standard SQL queries.
    - Thread-safe operations with internal mutexes.
    - Use the provided Database struct or implement your own via the DatabaseManager interface.

- Web client management: :
    - Make web requests more simple 
    - Choose the client of your choice : http, Tor sock, your own proxy
    - Fetch raw web content with GetContent().
    - Parse HTML documents with GetParsedContent().
    - Download files with DownloadDocument(), and with automatic sha-256 hash computation via DownloadDocumentReturnHash().

- Hash functions :
    - Supported algorithms:
        - sha-224, sha-256, sha-384, sha-512
        - sha3 family (224, 256, 384, 512)
        - shake-128, shake-256
        - blake2b (256, 384, 512)
        - blake2s (128, 256)
        - md5 (legacy)
        - sha1 (legacy)
    - File hashing with optional buffer sizes for efficiency:
        - Constants: buf_32_kb, buf_64_kb, buf_1_mb, buf_5_mb, buf_10_mb
    - Compare files or hash arbitrary data.

- Environnement variables :
    - Use the gotenv package to pick up your env variables easily


# Development

## v0 to v1 :
- 🚩 **v0.1.0** : first commit
- 🪜 **v0.2.0** : more feature for database and major rework of requests with a clean client management
- 🪜 **v0.2.1** : add Hash feature
- 🪜 **v0.2.2** : add Env feature + add Test and Mock + small improvements
- ➡️ **v0.2.3** : rework of Database feature + small fixes and improvements
- 👷‍♂️ **v0.2.4** : add new Context as a new feature for easy context creation 
- 🏁 **v1.0.0** : first release, full revision of the code + comments

## v1 and further :

- Better error model : no panic + set of sentinel errors ...
- Security by default policy (i.e : for database, always setup the sslmode or support certificate for requets)
- Support ORM for database 
- Integration of workflows
- Logging Hooks


# Dependencies

- Database features are based on the Go’s standard [database/sql](https://pkg.go.dev/database/sql) package
- Requests features are based on Go’s standard [net/http](https://pkg.go.dev/net/http) package
- Hash features are based on Go’s standard [crypto](golang.org/x/crypto) package
- Env features are based on [gotenv](https://github.com/subosito/gotenv) package
