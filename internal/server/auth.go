package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"cours-go/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type claims struct {
	UserID int64  `json:"uid"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *Server) login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAPIError(c, http.StatusBadRequest, "invalid_payload", "Payload de login invalide")
		return
	}

	var user domain.User
	var passwordHash string
	err := s.pool.QueryRow(
		c.Request.Context(),
		`SELECT id, email, password_hash FROM users WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Email, &passwordHash)
	if err != nil {
		writeAPIError(c, http.StatusUnauthorized, "invalid_credentials", "Identifiants invalides")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		writeAPIError(c, http.StatusUnauthorized, "invalid_credentials", "Identifiants invalides")
		return
	}

	token, err := s.signToken(user)
	if err != nil {
		writeAPIError(c, http.StatusInternalServerError, "token_error", "Impossible de créer le token")
		return
	}

	c.JSON(http.StatusOK, domain.LoginResponse{
		Token: token,
		User:  user,
	})
}

func (s *Server) me(c *gin.Context) {
	userID := currentUserID(c)

	var user domain.User
	err := s.pool.QueryRow(
		c.Request.Context(),
		`SELECT id, email FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Email)
	if err != nil {
		writeAPIError(c, http.StatusUnauthorized, "unknown_user", "Utilisateur introuvable")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) signToken(user domain.User) (string, error) {
	claims := claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.TokenTTL) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprint(user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func parseClaims(tokenString, secret string) (*claims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	customClaims, ok := parsed.Claims.(*claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	return customClaims, nil
}
