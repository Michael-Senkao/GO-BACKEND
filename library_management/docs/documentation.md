# Library Management System (Concurrent Reservation) - Documentation

## Overview
This project adds concurrent reservation support to the original library management system. It uses Goroutines, Channels, Mutexes, and timers to safely process multiple reservation requests.

## Key Concurrency Components
- **Mutex (sync.Mutex)**: `services.Library` uses a mutex `mu` to protect shared state (books, members, reservations, timers). All state-changing operations obtain the lock to prevent race conditions.
- **Channels**: `concurrency.ReservationRequest` channel (`reqCh`) is used to queue incoming reservation requests. Worker goroutines read from the channel and process requests concurrently.
- **Worker Pool (Goroutines)**: The `StartReservationWorkerPool` function spawns a configurable number of worker Goroutines which pull requests from the request channel and call `Library.ReserveBook`.
- **Timers (`time.Timer`)**: When a reservation is accepted, a `time.Timer` is created (5 seconds). If the member does not borrow the reserved book within 5 seconds, the timer's callback auto-cancels the reservation (cleans up internal state).
- **Auto-Cancellation**: Timer callbacks obtain the same mutex to safely mutate state. They verify the reservation still matches the expected member before cancellation.

## API (CLI)
- Add Book
- Remove Book (can't remove when borrowed/reserved)
- Add Member
- Borrow Book
- Return Book
- List Available Books
- List Borrowed Books by Member
- Reserve Book (single)
- Simulate Concurrent Reservations (creates many requests and processes them via worker pool)

## How to Run
1. Ensure Go is installed.
2. From the project root:
   ```bash
   go run ./...
