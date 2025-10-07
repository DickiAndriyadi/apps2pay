package worker

import (
	"apps2pay/handlers"
	"context"
	"time"
)

func StartSeatCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			// Cari kursi locked yang expired
			_, _ = handlers.DB.Exec(ctx,
				`UPDATE seats SET status = 'available', locked_until = NULL
                 WHERE status = 'locked' AND locked_until < NOW()`)
		case <-ctx.Done():
			return
		}
	}
}
