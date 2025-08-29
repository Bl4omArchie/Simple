# Changelog

**29/08/25**
- add hash feature : HashFile(), HashFileBuffer(), HashData(), CompareFile()
- buffer size constant : buf_32_kb, buf_64_kb, buf_1_mb, buf_5_mb and buf_10_mb
- working on generic interface for more hash functions. Only sha registry available.

**23/08/2025**
- add : OpenMysql(), OpenMysqlUnixSocket(), OpenPostgres()
- start working on OpenEnv()
- get your client more easily with httpClient() and onionClient()
- with the new client system, there is now only FetchBody(), GetContent() and GetParsedContent(). Now youu have to specify the client
- update DownloadDocumentReturnHash() with multiwriting
- Fix made on Mysql and Postgresql DSN
- add Clickhouse and GaussDB

(**21/08/2025**) First commit
- new feature database
- new feature requests 
