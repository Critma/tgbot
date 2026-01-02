package main

import (
	"database/sql"
	"log"

	"github.com/critma/tgsheduler/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Config load error", err)
	}

	db, err := sql.Open("postgres", config.DB_URL)
	if err != nil {
		log.Panic("fail to open db", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Panic("fail to ping db", err)
	}

	log.Println("Connected to database")
}
