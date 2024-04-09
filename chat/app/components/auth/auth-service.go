package auth

import "fmt"

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (authService *AuthService) AuthenticateUser(userID int32, accessToken string) string {
	tokens := make(map[string]bool)
	tokens["Test1234"] = true

	_, exists := tokens[accessToken]
	if exists {
		fmt.Println("User authentication success")
	} else {
		fmt.Println("User authentication error")
	}

	return ""
}

func (authService *AuthService) IsUserAuthenticated(token string) bool {
	return true
}
