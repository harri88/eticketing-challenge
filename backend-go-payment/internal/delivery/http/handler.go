package http

import (
	"net/http"

	"github.com/harri88/eticketing-challenge/internal/domain"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type PaymentHandler struct {
	UseCase domain.PaymentUsecase
}

func NewPaymentHandler(e *echo.Echo, us domain.PaymentUsecase) {
	handler := &PaymentHandler{
		UseCase: us,
	}
	e.POST("/api/v1/payments", handler.MakePayment)

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

// MakePayment godoc
// @Summary Make Payment
// @Description Process a payment transaction for an order
// @Tags payments
// @Accept json
// @Produce json
// @Param request body domain.PaymentRequest true "Payment Request"
// @Success 200 {object} domain.Transaction
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/payments [post]
func (h *PaymentHandler) MakePayment(c echo.Context) error {
	var req domain.PaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Basic Validation (in production use go-playground/validator)
	if req.OrderID == "" { //|| req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	// Get Idempotency Key from Header
	idempotencyKey := c.Request().Header.Get("Idempotency-Key")

	tx, err := h.UseCase.ProcessPayment(c.Request().Context(), req, idempotencyKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tx)
}
