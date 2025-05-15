package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/zjr71163356/simplebank/api"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/gapi"
	"github.com/zjr71163356/simplebank/pb"
	"github.com/zjr71163356/simplebank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGRPCServer(config, store)
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create server:", err)
	}
	server.Start(config.HTTPServerAddress)
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("can not create server:", err)
	}
	fmt.Println(config.GRPCServerAddress)
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("can not create listener:", err)
	}

	log.Printf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can not start grpc server:", err)
	}

}
