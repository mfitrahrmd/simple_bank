package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mfitrahrmd/simple_bank/database/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	s := Server{
		store:  store,
		router: gin.Default(),
	}

	s.router.POST("/accounts", s.CreateAccount)
	s.router.GET("/accounts/:id", s.GetAccount)
	s.router.GET("/accounts", s.ListAccounts)

	return &s
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
