package api

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/mfitrahrmd/simple_bank/database/sqlc"
	"github.com/pkg/errors"
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

type getAccountRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) GetAccount(c *gin.Context) {
	var req getAccountRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	foundAccount, err := s.store.GetAccount(context.Background(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, foundAccount)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) ListAccounts(c *gin.Context) {
	var req listAccountsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	foundAccount, err := s.store.ListAccounts(context.Background(), db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, foundAccount)
}
