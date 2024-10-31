package api

import (
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     utils.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupServerRoutes()
	return server, nil
}

func (Server *Server) Start(address string) error {
	return Server.router.Run(address)
}

func errorHandler(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (Server *Server) setupServerRoutes() {
	router := gin.Default()
	router.POST("/users", Server.CreateUser)
	router.POST("/users/login", Server.LoginUser)
	router.POST("/tokens/renew-access", Server.renewToken)

	routerGroup := router.Group("/", authMiddleware(Server.tokenMaker))

	routerGroup.POST("/accounts", Server.CreateAccount)
	routerGroup.GET("/accounts/:id", Server.GetAccount)
	routerGroup.GET("/accounts", Server.ListAccounts)

	routerGroup.POST("/transfer", Server.CreateTransfer)
	Server.router = router
}
