package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LockSeat(c *fiber.Ctx) error {
	scheduleID := c.Params("id")
	var req struct {
		SeatNumber string `json:"seat_number"`
		UserID     int    `json:"user_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Cek apakah kursi tersedia
	var status string
	err := DB.QueryRow(context.Background(),
		"SELECT status FROM seats WHERE schedule_id = $1 AND seat_number = $2",
		scheduleID, req.SeatNumber,
	).Scan(&status)
	if err != nil || status != "available" {
		return c.Status(409).JSON(fiber.Map{"error": "Seat not available"})
	}

	// Lock di Redis (15 menit)
	lockKey := "seat_lock:" + scheduleID + ":" + req.SeatNumber
	err = redisClient.Set(context.Background(), lockKey, req.UserID, 15*time.Minute).Err()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Lock failed"})
	}

	// Update status di DB
	lockedUntil := time.Now().Add(15 * time.Minute)
	_, err = DB.Exec(context.Background(),
		"UPDATE seats SET status = 'locked', locked_until = $1 WHERE schedule_id = $2 AND seat_number = $3",
		lockedUntil, scheduleID, req.SeatNumber,
	)
	if err != nil {
		redisClient.Del(context.Background(), lockKey)
		return c.Status(500).JSON(fiber.Map{"error": "DB update failed"})
	}

	return c.JSON(fiber.Map{"message": "Seat locked", "expires_in": "15 minutes"})
}

func ReleaseSeat(c *fiber.Ctx) error {
	scheduleID := c.Params("id")
	seatNumber := c.Query("seat")

	lockKey := "seat_lock:" + scheduleID + ":" + seatNumber
	redisClient.Del(context.Background(), lockKey)

	_, err := DB.Exec(context.Background(),
		"UPDATE seats SET status = 'available', locked_until = NULL WHERE schedule_id = $1 AND seat_number = $2",
		scheduleID, seatNumber,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Release failed"})
	}
	return c.JSON(fiber.Map{"message": "Seat released"})
}

func GetSeatsBySchedule(c *fiber.Ctx) error {
	scheduleID := c.Params("id")

	rows, err := DB.Query(context.Background(),
		"SELECT seat_number, status FROM seats WHERE schedule_id = $1 ORDER BY seat_number",
		scheduleID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch seats"})
	}
	defer rows.Close()

	var seats []map[string]interface{}
	for rows.Next() {
		var number, status string
		if err := rows.Scan(&number, &status); err != nil {
			continue
		}
		seats = append(seats, map[string]interface{}{
			"seat_number": number,
			"status":      status,
		})
	}

	return c.JSON(seats)
}

func ConfirmSeat(c *fiber.Ctx) error {
	scheduleID := c.Params("id")
	var req struct {
		SeatNumber string `json:"seat_number"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Pastikan kursi masih locked dan belum expired
	var status string
	var lockedUntil time.Time
	err := DB.QueryRow(context.Background(),
		"SELECT status, locked_until FROM seats WHERE schedule_id = $1 AND seat_number = $2",
		scheduleID, req.SeatNumber,
	).Scan(&status, &lockedUntil)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Seat not found"})
	}

	if status != "locked" || lockedUntil.Before(time.Now()) {
		return c.Status(409).JSON(fiber.Map{"error": "Seat is not locked or already expired"})
	}

	// Ubah status jadi 'sold'
	_, err = DB.Exec(context.Background(),
		"UPDATE seats SET status = 'sold', sold_at = NOW(), locked_until = NULL WHERE schedule_id = $1 AND seat_number = $2",
		scheduleID, req.SeatNumber,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to confirm seat"})
	}

	// Opsional: hapus lock dari Redis
	lockKey := "seat_lock:" + scheduleID + ":" + req.SeatNumber
	redisClient.Del(context.Background(), lockKey)

	return c.JSON(fiber.Map{"message": "Seat confirmed and sold"})
}
