package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"library_management/models"
	"library_management/services"
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
			fmt.Println("Exiting. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please choose a valid option.")
		}
		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("=== Library Management System ===")
	fmt.Println("1) Add Book")
	fmt.Println("2) Remove Book")
	fmt.Println("3) Add Member")
	fmt.Println("4) Borrow Book")
	fmt.Println("5) Return Book")
	fmt.Println("6) List Available Books")
	fmt.Println("7) List Borrowed Books by Member")
	fmt.Println("8) Exit")
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

