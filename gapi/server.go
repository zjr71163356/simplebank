package gapi

import (
	"fmt"

	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/pb"
	"github.com/zjr71163356/simplebank/token"
	"github.com/zjr71163356/simplebank/utils"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	maker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("fail to create token:%w", err)
	}
	server := &Server{config: config, store: store, tokenMaker: maker}
	
	return server, nil
}
