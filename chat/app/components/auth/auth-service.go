package auth

import "chat/app/components/user"

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (authService *AuthService) AuthenticateUser(user *user.User, accessToken string) bool {
	return user.Password == accessToken
}

func (authService *AuthService) IsUserAuthenticated(token string) bool {
	return true
}
