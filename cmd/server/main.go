package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"fullstackcms/backend/configs"
	"fullstackcms/backend/internal/router"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	dbCfg := mysql.Config{
		User:                 cfg.DB.Username,
		Passwd:               cfg.DB.Password,
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "recordings",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err = sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router.SetupRouter(db)
	log.Println("Servidor escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
