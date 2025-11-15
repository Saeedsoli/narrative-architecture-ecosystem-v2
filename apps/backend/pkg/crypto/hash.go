// apps/backend/pkg/crypto/hash.go

package crypto

import "golang.org/x/crypto/bcrypt"

// HashPassword یک رمز عبور را با استفاده از الگوریتم bcrypt هش می‌کند.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash بررسی می‌کند که آیا رمز عبور ارائه شده با هش ذخیره شده مطابقت دارد یا خیر.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}