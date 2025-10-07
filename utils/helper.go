package utils

import "fmt"

func GenerateSeatNumbers(total int) []string {
	var seats []string
	rows := "ABCDEFGHIJ" // 10 baris
	cols := 10           // 10 kolom → 100 kursi
	for i := 0; i < total && i < len(rows)*cols; i++ {
		row := rows[i/cols]
		col := (i % cols) + 1
		seats = append(seats, fmt.Sprintf("%c%d", row, col))
	}
	return seats
}
