package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/wafi04/backend/internal/handler/dto/types"
)

var jwtSecretKey = []byte("jsjxakabxjaigisyqyg189")

type JWTClaims struct {
	UserID          string `json:"user_id"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Role            string `json:"role"`
	IsActive        bool   `json:"is_active"`
	IsEmailVerified bool   `json:"is_email_verified"`
	jwt.StandardClaims
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
func GenerateToken(user *types.UserInfo, day int64) (string, error) {
	claims := JWTClaims{
		UserID:          user.UserID,
		Email:           user.Email,
		Name:            user.Name,
		Role:            user.Role,
		IsEmailVerified: user.IsEmailVerified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(day) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "wafiuddin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

type contextKey string

const UserContextKey contextKey = "user"

func GetUserFromGinContext(c *gin.Context) (*types.UserInfo, error) {
	user, exists := c.Get(string(UserContextKey))
	if !exists {
		return nil, errors.New("user not found in context")
	}
	userInfo, ok := user.(*types.UserInfo)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}
	return userInfo, nil
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, "Bearer ")
			if len(parts) == 2 {
				tokenString := strings.TrimSpace(parts[1])
				if tokenString != "" {
					if claims, err := ValidateToken(tokenString); err == nil {
						user := &types.UserInfo{
							UserID:          claims.UserID,
							Email:           claims.Email,
							Name:            claims.Name,
							Role:            claims.Role,
							IsEmailVerified: claims.IsEmailVerified,
						}
						c.Set(string(UserContextKey), user)
						c.Next()
						return
					}
				}
			}
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No valid tokens found",
			})
			return
		}

		// Validasi refresh token
		claims, err := ValidateToken(refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid refresh token",
			})
			return
		}

		user := &types.UserInfo{
			UserID:          claims.UserID,
			Email:           claims.Email,
			Name:            claims.Name,
			Role:            claims.Role,
			IsEmailVerified: claims.IsEmailVerified,
		}

		newAccessToken, err := GenerateToken(user, 24)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate new access token",
			})
			return
		}

		c.Header("New-Access-Token", newAccessToken)

		c.Set(string(UserContextKey), user)
		c.Next()
	}
}

func SetRefreshTokenCookie(c *gin.Context, token string) {
    c.SetCookie(
        "refresh_token",
        token,
        int(168*3600),    
        "/",               
        "192.168.100.6",  
        false,             
        true,              
    )
}

func SetSessionCookie(c *gin.Context, sessionID string) {
    c.SetCookie(
        "session",
        sessionID,
        int(168*3600),     
        "/",              
        "192.168.100.6",   
        false,             
        true,              
    )
}

func ClearTokens(c *gin.Context) {
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.Header("Authorization", "")
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(string(UserContextKey))
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not found in context",
			})
			return
		}

		userInfo, ok := user.(*types.UserInfo)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user type in context",
			})
			return
		}

		for _, role := range allowedRoles {
			if userInfo.Role == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
	}
}
