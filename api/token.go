package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewToken(ctx *gin.Context) {
	var req renewTokenRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorHandler(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	if session.IsBlocked {
		err := errors.New("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorHandler(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := errors.New("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorHandler(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := errors.New("invalid token")
		ctx.JSON(http.StatusUnauthorized, errorHandler(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := errors.New("token expired")
		ctx.JSON(http.StatusUnauthorized, errorHandler(err))
		return
	}

	token, accessTokenPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.ExpiryTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	res := renewTokenRes{
		AccessToken:          token,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, res)
}
