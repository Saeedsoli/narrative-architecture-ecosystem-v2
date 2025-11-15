// apps/backend/test/unit/user_test.go

package unit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"narrative-architecture/apps/backend/internal/domain/user"
)

func TestUser_SetPassword(t *testing.T) {
	u := user.NewUser("id", "email", "username", "fullName")
	
	err := u.SetPassword("MyStrongPassword123!")
	
	assert.NoError(t, err)
	assert.NotEmpty(t, u.PasswordHash)
	assert.NotEqual(t, "MyStrongPassword123!", u.PasswordHash)
}

func TestUser_CheckPassword(t *testing.T) {
	u := user.NewUser("id", "email", "username", "fullName")
	password := "MyStrongPassword123!"
	
	// هش کردن رمز عبور
	_ = u.SetPassword(password)
	
	// تست با رمز عبور صحیح
	assert.True(t, u.CheckPassword(password), "Should return true for correct password")
	
	// تست با رمز عبور اشتباه
	assert.False(t, u.CheckPassword("WrongPassword"), "Should return false for incorrect password")
}

func TestNewUser_Defaults(t *testing.T) {
	u := user.NewUser("id", "email", "username", "fullName")
	
	assert.Equal(t, user.StatusActive, u.Status)
	assert.Equal(t, []string{"user"}, u.Roles)
	assert.False(t, u.CreatedAt.IsZero())
	assert.False(t, u.UpdatedAt.IsZero())
}