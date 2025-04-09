package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/zjr71163356/simplebank/db/sqlc"
)

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR "`
}

func (server *Server) createAccount(ctx *gin.Context) {
	reqData := CreateAccountParams{}
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    reqData.Owner,
		Currency: reqData.Currency,
		Balance:  0,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, account)

}
