// apps/backend/internal/domain/user/user.go
// نسخه کامل و نهایی - بدون نیاز به تغییر

package user

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// UserStatus وضعیت یک کاربر را مشخص می‌کند.
type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusSuspended UserStatus = "suspended"
	StatusDeleted   UserStatus = "deleted"
)

// User موجودیت اصلی کاربر در لایه Domain است.
type User struct {
	ID             string
	Email          string
	EmailVerified  bool
	PasswordHash   string
	Username       string
	FullName       string
	Status         UserStatus
	Roles          []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// NewUser یک کاربر جدید با مقادیر پیش‌فرض ایجاد می‌کند.
func NewUser(id, email, username, fullName string) *User {
	return &User{
		ID:        id,
		Email:     email,
		Username:  username,
		FullName:  fullName,
		Status:    StatusActive,
		Roles:     []string{"user"}, // نقش پیش‌فرض
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SetPassword رمز عبور کاربر را هش کرده و تنظیم می‌کند.
func (u *User) SetPassword(password string) error {
	// هزینه هش کردن را می‌توان قابل تنظیم کرد
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword بررسی می‌کند که آیا رمز عبور ارائه شده با هش ذخیره شده مطابقت دارد یا خیر.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}