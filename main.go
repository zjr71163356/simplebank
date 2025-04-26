package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/zjr71163356/simplebank/api"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/utils"
)

func main() {
	var err error
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config file:", err)
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	// fmt.Print(connManage)
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}
	store := db.NewStore(testDB)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create server:", err)
	}
	server.Start(config.Address)
}
