package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/zjr71163356/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/CreateAccount", server.createAccount)
	router.GET("/GetAccount/:id", server.getAccount)
	router.GET("/GetAccountList", server.getAccountList)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	err := server.router.Run(address)
	if err != nil {
		return err
	}
	return err
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
