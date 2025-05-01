package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/token"
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
	payload, ok := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid token payload")))
	}
	isVaild, err := server.validAccount(ctx, reqData.FromAccountID, reqData.Currency, payload.Username)
	if !isVaild {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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
func (server *Server) validAccount(ctx *gin.Context, accountId int64, curreny string, username string) (bool, error) {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("account %d does not exist", accountId)
		}
		return false, fmt.Errorf("failed to get account %d: %w", accountId, err)
	}

	if account.Owner != username {
		return false, fmt.Errorf("account %d does not belong to user %s", accountId, username)
	}

	if account.Currency != curreny {
		return false, fmt.Errorf("account %d has a different currency %s", accountId, account.Currency)
	}
	return true, nil

}
