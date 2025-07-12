package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gunsu12/go-wallet/internal/domain"
	"github.com/gunsu12/go-wallet/internal/middleware"
	"github.com/gunsu12/go-wallet/internal/usecase"
)

type WalletHandler struct {
	usecase *usecase.WalletUsecase
}

func NewWalletHandler(r *gin.Engine, uc *usecase.WalletUsecase) {
	handler := &WalletHandler{usecase: uc}

	auth := r.Group("/wallets")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/list", handler.List)
		auth.GET("/detail/:id", handler.Detail)
		auth.POST("/add", handler.Create)
		auth.PUT("/update/:id", handler.Update)
		auth.DELETE("/delete/:id", handler.Delete)
	}
}

func (h *WalletHandler) List(c *gin.Context) {
	uid, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id tidak ditemukan silahkan login terlebih dahulu: "})
		return
	}

	uidStr, ok := uid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID invalid "})
		return
	}

	wallets, err := h.usecase.FindByUser(uidStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "dompet ditemukan",
		"data":    wallets,
	})
}

func (h *WalletHandler) Detail(c *gin.Context) {
	id := c.Param("id")

	wallet, err := h.usecase.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// validasi jika dompet memang dimiliki oleh user
	user_id, _ := c.Get("user_id")

	if wallet.UserID != user_id.(string) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tidak di ijinkan melihat dompet ini"})
		return
	}

	c.JSON(http.StatusFound, gin.H{
		"message": "data found",
		"data": gin.H{
			"id":          wallet.ID,
			"name":        wallet.Name,
			"description": wallet.Description,
			"amount":      wallet.Amount,
			"created_at":  wallet.CreatedAt,
			"updated_at":  wallet.UpdatedAt,
		},
	})
}

func (h *WalletHandler) Create(c *gin.Context) {
	var wallet domain.Wallet
	uid, exist := c.Get("user_id")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "tidak ditemukan user id"})
	}
	// konfersi body ke json yang dikonversi langsung ke struct domain
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validasi gagal",
			"error":   err.Error(),
		})
		return
	}

	// lakukan inject ke userID setelah bind json sehingga tidak terjadi penumpukan dengan user request
	wallet.UserID = uid.(string)

	// panggil usecase
	err := h.usecase.Create(&wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "something wrong ",
		})
	}

	// jika berhasil input
	c.JSON(http.StatusCreated, gin.H{
		"message": "wallet berhasil didaftarkan",
		"data": gin.H{
			"id":          wallet.ID,
			"name":        wallet.Name,
			"description": wallet.Description,
			"amount":      wallet.Amount,
			"user_id":     wallet.UserID,
		},
	})
}

func (h *WalletHandler) Update(c *gin.Context) {
	walletID := c.Param("id")

	// Ambil wallet lama
	existingWallet, err := h.usecase.FindByID(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingWallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet tidak ditemukan"})
		return
	}

	// Cek user yang login
	userID, _ := c.Get("user_id")
	if existingWallet.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tidak diizinkan mengubah dompet ini"})
		return
	}

	// Bind hanya field yang diizinkan
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
		return
	}

	// Siapkan data update
	existingWallet.Name = req.Name
	existingWallet.Description = req.Description

	// Panggil usecase untuk update
	if err := h.usecase.Update(existingWallet, walletID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Wallet berhasil diperbarui",
		"data":    existingWallet,
	})
}

func (h *WalletHandler) Delete(c *gin.Context) {
	walletID := c.Param("id")

	// Ambil wallet lama
	existingWallet, err := h.usecase.FindByID(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingWallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet tidak ditemukan"})
		return
	}

	// Cek user yang login
	userID := c.GetString("user_id")
	if existingWallet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tidak diizinkan menghapus dompet ini"})
		return
	}

	if err := h.usecase.Delete(walletID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "berhasil dihapus"})
}
