package handlers

import (
	"apps2pay/models"
	"apps2pay/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func CreateSchedule(c *fiber.Ctx) error {
	var s models.Schedule
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if s.TotalSeats <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Total seats must be > 0"})
	}

	// Mulai transaksi
	tx, err := DB.Begin(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Transaction failed"})
	}
	defer tx.Rollback(context.Background())

	// Simpan jadwal
	var scheduleID int
	err = tx.QueryRow(context.Background(),
		`INSERT INTO schedules (movie_title, cinema_branch, city, show_time, total_seats)
         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		s.MovieTitle, s.CinemaBranch, s.City, s.ShowTime, s.TotalSeats,
	).Scan(&scheduleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create schedule"})
	}

	// Generate & insert seats
	seatNumbers := utils.GenerateSeatNumbers(s.TotalSeats)
	for _, num := range seatNumbers {
		_, err := tx.Exec(context.Background(),
			"INSERT INTO seats (schedule_id, seat_number, status) VALUES ($1, $2, 'available')",
			scheduleID, num,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create seats"})
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Commit failed"})
	}

	s.ID = scheduleID
	return c.Status(201).JSON(s)
}

func GetSchedules(c *fiber.Ctx) error {
	var schedules []models.Schedule
	rows, err := DB.Query(context.Background(), "SELECT id, movie_title, cinema_branch, city, show_time, total_seats, status FROM schedules")
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(schedules)
		}

		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Schedule
		err := rows.Scan(&s.ID, &s.MovieTitle, &s.CinemaBranch, &s.City, &s.ShowTime, &s.TotalSeats, &s.Status)
		if err != nil {
			continue
		}
		schedules = append(schedules, s)
	}
	return c.JSON(schedules)
}

func GetSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	var s models.Schedule
	err := DB.QueryRow(context.Background(),
		"SELECT id, movie_title, cinema_branch, city, show_time, total_seats, status FROM schedules WHERE id = $1",
		id,
	).Scan(&s.ID, &s.MovieTitle, &s.CinemaBranch, &s.City, &s.ShowTime, &s.TotalSeats, &s.Status)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Schedule not found"})
	}
	return c.JSON(s)
}

func UpdateSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	var s models.Schedule
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	query := `UPDATE schedules SET movie_title=$1, cinema_branch=$2, city=$3, show_time=$4, total_seats=$5, status=$6, updated_at=$7
              WHERE id = $8`
	_, err := DB.Exec(context.Background(), query,
		s.MovieTitle, s.CinemaBranch, s.City, s.ShowTime, s.TotalSeats, s.Status, time.Now(), id,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
	}
	s.ID, _ = c.ParamsInt("id")
	return c.JSON(s)
}

func DeleteSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := DB.Exec(context.Background(), "DELETE FROM schedules WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Delete failed"})
	}
	return c.SendStatus(204)
}

func CancelSchedule(c *fiber.Ctx) error {
	id := c.Params("id")

	// Update status jadwal
	_, err := DB.Exec(context.Background(), "UPDATE schedules SET status = 'cancelled' WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to cancel schedule"})
	}

	// Restok semua kursi & catat refund
	rows, err := DB.Query(context.Background(),
		"SELECT id FROM seats WHERE schedule_id = $1 AND status = 'sold'", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Refund lookup failed"})
	}
	defer rows.Close()

	for rows.Next() {
		var seatID int
		rows.Scan(&seatID)
		// Simpan ke log refund
		DB.Exec(context.Background(),
			"INSERT INTO refunds (schedule_id, seat_id, reason) VALUES ($1, $2, 'cinema_cancel')",
			id, seatID,
		)
	}

	// Kembalikan semua kursi ke available
	DB.Exec(context.Background(),
		"UPDATE seats SET status = 'available', sold_at = NULL WHERE schedule_id = $1", id)

	return c.JSON(fiber.Map{"message": "Schedule cancelled and all tickets refunded"})
}
