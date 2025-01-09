# GoLang Blog Backend

## Setup

If new to Golang, follow [tutorial to setup and run GO provided on main doc page](https://go.dev/learn/). It's documentation is fine

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

## New dependencies

To use dependencies, as for example `google uuid` must update go.mod by using following command

```bash
go get github.com/google/uuid
go get github.com/joho/godotenv
```