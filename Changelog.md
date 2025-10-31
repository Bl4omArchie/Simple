# Changelog

**27/10/2025**
- update GetTable() for Orm
- add GetColumn() and OpenDatabase() for Orm
- UpdateTable() -> UpdateColumnWhereValue()

**27/10/2025**
- update Database to Orm
- new folder for tests
- update for Hash

**10/09/2025**
- add new feature, Deserializer.
- support the following file format : json, yaml, toml and xml. 

**02/09/2025**
- add DatabaseManager interface for utilities task such as connection, sql requests
- add easy database driver selection for mysql, postgresql and sqlite.
- add DownloadDocument()
- add new registry for legacy hash functions
- update roadmap in the README

**31/09/2025**
- add Registry and RegistryKey for hashl
- add more hash functions : blake2b-256, blake2b-384, blake2b-512, blake2s-128, blake2s-256, shake-128 and shake-256
- add Test for hash and env
- add Mock for requests
- add MIT license

**29/08/2025**
- add hash feature : HashFile(), HashFileBuffer(), HashData(), CompareFile()gvfd
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
