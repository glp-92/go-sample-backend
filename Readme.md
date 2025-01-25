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

### SQLite

```bash
sudo apt install sqlite3
```
Create database

```bash
sqlite3 recordings.db
sqlite3 .read db_entrypoint.sql
```

### MySQL

Install MySQL on system (Docker works fine too)

```bash
sudo apt install mysql-server -y
sudo systemctl enable mysql.service
```

Change root usr
```bash
sudo mysql
```

```sql
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'rootpwd';
CREATE USER 'customuser'@'localhost' IDENTIFIED BY 'password'
CREATE DATABASE db_name;'
GRANT ALL PRIVILEGES ON db_name.* TO 'customuser'@'localhost';
```

```bash
mysql -u customuser -p password db_name < db_entrypoint.sql
```

## New dependencies

To use dependencies, as for example `google uuid` must update go.mod by using following command

```bash
go get -u github.com/google/uuid
go get -u github.com/joho/godotenv
go get -u golang.org/x/crypto
go get -u github.com/go-sql-driver/mysql
go get -u github.com/golang-jwt/jwt/v5
```