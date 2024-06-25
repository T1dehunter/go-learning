package auth

import "chat/server/components/user"

type AuthService struct {
	userRepository *user.UserRepository
}

func NewAuthService(userRepo *user.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepo}
}

func (authService *AuthService) AuthenticateUser(user *user.User, accessToken string) bool {
	return user.Password == accessToken
}

func (authService *AuthService) IsUserAuthenticated(token string) bool {
	return true
}
