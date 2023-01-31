package api

import (
	"context"
	"github.com/gin-gonic/gin"
	db "github.com/mfitrahrmd/simple_bank/database/sqlc"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner,omitempty" binding:"required"`
	Currency string `json:"currency,omitempty" binding:"required,oneof=USD IDR"`
}

func (s *Server) CreateAccount(c *gin.Context) {
	var req createAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	createdAccount, err := s.store.CreateAccount(context.Background(), db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdAccount)
}
