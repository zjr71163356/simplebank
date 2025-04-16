package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/zjr71163356/simplebank/db/sqlc"
)

type CreateTransfer struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	reqData := CreateTransfer{}
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.CreateTransfer(ctx, db.CreateTransferParams{
		FromAccountID: reqData.FromAccountID,
		ToAccountID:   reqData.ToAccountID,
		Amount:        reqData.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)

}
