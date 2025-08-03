// ================== internal/features/auth/validator.go ==================
package auth

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateRegister(req *RegisterRequest) error {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" {
		return errors.New("Email is required")
	}

	if !emailRegex.MatchString(req.Email) {
		return errors.New("Invalid email format")
	}

	if len(req.Password) < 6 {
		return errors.New("Password must be at least 6 characters")
	}

	if req.Name == "" {
		return errors.New("Name is required")
	}

	if len(req.Name) < 2 {
		return errors.New("Name must be at least 2 characters")
	}

	if len(req.Name) > 100 {
		return errors.New("Name cannot exceed 100 characters")
	}

	return nil
}

func ValidateLogin(req *LoginRequest) error {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" {
		return errors.New("Email is required")
	}

	if req.Password == "" {
		return errors.New("Password is required")
	}

	return nil
}

// TranslateAuthError converts database errors to user-friendly messages
func TranslateAuthError(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()

	if strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "already exists") {
		return "Email already registered"
	}

	if strings.Contains(errStr, "invalid credentials") {
		return "Invalid email or password"
	}

	if strings.Contains(errStr, "not found") {
		return "User not found"
	}

	return "Something went wrong. Please try again"
}
