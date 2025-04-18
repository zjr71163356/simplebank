package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/zjr71163356/simplebank/db/sqlc"
)

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR "`
}

type GetAccountParams struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

type GetAccountListParams struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10" `
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
		//foreign_key_violation
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

func (server *Server) getAccount(ctx *gin.Context) {
	reqData := GetAccountParams{}
	if err := ctx.ShouldBindUri(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, reqData.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

func (server *Server) getAccountList(ctx *gin.Context) {
	reqData := GetAccountListParams{}
	if err := ctx.ShouldBindQuery(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var accountList []db.Account
	accountList, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  reqData.PageSize,
		Offset: (reqData.PageId - 1) * reqData.PageSize,
		Owner:  "ofudn",
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accountList)

}
