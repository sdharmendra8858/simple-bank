package api

import (
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccounts)
	server.router = router

	return server
}

func (Server *Server) Start(address string) error {
	return Server.router.Run(address)
}

func errorHandler(err error) gin.H {
	return gin.H{"error": err.Error()}
}
