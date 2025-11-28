package concurrency

import (
	"fmt"
	"sync"

	"library_management/services"
)

// ReservationRequest represents a single reservation attempt.
type ReservationRequest struct {
	BookID   int
	MemberID int
	Resp     chan error
}

// StartReservationWorkerPool starts a pool of worker goroutines that process reservation requests concurrently.
// - lib: the Library instance to use
// - reqCh: incoming requests channel
// - workerCount: number of worker goroutines to spawn
// - wg: WaitGroup for graceful shutdown (optional - can be nil)
func StartReservationWorkerPool(lib *services.Library, reqCh <-chan ReservationRequest, workerCount int, wg *sync.WaitGroup) {
	for i := 0; i < workerCount; i++ {
		if wg != nil {
			wg.Add(1)
		}
		go func(workerID int) {
			defer func() {
				if wg != nil {
					wg.Done()
				}
			}()
			for req := range reqCh {
				err := lib.ReserveBook(req.BookID, req.MemberID)
				// send result back (non-blocking safety if receiver not listening)
				select {
				case req.Resp <- err:
				default:
					// If nobody is listening, swallow
				}
				// Also log processing
				if err != nil {
					fmt.Printf("[Worker %d] Failed to reserve Book %d for Member %d: %v\n", workerID, req.BookID, req.MemberID, err)
				} else {
					fmt.Printf("[Worker %d] Reserved Book %d for Member %d\n", workerID, req.BookID, req.MemberID)
				}
			}
		}(i + 1)
	}
}
