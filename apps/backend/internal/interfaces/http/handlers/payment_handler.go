// apps/backend/internal/interfaces/http/handlers/payment_handler.go

package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	app_commerce "narrative-architecture/apps/backend/internal/application/commerce"
)

type PaymentHandler struct {
	processPaymentUC *app_commerce.ProcessPaymentUseCase
	verifyPurchaseUC *app_commerce.VerifyPurchaseUseCase
}

func NewPaymentHandler(processUC *app_commerce.ProcessPaymentUseCase, verifyUC *app_commerce.VerifyPurchaseUseCase) *PaymentHandler {
	return &PaymentHandler{
		processPaymentUC: processUC,
		verifyPurchaseUC: verifyUC,
	}
}

// CreatePayment یک درخواست پرداخت جدید ایجاد کرده و کاربر را به درگاه هدایت می‌کند.
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req struct {
		OrderID string `json:"orderId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: orderId is required"})
		return
	}

	userID := c.GetString("userID")
	
	useCaseReq := app_commerce.ProcessPaymentRequest{
		OrderID: req.OrderID,
		UserID:  userID,
	}

	res, err := h.processPaymentUC.Execute(c.Request.Context(), useCaseReq)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"paymentUrl": res.PaymentURL})
}

// VerifyPaymentCallback پس از بازگشت کاربر از درگاه، پرداخت را تأیید می‌کند.
func (h *PaymentHandler) VerifyPaymentCallback(c *gin.Context) {
	orderID := c.Query("orderId")
	authority := c.Query("Authority")
	status := c.Query("Status")

	if orderID == "" || authority == "" || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback parameters"})
		return
	}

	req := app_commerce.VerifyPurchaseRequest{
		OrderID:   orderID,
		Authority: authority,
		Status:    status,
	}

	if err := h.verifyPurchaseUC.Execute(c.Request.Context(), req); err != nil {
		// هدایت به صفحه شکست در پرداخت
		c.Redirect(http.StatusTemporaryRedirect, "/payment/failed?error=verification_failed")
		return
	}

	// هدایت به صفحه موفقیت در پرداخت
	c.Redirect(http.StatusTemporaryRedirect, "/payment/success?orderId=" + orderID)
}