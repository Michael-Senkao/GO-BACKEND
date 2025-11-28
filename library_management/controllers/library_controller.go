package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"library_management/models"
	"library_management/services"
	"library_management/concurrency"
)

// Controller wraps service for CLI interactions.
type Controller struct {
	lib *services.Library
}

// NewController returns a new Controller instance.
func NewController(lib *services.Library) *Controller {
	return &Controller{lib: lib}
}

// Start runs the console menu loop.
func (c *Controller) Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		printMenu()
		fmt.Print("Enter choice: ")
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			c.handleAddBook(reader)
		case "2":
			c.handleRemoveBook(reader)
		case "3":
			c.handleAddMember(reader)
		case "4":
			c.handleBorrowBook(reader)
		case "5":
			c.handleReturnBook(reader)
		case "6":
			c.handleListAvailableBooks()
		case "7":
			c.handleListBorrowedByMember(reader)
		case "8":
			c.handleReserveBook(reader)
		case "9":
			c.handleSimulateConcurrentReservations(reader)
		case "10":
			fmt.Println("Exiting. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please choose a valid option.")
		}
		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("=== Library Management System (Concurrent Reservation) ===")
	fmt.Println("1) Add Book")
	fmt.Println("2) Remove Book")
	fmt.Println("3) Add Member")
	fmt.Println("4) Borrow Book")
	fmt.Println("5) Return Book")
	fmt.Println("6) List Available Books")
	fmt.Println("7) List Borrowed Books by Member")
	fmt.Println("8) Reserve Book (single request)")
	fmt.Println("9) Simulate Concurrent Reservations")
	fmt.Println("10) Exit")
}

func (c *Controller) handleAddBook(reader *bufio.Reader) {
	fmt.Println("--- Add Book ---")
	id := promptInt(reader, "Book ID: ")
	title := promptString(reader, "Title: ")
	author := promptString(reader, "Author: ")

	book := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: "Available",
	}
	c.lib.AddBook(book)
	fmt.Println("Book added successfully.")
}

func (c *Controller) handleRemoveBook(reader *bufio.Reader) {
	fmt.Println("--- Remove Book ---")
	id := promptInt(reader, "Book ID to remove: ")
	err := c.lib.RemoveBook(id)
	if err != nil {
		fmt.Println("Error removing book:", err)
	} else {
		fmt.Println("Book removed successfully.")
	}
}

func (c *Controller) handleAddMember(reader *bufio.Reader) {
	fmt.Println("--- Add Member ---")
	id := promptInt(reader, "Member ID: ")
	name := promptString(reader, "Member Name: ")
	member := models.Member{
		ID:   id,
		Name: name,
	}
	if err := c.lib.AddMember(member); err != nil {
		fmt.Println("Error adding member:", err)
	} else {
		fmt.Println("Member added successfully.")
	}
}

func (c *Controller) handleBorrowBook(reader *bufio.Reader) {
	fmt.Println("--- Borrow Book ---")
	bookID := promptInt(reader, "Book ID: ")
	memberID := promptInt(reader, "Member ID: ")
	if err := c.lib.BorrowBook(bookID, memberID); err != nil {
		fmt.Println("Error borrowing book:", err)
	} else {
		fmt.Println("Book borrowed successfully.")
	}
}

func (c *Controller) handleReturnBook(reader *bufio.Reader) {
	fmt.Println("--- Return Book ---")
	bookID := promptInt(reader, "Book ID: ")
	memberID := promptInt(reader, "Member ID: ")
	if err := c.lib.ReturnBook(bookID, memberID); err != nil {
		fmt.Println("Error returning book:", err)
	} else {
		fmt.Println("Book returned successfully.")
	}
}

func (c *Controller) handleListAvailableBooks() {
	fmt.Println("--- Available Books ---")
	books := c.lib.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}
	for _, b := range books {
		fmt.Printf("ID: %d | Title: %s | Author: %s\n", b.ID, b.Title, b.Author)
	}
}

func (c *Controller) handleListBorrowedByMember(reader *bufio.Reader) {
	fmt.Println("--- List Borrowed Books by Member ---")
	memberID := promptInt(reader, "Member ID: ")
	books, err := c.lib.ListBorrowedBooks(memberID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(books) == 0 {
		fmt.Println("Member has not borrowed any books.")
		return
	}
	for _, b := range books {
		fmt.Printf("ID: %d | Title: %s | Author: %s\n", b.ID, b.Title, b.Author)
	}
}

func (c *Controller) handleReserveBook(reader *bufio.Reader) {
	fmt.Println("--- Reserve Book ---")
	bookID := promptInt(reader, "Book ID: ")
	memberID := promptInt(reader, "Member ID: ")
	err := c.lib.ReserveBook(bookID, memberID)
	if err != nil {
		fmt.Println("Reservation failed:", err)
	} else {
		fmt.Println("Reservation successful. You have 5 seconds to borrow the book before auto-cancel.")
	}
}

// Simulate many members simultaneously trying to reserve the same (or different) books
func (c *Controller) handleSimulateConcurrentReservations(reader *bufio.Reader) {
	fmt.Println("--- Simulate Concurrent Reservations ---")
	bookID := promptInt(reader, "Book ID to contest: ")
	fmt.Println("We'll simulate multiple members trying to reserve the same book simultaneously.")
	count := promptInt(reader, "How many concurrent attempts? (e.g., 5): ")
	workerCount := promptInt(reader, "How many worker goroutines to process requests? (e.g., 3): ")

	// Create request channel and worker pool
	reqCh := make(chan concurrency.ReservationRequest)
	var wg sync.WaitGroup
	concurrency.StartReservationWorkerPool(c.lib, reqCh, workerCount, &wg)

	// Launch goroutines that send reservation requests nearly simultaneously
	respChans := make([]chan error, count)
	for i := 0; i < count; i++ {
		resp := make(chan error, 1)
		respChans[i] = resp
		memberID := 100 + i // create simulated member IDs (100,101,...)
		// ensure members exist
		_ = c.lib.AddMember(models.Member{ID: memberID, Name: fmt.Sprintf("SimMember-%d", memberID)})
		req := concurrency.ReservationRequest{
			BookID:   bookID,
			MemberID: memberID,
			Resp:     resp,
		}
		// Send requests in separate goroutines to simulate near-simultaneous arrivals
		go func(r concurrency.ReservationRequest) {
			reqCh <- r
		}(req)
		// tiny sleep to better simulate near-simultaneous but not perfectly ordered bursts
		time.Sleep(10 * time.Millisecond)
	}

	// Collect responses with a small timeout window
	for i, ch := range respChans {
		select {
		case err := <-ch:
			if err != nil {
				fmt.Printf("Simulated Member %d: reservation failed: %v\n", 100+i, err)
			} else {
				fmt.Printf("Simulated Member %d: reservation succeeded\n", 100+i)
			}
		case <-time.After(2 * time.Second):
			fmt.Printf("Simulated Member %d: no response (timed out)\n", 100+i)
		}
	}

	// close the reqCh and wait for workers to finish processing queued messages
	close(reqCh)
	wg.Wait()

	fmt.Println("Simulation complete. Waiting 6 seconds to observe any auto-cancellations (if any).")
	time.Sleep(6 * time.Second)
	fmt.Println("Done waiting. Simulation finished.")
}

// Helper prompts
func promptString(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptInt(reader *bufio.Reader, prompt string) int {
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		n, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Please enter a valid integer.")
			continue
		}
		return n
	}
}
