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
	config, _ := utils.LoadConfig(".")

	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	// fmt.Print(connManage)
	if err != nil {
		log.Fatal(err)
	}
	store := db.NewStore(testDB)
	server := api.NewServer(store)
	server.Start(config.Address)
}
