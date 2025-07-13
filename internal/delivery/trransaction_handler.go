package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gunsu12/go-wallet/internal/domain"
	"github.com/gunsu12/go-wallet/internal/middleware"
	"github.com/gunsu12/go-wallet/internal/usecase"
)

type TransactionHandler struct {
	usecase *usecase.TransactionUsecase
}

func NewTransactionHandler(r *gin.Engine, uc *usecase.TransactionUsecase) {
	handler := &TransactionHandler{usecase: uc}

	auth := r.Group("/transactions")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/by_wallet/:id", handler.ListByWallet)
		auth.GET("/by_user/:id", handler.ListByUser)
		auth.POST("/create", handler.Create)
	}
}

func (h *TransactionHandler) ListByWallet(c *gin.Context) {
	waletId := c.Param("id")

	// get transaction from db
	trs, err := h.usecase.FindByWallet(waletId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if trs == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found on this wallet"})
		return
	}

	c.JSON(http.StatusFound, gin.H{
		"message": "Data found",
		"data":    trs,
	})
}

func (h *TransactionHandler) ListByUser(c *gin.Context) {
	userId := c.Param("id")

	// get transaction from db
	trs, err := h.usecase.FindByUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if trs == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found on this wallet"})
		return
	}

	c.JSON(http.StatusFound, gin.H{
		"message": "Data found",
		"data":    trs,
	})
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var transaction domain.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validasi gagal",
			"error":   err.Error(),
		})
		return
	}

	err := h.usecase.Create(&transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "gagal menyimpan",
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "transaksi berhasil sibuat",
		"data": gin.H{
			"id":                      transaction.ID,
			"transaction_description": transaction.TransactionDescription,
			"transaction_amount":      transaction.TransactionAmount,
			"transaction_type":        transaction.TransactionType,
			"wallet_id":               transaction.WalletID,
		},
	})
}
