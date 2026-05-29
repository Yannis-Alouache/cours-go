package server

import (
	"net/http"
	"strconv"

	"cours-go/internal/domain"

	"github.com/gin-gonic/gin"
)

func (s *Server) listRooms(c *gin.Context) {
	rows, err := s.pool.Query(c.Request.Context(), `SELECT id, name FROM rooms ORDER BY name`)
	if err != nil {
		writeAPIError(c, http.StatusInternalServerError, "rooms_query_failed", "Impossible de charger les salles")
		return
	}
	defer rows.Close()

	rooms := make([]domain.Room, 0)
	for rows.Next() {
		var room domain.Room
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			writeAPIError(c, http.StatusInternalServerError, "rooms_scan_failed", "Impossible de lire les salles")
			return
		}
		rooms = append(rooms, room)
	}

	c.JSON(http.StatusOK, rooms)
}

func (s *Server) roomAvailability(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("roomID"), 10, 64)
	if err != nil {
		writeAPIError(c, http.StatusBadRequest, "invalid_room_id", "Identifiant de salle invalide")
		return
	}

	day, err := domain.NormalizeDate(c.Query("date"))
	if err != nil {
		writeAPIError(c, http.StatusBadRequest, "invalid_date", err.Error())
		return
	}

	var exists bool
	if err := s.pool.QueryRow(c.Request.Context(), `SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)`, roomID).Scan(&exists); err != nil {
		writeAPIError(c, http.StatusInternalServerError, "room_lookup_failed", "Impossible de vérifier la salle")
		return
	}
	if !exists {
		writeAPIError(c, http.StatusNotFound, "room_not_found", "Salle introuvable")
		return
	}

	rows, err := s.pool.Query(
		c.Request.Context(),
		`SELECT id, room_id, to_char(start_time, 'HH24:MI:SS'), to_char(end_time, 'HH24:MI:SS')
		 FROM reservations
		 WHERE room_id = $1 AND day = $2::date
		 ORDER BY start_time`,
		roomID,
		day.Format("2006-01-02"),
	)
	if err != nil {
		writeAPIError(c, http.StatusInternalServerError, "availability_query_failed", "Impossible de charger les disponibilités")
		return
	}
	defer rows.Close()

	reservations := make([]domain.Reservation, 0)
	for rows.Next() {
		var reservation domain.Reservation
		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.StartTime, &reservation.EndTime); err != nil {
			writeAPIError(c, http.StatusInternalServerError, "availability_scan_failed", "Impossible de lire les réservations")
			return
		}
		reservation.Day = day.Format("2006-01-02")
		reservations = append(reservations, reservation)
	}

	slots := domain.MarkReservedSlots(domain.DailySlots(day), reservations)
	c.JSON(http.StatusOK, slots)
}
