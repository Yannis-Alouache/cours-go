package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"cours-go/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
)

func (s *Server) listMyReservations(c *gin.Context) {
	rows, err := s.pool.Query(
		c.Request.Context(),
		`SELECT r.id, r.room_id, rooms.name, to_char(r.day, 'YYYY-MM-DD'), to_char(r.start_time, 'HH24:MI:SS'), to_char(r.end_time, 'HH24:MI:SS')
		 FROM reservations r
		 JOIN rooms ON rooms.id = r.room_id
		 WHERE r.user_id = $1
		 ORDER BY r.day, r.start_time`,
		currentUserID(c),
	)
	if err != nil {
		writeAPIError(c, http.StatusInternalServerError, "reservations_query_failed", "Impossible de charger les réservations")
		return
	}
	defer rows.Close()

	reservations := make([]domain.Reservation, 0)
	for rows.Next() {
		var reservation domain.Reservation
		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.RoomName, &reservation.Day, &reservation.StartTime, &reservation.EndTime); err != nil {
			writeAPIError(c, http.StatusInternalServerError, "reservations_scan_failed", "Impossible de lire les réservations")
			return
		}
		reservations = append(reservations, reservation)
	}

	c.JSON(http.StatusOK, reservations)
}

func (s *Server) createReservation(c *gin.Context) {
	var req domain.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAPIError(c, http.StatusBadRequest, "invalid_payload", "Payload de réservation invalide")
		return
	}

	day, err := domain.NormalizeDate(req.Day)
	if err != nil {
		writeAPIError(c, http.StatusBadRequest, "invalid_date", err.Error())
		return
	}

	startTime, endTime, err := domain.BuildSlotTime(day, req.StartHour)
	if err != nil {
		writeAPIError(c, http.StatusBadRequest, "invalid_slot", err.Error())
		return
	}

	if err := domain.ValidateFutureSlot(day, startTime, time.Now()); err != nil {
		writeAPIError(c, http.StatusBadRequest, "past_slot", err.Error())
		return
	}

	var reservation domain.Reservation
	err = s.pool.QueryRow(
		c.Request.Context(),
		`INSERT INTO reservations (user_id, room_id, day, start_time, end_time)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, room_id, to_char(day, 'YYYY-MM-DD'), to_char(start_time, 'HH24:MI:SS'), to_char(end_time, 'HH24:MI:SS')`,
		currentUserID(c),
		req.RoomID,
		day.Format("2006-01-02"),
		startTime,
		endTime,
	).Scan(&reservation.ID, &reservation.RoomID, &reservation.Day, &reservation.StartTime, &reservation.EndTime)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeAPIError(c, http.StatusConflict, "slot_taken", "Créneau déjà réservé")
			return
		}

		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			writeAPIError(c, http.StatusBadRequest, "room_not_found", "Salle introuvable")
			return
		}

		writeAPIError(c, http.StatusInternalServerError, "reservation_create_failed", "Impossible de créer la réservation")
		return
	}

	if err := s.pool.QueryRow(c.Request.Context(), `SELECT name FROM rooms WHERE id = $1`, req.RoomID).Scan(&reservation.RoomName); err != nil {
		writeAPIError(c, http.StatusInternalServerError, "room_lookup_failed", "Réservation créée mais salle introuvable")
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

func (s *Server) deleteReservation(c *gin.Context) {
	reservationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		writeAPIError(c, http.StatusBadRequest, "invalid_reservation_id", "Identifiant de réservation invalide")
		return
	}

	commandTag, err := s.pool.Exec(
		c.Request.Context(),
		`DELETE FROM reservations WHERE id = $1 AND user_id = $2`,
		reservationID,
		currentUserID(c),
	)
	if err != nil {
		writeAPIError(c, http.StatusInternalServerError, "reservation_delete_failed", "Impossible d'annuler la réservation")
		return
	}

	if commandTag.RowsAffected() == 0 {
		writeAPIError(c, http.StatusNotFound, "reservation_not_found", "Réservation introuvable")
		return
	}

	c.Status(http.StatusNoContent)
}
