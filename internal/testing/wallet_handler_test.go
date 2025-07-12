package testing

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gunsu12/go-wallet/internal/delivery"
	"github.com/gunsu12/go-wallet/internal/domain"
	"github.com/stretchr/testify/assert"
)

// --- Mock Usecase Implementation ---
type mockWalletUsecase struct {
	CreateFn     func(wallet *domain.Wallet) error
	FindByUserFn func(userID string) ([]domain.Wallet, error)
}

func (m *mockWalletUsecase) Create(wallet *domain.Wallet) error {
	return m.CreateFn(wallet)
}

func (m *mockWalletUsecase) Update(wallet *domain.Wallet, id string) error {
	return nil
}

func (m *mockWalletUsecase) FindByID(id string) (*domain.Wallet, error) {
	return nil, nil
}

func (m *mockWalletUsecase) FindByUser(userID string) ([]domain.Wallet, error) {
	return m.FindByUserFn(userID)
}

// --- Setup Router Helper ---
func setupTestRouter(mockUC *mockWalletUsecase) *gin.Engine {
	r := gin.Default()
	handler := &delivery.WalletHandler{Usecase: mockUC}

	// Public
	r.POST("/wallets/add", handler.Create)

	// Simulasi Auth Middleware: inject user_id manual
	r.GET("/wallets/list", func(c *gin.Context) {
		c.Set("user_id", "mock-user-id")
		handler.List(c)
	})

	return r
}

// --- Test Create Wallet Success ---
func TestWalletHandler_Create_Success(t *testing.T) {
	mockUC := &mockWalletUsecase{
		CreateFn: func(wallet *domain.Wallet) error {
			wallet.ID = "mock-id"
			return nil
		},
	}

	router := setupTestRouter(mockUC)

	body := domain.Wallet{
		Name:        "Dompet A",
		Description: "Dompet Bulanan",
		Amount:      "100000",
		UserID:      "mock-user-id",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/wallets/add", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

// --- Test Get Wallet List Success ---
func TestWalletHandler_List_Success(t *testing.T) {
	mockUC := &mockWalletUsecase{
		FindByUserFn: func(userID string) ([]domain.Wallet, error) {
			return []domain.Wallet{
				{ID: "wallet1", Name: "Dompet A", Amount: "100000"},
				{ID: "wallet2", Name: "Dompet B", Amount: "50000"},
			}, nil
		},
	}

	router := setupTestRouter(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/wallets/list", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// --- Test Get Wallet List Failure ---
func TestWalletHandler_List_Error(t *testing.T) {
	mockUC := &mockWalletUsecase{
		FindByUserFn: func(userID string) ([]domain.Wallet, error) {
			return nil, errors.New("db error")
		},
	}

	router := setupTestRouter(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/wallets/list", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
