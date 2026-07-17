package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (s *Server) getConfig(c *gin.Context) {
	s.getSyncJSON(c, "config_json")
}

func (s *Server) putConfig(c *gin.Context) {
	s.putSyncJSON(c, "config_json")
}

func (s *Server) getSyncJSON(c *gin.Context, column string) {
	var payload []byte
	err := s.pool.QueryRow(
		c.Request.Context(),
		`SELECT `+column+`::text
		 FROM client_sync
		 WHERE user_id = $1`,
		currentUserID(c),
	).Scan(&payload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.Data(http.StatusOK, "application/json", []byte("{}"))
			return
		}

		writeAPIError(c, http.StatusInternalServerError, "sync_read_failed", "Impossible de charger le JSON")
		return
	}

	c.Data(http.StatusOK, "application/json", payload)
}

func (s *Server) putSyncJSON(c *gin.Context, column string) {
	body, err := c.GetRawData()
	if err != nil || len(body) == 0 {
		writeAPIError(c, http.StatusBadRequest, "invalid_payload", "Payload JSON requis")
		return
	}

	query := `INSERT INTO client_sync (user_id, ` + column + `, updated_at)
		VALUES ($1, $2::jsonb, NOW())
		ON CONFLICT (user_id) DO UPDATE
		SET ` + column + ` = EXCLUDED.` + column + `, updated_at = NOW()`

	if _, err := s.pool.Exec(c.Request.Context(), query, currentUserID(c), string(body)); err != nil {
		writeAPIError(c, http.StatusBadRequest, "sync_update_failed", "Impossible d'enregistrer le JSON")
		return
	}

	c.Status(http.StatusNoContent)
}
