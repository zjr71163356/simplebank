package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/token"
	"github.com/zjr71163356/simplebank/utils"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	maker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("fail to create token:%w", err)
	}
	server := &Server{config: config, store: store, tokenMaker: maker}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", vaildatorCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	authRouter := router.Group("/").Use(authMiddleWare(server.tokenMaker))
	authRouter.POST("/CreateAccount", server.createAccount)
	authRouter.GET("/GetAccount/:id", server.getAccount)
	authRouter.GET("/GetAccountList", server.getAccountList)

	router.POST("/CreateTransfer", server.createTransfer)
	router.POST("/User/Create", server.createUser)
	router.POST("/User/Login", server.LoginUser)
	server.router = router
}

func (server *Server) Start(address string) error {
	err := server.router.Run(address)
	if err != nil {
		return err
	}
	return err
}

// errorResponse formats the error message to be returned in the response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
