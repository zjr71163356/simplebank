package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	_ "github.com/zjr71163356/simplebank/doc/statik"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/zjr71163356/simplebank/api"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/gapi"
	"github.com/zjr71163356/simplebank/pb"
	"github.com/zjr71163356/simplebank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGRPCServer(config, store)
	runGRPCGatewayServer(config, store)
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
func runGRPCGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create server:", err)
	}

	opts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(opts)
	ctx := context.Background()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("can not register handler server:", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("can not create statik fs:", err)
	}

	hfs := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", hfs)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	log.Printf("start gataway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("can not start http server:", err)
	}
}
