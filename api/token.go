package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type NewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) RenewToken(ctx *gin.Context) {
	reqData := RenewTokenRequest{}
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(reqData.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, payload.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.RefreshToken != reqData.RefreshToken {
		err = fmt.Errorf("session refresh token does not match request")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != payload.Username {
		err = fmt.Errorf("session user does not match token user ")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err = fmt.Errorf("session is blocked")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err = fmt.Errorf("session is expired")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	accessToken, payload, err := server.tokenMaker.CreateToken(payload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := NewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: payload.ExpiredAt,
	}
	ctx.JSON(http.StatusAccepted, rsp)

}
