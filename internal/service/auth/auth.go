package authservice

import (
	"context"

	"github.com/wafi04/backend/internal/handler/dto/request"
	"github.com/wafi04/backend/internal/handler/dto/response"
	authrepository "github.com/wafi04/backend/internal/repository/auth"
	"github.com/wafi04/backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepository *authrepository.AuthRepository
	log            logger.Logger
}

func NewAuthService(authrepo *authrepository.AuthRepository) *AuthService {
	return &AuthService{
		AuthRepository: authrepo,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (response.CreateUserResponse, error) {

	hashPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		s.log.Log(logger.ErrorLevel, "Failes Password : %v", err)
	}
	return s.AuthRepository.CreateUser(ctx, &request.CreateUserRequest{
		Name:       req.Name,
		Email:      req.Email,
		Password:   string(hashPw),
		Role:       "",
		IPAddress:  req.IPAddress,
		DeviceInfo: req.DeviceInfo,
	})
}

func (s *AuthService) Login(ctx context.Context, login *request.LoginRequest) (*response.LoginResponse, error) {
	return s.AuthRepository.Login(ctx, login)
}

func (s *AuthService) GetUser(ctx context.Context, req *request.GetUserRequest) (*response.GetUserResponse, error) {
	return s.AuthRepository.GetUser(ctx, req)
}

func (s *AuthService) VerifyEmail(ctx context.Context, req *request.VerifyEmailRequest) (*response.VerifyEmailResponse, error) {
	return s.AuthRepository.VerifyEmail(ctx, req)
}
func (s *AuthService) ResendVerification(ctx context.Context, req *request.ResendVerificationRequest) (*response.ResendVerificationResponse, error) {
	return s.AuthRepository.ResendVerification(ctx, req)
}

func (s *AuthService) Logout(ctx context.Context, req *request.LogoutRequest) (*response.LogoutResponse, error) {
	return s.AuthRepository.Logout(ctx, req)
}
func (s *AuthService) RevokeSession(ctx context.Context, req *request.RevokeSessionRequest) (*response.RevokeSessionResponse, error) {
	return s.AuthRepository.RevokeSession(ctx, req)
}
func (s *AuthService) RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.RefreshTokenResponse, error) {
	return s.AuthRepository.RefreshToken(ctx, req)
}
func (s *AuthService) ListSessions(ctx context.Context, req *request.ListSessionRequest) (*response.ListSessionResponse, error) {
	return s.AuthRepository.ListSessions(ctx, req)
}
func (s *AuthService) GetCurrentSession(ctx context.Context, ) (*response.ListSessionResponse, error) {
	return s.AuthRepository.GetCurrentSession(ctx, req)
}
