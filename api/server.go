package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/zjr71163356/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", vaildatorCurrency)
	}
	router.POST("/CreateAccount", server.createAccount)
	router.GET("/GetAccount/:id", server.getAccount)
	router.GET("/GetAccountList", server.getAccountList)

	router.POST("/CreateTransfer", server.createTransfer)
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
