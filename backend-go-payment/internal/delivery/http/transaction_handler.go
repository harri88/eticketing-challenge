package http

import (
	"net/http"
	"strconv"

	"github.com/harri88/eticketing-challenge/internal/usecase"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	UseCase *usecase.TransactionUsecase
}

func NewTransactionHandler(e *echo.Echo, uc *usecase.TransactionUsecase) {
	handler := &TransactionHandler{
		UseCase: uc,
	}
	e.GET("/api/v1/transactions", handler.GetAll)
	e.GET("/api/v1/transactions/:id", handler.GetByID)
	e.GET("/api/v1/transactions/txn/:transaction_id", handler.GetByTransactionID)
	e.GET("/api/v1/transactions/order/:order_id", handler.GetByOrderID)
}

// GetAll godoc
// @Summary Get all transactions
// @Description Retrieve all payment transactions with pagination support
// @Tags transactions
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) GetAll(c echo.Context) error {
	transactions, err := h.UseCase.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var data interface{} = transactions
	if transactions == nil {
		data = []interface{}{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

// GetByID godoc
// @Summary Get transaction by ID
// @Description Retrieve a specific transaction by its primary key ID
// @Tags transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid transaction ID"})
	}

	transaction, err := h.UseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if transaction == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "transaction not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": transaction,
	})
}

// GetByTransactionID godoc
// @Summary Get transaction by transaction ID
// @Description Retrieve a transaction by its unique transaction_id (e.g., TXN-CC-17000000)
// @Tags transactions
// @Produce json
// @Param transaction_id path string true "Transaction ID (e.g., TXN-CC-17000000)"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions/txn/{transaction_id} [get]
func (h *TransactionHandler) GetByTransactionID(c echo.Context) error {
	transactionID := c.Param("transaction_id")

	transaction, err := h.UseCase.GetByTransactionID(c.Request().Context(), transactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if transaction == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "transaction not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": transaction,
	})
}

// GetByOrderID godoc
// @Summary Get transaction by order ID
// @Description Retrieve a transaction by its associated order_id
// @Tags transactions
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions/order/{order_id} [get]
func (h *TransactionHandler) GetByOrderID(c echo.Context) error {
	orderID := c.Param("order_id")

	transaction, err := h.UseCase.GetByOrderID(c.Request().Context(), orderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if transaction == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "transaction not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": transaction,
	})
}
