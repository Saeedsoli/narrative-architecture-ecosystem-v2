// apps/backend/internal/interfaces/http/dto/auth_dto.go

package dto

import "narrative-architecture/apps/backend/internal/domain/user"

// UserResponse ساختار پاسخ API برای اطلاعات کاربر است.
type UserResponse struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Username  string   `json:"username"`
	FullName  string   `json:"fullName"`
	AvatarURL string   `json:"avatarUrl,omitempty"`
	Roles     []string `json:"roles"`
}

// ToUserResponse موجودیت User دامنه را به DTO پاسخ تبدیل می‌کند.
func ToUserResponse(u *user.User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FullName:  u.FullName,
		// AvatarURL: u.AvatarURL, // فرض بر اینکه این فیلد در User struct وجود دارد
		Roles:     u.Roles,
	}
}