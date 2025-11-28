package services

import (
	"errors"
	"fmt"

	"library_management/models"
)

// LibraryManager defines the operations for the library.
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
	AddMember(m models.Member) error
	GetMember(memberID int) (*models.Member, error)
}

// Library implements LibraryManager.
type Library struct {
	books   map[int]models.Book
	members map[int]models.Member
}

// NewLibrary creates a new Library instance.
func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}

// AddBook adds a new book to the library.
func (l *Library) AddBook(book models.Book) {
	// If ID exists, overwrite with new details
	if _, exists := l.books[book.ID]; exists {
		// keep existing status if unspecified
		if book.Status == "" {
			book.Status = l.books[book.ID].Status
		}
	}
	if book.Status == "" {
		book.Status = "Available"
	}
	l.books[book.ID] = book
}

// RemoveBook removes a book from the library by its ID.
func (l *Library) RemoveBook(bookID int) error {
	b, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	if b.Status == "Borrowed" {
		return errors.New("cannot remove a borrowed book")
	}
	delete(l.books, bookID)
	return nil
}

// AddMember adds a new member to the library.
func (l *Library) AddMember(m models.Member) error {
	if _, exists := l.members[m.ID]; exists {
		return errors.New("member with this ID already exists")
	}
	m.BorrowedBooks = []models.Book{}
	l.members[m.ID] = m
	return nil
}

// GetMember returns a pointer to a member if exists.
func (l *Library) GetMember(memberID int) (*models.Member, error) {
	m, ok := l.members[memberID]
	if !ok {
		return nil, errors.New("member not found")
	}
	return &m, nil
}

// BorrowBook allows a member to borrow a book if it is available.
func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}
	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}
	// mark book as borrowed
	book.Status = "Borrowed"
	l.books[bookID] = book

	// add a copy to member's borrowed books
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member
	return nil
}

// ReturnBook allows a member to return a borrowed book.
func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	// find the book in the member's borrowed slice
	foundIdx := -1
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			foundIdx = i
			break
		}
	}
	if foundIdx == -1 {
		return errors.New("member did not borrow this book")
	}

	// remove the book from member's borrowed slice
	member.BorrowedBooks = append(member.BorrowedBooks[:foundIdx], member.BorrowedBooks[foundIdx+1:]...)
	l.members[memberID] = member

	// update book status to available
	book.Status = "Available"
	l.books[bookID] = book

	return nil
}

// ListAvailableBooks lists all available books.
func (l *Library) ListAvailableBooks() []models.Book {
	list := make([]models.Book, 0)
	for _, b := range l.books {
		if b.Status == "Available" {
			list = append(list, b)
		}
	}
	return list
}

// ListBorrowedBooks lists all books borrowed by a specific member.
func (l *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	member, ok := l.members[memberID]
	if !ok {
		return nil, errors.New("member not found")
	}
	return member.BorrowedBooks, nil
}

// Helper: Seed sample data (optional)
func (l *Library) SeedSampleData() {
	l.AddBook(models.Book{ID: 1, Title: "1984", Author: "George Orwell", Status: "Available"})
	l.AddBook(models.Book{ID: 2, Title: "The Hobbit", Author: "J.R.R. Tolkien", Status: "Available"})
	l.AddBook(models.Book{ID: 3, Title: "Clean Code", Author: "Robert C. Martin", Status: "Available"})
	_ = l.AddMember(models.Member{ID: 1, Name: "Alice"})
	_ = l.AddMember(models.Member{ID: 2, Name: "Bob"})

	// Borrow one for demonstration
	err := l.BorrowBook(2, 1)
	if err != nil {
		fmt.Println("seed borrow error:", err)
	}
}

