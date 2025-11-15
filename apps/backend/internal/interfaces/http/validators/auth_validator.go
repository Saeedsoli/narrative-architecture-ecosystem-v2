// apps/backend/internal/interfaces/http/validator/auth_validator.go

package validator

import (
	"regexp"
	"github.com/go-playground/validator/v10"
)

// PasswordPolicy یک تابع اعتبارسنجی سفارشی برای رمز عبور است.
// این تابع بررسی می‌کند که آیا رمز عبور شامل حداقل یک حرف بزرگ، یک حرف کوچک، یک عدد، و یک کاراکتر خاص است.
var PasswordPolicy validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSymbol := regexp.MustCompile(`[\W_]`).MatchString(password) // Non-alphanumeric

	return hasUpper && hasLower && hasNumber && hasSymbol
}