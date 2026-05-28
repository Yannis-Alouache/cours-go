package server

import (
	"net/http"

	"cours-go/internal/config"
	"cours-go/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	cfg  config.Server
	pool *pgxpool.Pool
}

func NewRouter(cfg config.Server, pool *pgxpool.Pool) http.Handler {
	srv := &Server{
		cfg:  cfg,
		pool: pool,
	}

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := router.Group("/auth")
	auth.POST("/login", srv.login)

	protected := router.Group("/")
	protected.Use(srv.authMiddleware())
	protected.GET("/auth/me", srv.me)
	protected.GET("/rooms", srv.listRooms)
	protected.GET("/rooms/:roomID/availability", srv.roomAvailability)
	protected.GET("/reservations/me", srv.listMyReservations)
	protected.POST("/reservations", srv.createReservation)
	protected.DELETE("/reservations/:id", srv.deleteReservation)
	protected.GET("/me/config", srv.getConfig)
	protected.PUT("/me/config", srv.putConfig)
	protected.GET("/me/state", srv.getState)
	protected.PUT("/me/state", srv.putState)

	return router
}

func writeAPIError(c *gin.Context, status int, code, message string) {
	c.JSON(status, domain.APIError{Code: code, Message: message})
}
