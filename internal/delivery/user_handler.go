package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gunsu12/go-wallet/internal/domain"
	"github.com/gunsu12/go-wallet/internal/middleware"
	"github.com/gunsu12/go-wallet/internal/usecase"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// NewUserHandler mendaftarkan endpoint user ke router
func NewUserHandler(r *gin.Engine, uc *usecase.UserUsecase) {
	handler := &UserHandler{usecase: uc}

	r.POST("/users/register", handler.Register)
	r.POST("/users/login", handler.Login) // 👈 ini

	auth := r.Group("/users")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/detail", handler.Detail)
	}
}

// Register menangani endpoint POST /users/register
func (h *UserHandler) Register(c *gin.Context) {
	var user domain.User

	//Bind JSON dari body ke struct user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	// Panggil usecase untuk registrasi user
	err := h.usecase.Register(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Berhasil dibuat
	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil didaftarkan",
		"user": gin.H{
			"id":           user.ID,
			"name":         user.Name,
			"username":     user.Username,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest

	// Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data login tidak valid: " + err.Error()})
		return
	}

	// Coba login via usecase
	user, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 👉 Simpel dulu: belum pakai JWT
	// Bisa ditambah token nanti
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"token":    user.Token,
		},
	})
}

func (h *UserHandler) Detail(c *gin.Context) {
	uid, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id tidak ditemukan silahkan login terlebih dahulu: "})
		return
	}

	// konversi lagi ini ya ke string?
	// karena return dari c.Get("user_id") adalah interface jadi harus di type assertion lagi
	uidStr, ok := uid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID invalid "})
		return
	}

	// panggil usecase
	user, err := h.usecase.GetByID(uidStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berikut detil account",
		"data": gin.H{
			"id":           user.ID,
			"name":         user.Name,
			"username":     user.Username,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"created_at":   user.CreatedAt,
		},
	})
}
