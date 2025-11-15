// apps/backend/internal/application/auth/register_user.go
// نسخه کامل و نهایی - بدون نیاز به تغییر

package auth

import (
	"context"
	"errors"
	"time"

	"narrative-architecture/apps/backend/internal/domain/user"
	"narrative-architecture/apps/backend/pkg/jwt"
	"github.com/oklog/ulid/v2"
)

// RegisterUserRequest ساختار درخواست برای ثبت‌نام کاربر است.
type RegisterUserRequest struct {
	Email    string
	Password string
	Username string
	FullName string
}

// RegisterUserResponse ساختار پاسخ پس از ثبت‌نام موفق است.
type RegisterUserResponse struct {
	User         *user.User
	AccessToken  string
	RefreshToken string
}

// UserRepository رابطی برای دسترسی به داده‌های کاربران است.
type UserRepository interface {
    FindByEmail(ctx context.Context, email string) (*user.User, error)
    Save(ctx context.Context, u *user.User) error
}

// RegisterUserUseCase منطق تجاری برای ثبت‌نام کاربر را کپسوله می‌کند.
type RegisterUserUseCase struct {
	userRepo UserRepository
	// eventBus EventBus // برای ارسال ایمیل خوش‌آمدگویی
}

// NewRegisterUserUseCase یک نمونه جدید از RegisterUserUseCase ایجاد می‌کند.
func NewRegisterUserUseCase(userRepo UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepo: userRepo}
}

// Execute متد اصلی برای اجرای Use Case است.
func (uc *RegisterUserUseCase) Execute(ctx context.Context, req RegisterUserRequest) (*RegisterUserResponse, error) {
	// 1. بررسی تکراری بودن ایمیل
	if existing, _ := uc.userRepo.FindByEmail(ctx, req.Email); existing != nil {
		return nil, errors.New("email already exists")
	}
	// TODO: بررسی تکراری بودن نام کاربری

	// 2. ایجاد موجودیت کاربر جدید
	newUser := user.NewUser(
		ulid.New().String(),
		req.Email,
		req.Username,
		req.FullName,
	)
	if err := newUser.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// 3. ذخیره کاربر در دیتابیس
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}
	
	// TODO: اختصاص نقش پیش‌فرض 'user'

	// 4. تولید توکن‌ها
	// TODO: خواندن secrets و TTLs از پیکربندی
	accessToken, refreshToken, err := jwt.GenerateTokens(newUser.ID, newUser.Email, newUser.Roles, "secret", "refresh_secret", 15*time.Minute, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}
	
	// TODO: ذخیره Refresh Token در جدول auth_tokens برای امنیت بیشتر

	return &RegisterUserResponse{
		User:         newUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}