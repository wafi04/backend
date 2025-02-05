package authhandler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/middleware"
	httpresponse "github.com/wafi04/backend/pkg/response"
	request "github.com/wafi04/backend/pkg/types/req"
	authrepo "github.com/wafi04/backend/services/auth/repository"
	authservice "github.com/wafi04/backend/services/auth/service"
)

type AuthHandler struct {
	AuthService *authservice.AuthService
	log         logger.Logger
}

func NewAuthHandler(authservice *authservice.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authservice,
	}
}

func (s *AuthHandler) Verify(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "No refresh token found",
		})
		return
	}
	log.Printf("Received refresh token: %s", refreshToken)
	log.Printf("Validating refresh token...")
	claims, err := middleware.ValidateToken(refreshToken)
	if err != nil {
		log.Printf("Validation failed: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("Invalid refresh token: %v", err),
		})
		return
	}
	log.Printf("Refresh token validated successfully: %+v", claims)
	sessionID, err := c.Cookie("session")
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "No session found")
		return
	}

	session, err := s.AuthService.AuthRepository.GetCurrentSession(c.Request.Context(), &authrepo.ReqSession{
		SessionID: sessionID,
	})
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Invalid session")
		return
	}

	// Validasi session
	if !session.IsActive {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Session inactive")
		return
	}

	if session.ExpiresAt < time.Now().Unix() {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Session expired")
		return
	}

	userResp, err := s.AuthService.GetUser(c.Request.Context(), &request.GetUserRequest{
		UserID: session.UserID,
	})
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed to fetch user")
		return
	}

	// _, err = s.AuthService.AuthRepository.UpdateSessionActivity(c.Request.Context(), session.SessionID)
	// if err != nil {
	// 	log.Printf("Failed to update session activity: %v", err)
	// }

	httpresponse.SendSuccessResponse(c, http.StatusOK, "User retrieved successfully", userResp)

}

func (s *AuthHandler) CreateUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	log.Printf("Received CreateUser request for user: %s", req.Name)
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to create user : %v", err.Error())
		return
	}
	user, err := s.AuthService.CreateUser(c.Request.Context(), &request.CreateUserRequest{
		Name:       req.Name,
		Email:      req.Email,
		Role:       "",
		Password:   req.Password,
		Picture:    "",
		IPAddress:  clientIP,
		DeviceInfo: userAgent,
	})
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed to create user")
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusCreated, "Created user successfully", user)
}

func (s *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed to login")
		return
	}

	resp, err := s.AuthService.Login(c.Request.Context(), &request.LoginRequest{
		Email:      req.Email,
		Password:   req.Password,
		DeviceInfo: userAgent,
		IPAddress:  clientIP,
	})
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "user not found : %s", err.Error())
		return
	}

	c.Header("Access-Control-Allow-Credentials", "true")
	middleware.SetRefreshTokenCookie(c, resp.Refresh_token)
	middleware.SetSessionCookie(c, resp.SessionInfo.SessionID)
	httpresponse.SendSuccessResponse(c, http.StatusOK, "Login user successfully", resp)
}

func (s *AuthHandler) VerifyEmail(c *gin.Context) {
	var req request.VerifyEmailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Invalid verification token")
		return
	}

	resp, err := s.AuthService.VerifyEmail(c.Request.Context(), &req)
	if err != nil {
		log.Printf("Failed to verify email: %v", err)
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Email verification failed")
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Email verified successfully", resp)
}

func (s *AuthHandler) ResendVerification(c *gin.Context) {
	var req request.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := s.AuthService.ResendVerification(c.Request.Context(), &req)
	if err != nil {
		log.Printf("Failed to resend verification: %v", err)
		httpresponse.SendErrorResponse(c, http.StatusInternalServerError, "Failed to resend verification")
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Verification email resent", resp)
}
func (s *AuthHandler) GetUser(c *gin.Context) {
	s.log.Log(logger.InfoLevel, "Incoming request")
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	sessionID, err := c.Cookie("session")
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "No session found")
		return
	}

	session, err := s.AuthService.AuthRepository.GetCurrentSession(c.Request.Context(), &authrepo.ReqSession{
		SessionID: sessionID,
	})
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Invalid session")
		return
	}

	if !session.IsActive {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Session inactive")
		return
	}
	s.log.Log(logger.InfoLevel, "SESSION TIME  : %d", session.ExpiresAt)
	s.log.Log(logger.InfoLevel, "TIME  : %d", time.Now().Unix())
	if session.ExpiresAt < time.Now().Unix() {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Session expired")
		return
	}
	if session.UserID != user.UserID {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Invalid session ownership")
		return
	}

	// Ambil data pengguna
	userResp, err := s.AuthService.GetUser(c.Request.Context(), &request.GetUserRequest{
		UserID: session.UserID,
	})
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed to fetch user")
		return
	}

	// _, err = s.AuthService.AuthRepository.UpdateSessionActivity(c.Request.Context(), session.SessionID)
	// if err != nil {
	// 	log.Printf("Failed to update session activity: %v", err)
	// }

	httpresponse.SendSuccessResponse(c, http.StatusOK, "User retrieved successfully", userResp)
}

func (s *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		AccessToken string `json:"access_token"`
	}

	sessionID, err := c.Cookie("session")
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "No session found")
		return
	}
	if err = c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	middleware.ClearTokens(c)

	resp, err := s.AuthService.Logout(c.Request.Context(), &request.LogoutRequest{
		AccessToken: req.AccessToken,
		SessionID:   sessionID,
	})
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusInternalServerError, "Logout failed")
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Logged out successfully", resp)
}

func (s *AuthHandler) RevokeSession(c *gin.Context) {
	var req request.RevokeSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := s.AuthService.RevokeSession(c.Request.Context(), &req)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusInternalServerError, "Failed to revoke session")
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Session revoked successfully", resp)
}

func (s *AuthHandler) RefreshToken(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		s.log.Log(logger.DebugLevel, "err : %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	session, err := c.Cookie("session")
	if err != nil {
		s.log.Log(logger.DebugLevel, "err : %v", err)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session token"})
		return
	}
	_, err = middleware.ValidateToken(cookie)
	if err != nil {
		s.log.Log(logger.DebugLevel, "err : %v", err)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	resp, err := s.AuthService.RefreshToken(c.Request.Context(), &request.RefreshTokenRequest{
		RefreshToken: cookie,
		SessionID:    session,
	})
	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Token refresh failed")
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Token refreshed successfully", resp)
}

func (s *AuthHandler) ListSessions(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	resp, err := s.AuthService.ListSessions(c.Request.Context(), &request.ListSessionRequest{
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Failed to list sessions: %v", err)
		httpresponse.SendErrorResponse(c, http.StatusInternalServerError, "Failed to list sessions")
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Sessions listed successfully", resp)
}
