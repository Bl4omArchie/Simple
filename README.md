# Simple package

Simple is a Go package that aims to avoid repetitive operations such as opening a database or web requests.
For the moment it is about simple functions aimed at helping the development of small scripts.

# Features

- database :
    - open database with your favourite driver (sqlite, postgresql, mysql)
    - close database
    - migrate tables
    - insert rows

- requests : 
    - get content of a website
    - get content of a .onion website
    - get parsed content of a website
    - download a document and retun its hash value

# Development

Currently, I am writing clean tests for my package.
Once the tests released I could safely claim the v1 of the package.

TODO :
- make tests
- add OpenMysql
- add OpenPostgres
- add OpenWithEnv
