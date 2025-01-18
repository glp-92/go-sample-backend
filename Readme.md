# GoLang Blog Backend

## Setup

If new to Golang, follow [tutorial to setup and run GO provided on main doc page](https://go.dev/learn/). It's documentation is fine

Once installed, add `export PATH=$PATH:/usr/local/go/bin` in the end of `~/.profile` and reboot.

For a new project

```bash
mkdir projectName
cd projectName
go mod init projectName
```

Setup SQLite database

```bash
sudo apt install sqlite3
```

Create database

```bash
sqlite3 recordings.db
```

It will open DB bash, so create a table

```sql
CREATE TABLE posts (
    id CHAR(36) NOT NULL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    excerpt TEXT,
    content TEXT,
    date DATETIME NOT NULL
);
```

```sql
CREATE TABLE users (
    id CHAR(36) NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);
```

## New dependencies

To use dependencies, as for example `google uuid` must update go.mod by using following command

```bash
go get github.com/google/uuid
go get github.com/joho/godotenv
go get golang.org/x/crypto
go get github.com/go-sql-driver/mysql
go get -u github.com/golang-jwt/jwt/v5
```