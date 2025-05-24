package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/zjr71163356/simplebank/api"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	_ "github.com/zjr71163356/simplebank/doc/statik"
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
		log.Fatal().Err(err).Msg("can not load config file")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	// fmt.Print(connManage)
	if err != nil {
		log.Fatal().Err(err).Msg("can not connect to db")
	}

	runDBMigration(config.MigrateURL, config.DBSource)

	store := db.NewStore(conn)
	go runGRPCServer(config, store)
	runGRPCGatewayServer(config, store)
}
func runDBMigration(migrateURL string, dbSource string) {
	migration, err := migrate.New(migrateURL, dbSource)

	if err != nil {
		log.Fatal().Err(err).Msg("can not create migration")
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("can not run migration")
	}

	log.Info().Msg("db migration completed")
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("can not create server")
	}
	server.Start(config.HTTPServerAddress)
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal().Err(err).Msg("can not create server")
	}

	interceptor := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(interceptor)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("can not create listener")
	}

	log.Info().Msgf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("can not start grpc server")
	}

}
func runGRPCGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("can not create server")
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
		log.Fatal().Err(err).Msg("can not register handler server")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("can not create statik fs")
	}

	hfs := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", hfs)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	log.Info().Msgf("start gataway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("can not start http server")
	}
}
