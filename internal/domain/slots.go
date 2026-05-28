package domain

import (
	"errors"
	"fmt"
	"time"
)

const (
	OpeningHour = 8
	ClosingHour = 18
)

func NormalizeDate(value string) (time.Time, error) {
	day, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	return day.UTC(), nil
}

func DailySlots(day time.Time) []Slot {
	slots := make([]Slot, 0, ClosingHour-OpeningHour)
	baseDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)

	for hour := OpeningHour; hour < ClosingHour; hour++ {
		start := baseDay.Add(time.Duration(hour) * time.Hour)
		end := start.Add(time.Hour)
		slots = append(slots, Slot{
			Start:     start.Format("15:04"),
			End:       end.Format("15:04"),
			Available: true,
		})
	}

	return slots
}

func BuildSlotTime(day time.Time, startHour int) (string, string, error) {
	if startHour < OpeningHour || startHour >= ClosingHour {
		return "", "", fmt.Errorf("start hour must be between %02d:00 and %02d:00", OpeningHour, ClosingHour-1)
	}

	start := time.Date(day.Year(), day.Month(), day.Day(), startHour, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	return start.Format("15:04:05"), end.Format("15:04:05"), nil
}

func MarkReservedSlots(slots []Slot, reservations []Reservation) []Slot {
	lookup := make(map[string]struct{}, len(reservations))
	for _, reservation := range reservations {
		lookup[reservation.StartTime+"-"+reservation.EndTime] = struct{}{}
	}

	result := make([]Slot, len(slots))
	copy(result, slots)

	for i := range result {
		_, taken := lookup[result[i].Start+":00-"+result[i].End+":00"]
		result[i].Available = !taken
	}

	return result
}

func ValidateFutureSlot(day time.Time, startTime string, now time.Time) error {
	start, err := time.Parse("2006-01-02 15:04:05", day.Format("2006-01-02")+" "+startTime)
	if err != nil {
		return err
	}

	if start.Before(now.UTC()) {
		return errors.New("cannot reserve a slot in the past")
	}

	return nil
}
