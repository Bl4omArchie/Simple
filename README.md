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


# Development

The package is currently at version v0.1.0 as I continue to write tests and improve the codebase.
Once the tests are stable, I'll release the v1 of the package.

✅ v0.1.0 : current version

🪜 v0.2.0 : more feature for database and major rework of requests with a clean client management

➡️ v0.2.1 : for database feature, add optional parameters for more flexible configuration

🪜 v0.2.2 : add Context for database and requests

🪜 v0.3.0 : test for database and requests

🚩 v1.0.0 : first release, full revision of the code and comments above the functions
