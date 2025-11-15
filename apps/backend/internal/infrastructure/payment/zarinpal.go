// apps/backend/internal/infrastructure/payment/zarinpal.go

package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"narrative-architecture/apps/backend/internal/domain/commerce"
)

const (
	zarinpalRequestURL  = "https://api.zarinpal.com/pg/v4/payment/request.json"
	zarinpalVerifyURL   = "https://api.zarinpal.com/pg/v4/payment/verify.json"
	zarinpalSandboxRequestURL = "https://sandbox.zarinpal.com/pg/v4/payment/request.json"
	zarinpalSandboxVerifyURL  = "https://sandbox.zarinpal.com/pg/v4/payment/verify.json"
)

type ZarinpalGateway struct {
	merchantID string
	isSandbox  bool
	httpClient *http.Client
}

func NewZarinpalGateway(merchantID string, isSandbox bool) *ZarinpalGateway {
	return &ZarinpalGateway{
		merchantID: merchantID,
		isSandbox:  isSandbox,
		httpClient: &http.Client{},
	}
}

type zarinpalRequest struct {
	MerchantID  string `json:"merchant_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	CallbackURL string `json:"callback_url"`
}

type zarinpalRequestResponse struct {
	Data struct {
		Authority string `json:"authority"`
		Fee       int    `json:"fee"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (z *ZarinpalGateway) CreatePaymentRequest(ctx context.Context, amount int, currency, description, callbackURL string) (*commerce.PaymentRequestResponse, error) {
	// زرین‌پال فقط از ریال/تومان پشتیبانی می‌کند
	// قیمت‌ها به تومان است، پس همان را ارسال می‌کنیم
	
	payload := zarinpalRequest{
		MerchantID:  z.merchantID,
		Amount:      amount / 10, // تبدیل از سنت (ریال) به تومان
		Description: description,
		CallbackURL: callbackURL,
	}

	body, _ := json.Marshal(payload)
	requestURL := zarinpalRequestURL
	if z.isSandbox {
		requestURL = zarinpalSandboxRequestURL
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.httpClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()

	var zarinpalResp zarinpalRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&zarinpalResp); err != nil { return nil, err }

	if len(zarinpalResp.Errors) > 0 {
		return nil, errors.New("zarinpal request failed")
	}

	paymentURL := "https://www.zarinpal.com/pg/StartPay/" + zarinpalResp.Data.Authority
	if z.isSandbox {
		paymentURL = "https://sandbox.zarinpal.com/pg/StartPay/" + zarinpalResp.Data.Authority
	}

	return &commerce.PaymentRequestResponse{
		PaymentID:  zarinpalResp.Data.Authority, // از Authority به‌عنوان شناسه موقت استفاده می‌کنیم
		PaymentURL: paymentURL,
		Authority:  zarinpalResp.Data.Authority,
	}, nil
}

type zarinpalVerifyRequest struct {
	MerchantID string `json:"merchant_id"`
	Amount     int    `json:"amount"`
	Authority  string `json:"authority"`
}

type zarinpalVerifyResponse struct {
	Data struct {
		Code   int    `json:"code"`
		RefID  int64  `json:"ref_id"`
		CardPan string `json:"card_pan"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (z *ZarinpalGateway) VerifyPayment(ctx context.Context, amount int, authority string) (*commerce.VerificationResponse, error) {
	payload := zarinpalVerifyRequest{
		MerchantID: z.merchantID,
		Amount:     amount / 10,
		Authority:  authority,
	}
	// ... (منطق ارسال درخواست POST به Verify URL)

	// پس از دریافت پاسخ موفق
	// var verifyResp zarinpalVerifyResponse
	// ...
	
	// if verifyResp.Data.Code == 100 || verifyResp.Data.Code == 101 {
	//	return &commerce.VerificationResponse{RefID: strconv.FormatInt(verifyResp.Data.RefID, 10)}, nil
	// }

	return nil, errors.New("payment verification failed")
}