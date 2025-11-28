package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

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
	ReserveBook(bookID int, memberID int) error
}

// Library implements LibraryManager with concurrency support.
type Library struct {
	books        map[int]models.Book
	members      map[int]models.Member
	reservations map[int]int          // bookID -> memberID
	timers       map[int]*time.Timer  // bookID -> auto-cancel timer
	mu           sync.Mutex
}

// NewLibrary creates a new Library instance.
func NewLibrary() *Library {
	return &Library{
		books:        make(map[int]models.Book),
		members:      make(map[int]models.Member),
		reservations: make(map[int]int),
		timers:       make(map[int]*time.Timer),
	}
}

// AddBook adds a new book to the library.
func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if book.Status == "" {
		book.Status = "Available"
	}
	l.books[book.ID] = book
}

// RemoveBook removes a book from the library by its ID.
func (l *Library) RemoveBook(bookID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	if b.Status == "Borrowed" {
		return errors.New("cannot remove a borrowed book")
	}
	if _, reserved := l.reservations[bookID]; reserved {
		return errors.New("cannot remove a reserved book")
	}
	delete(l.books, bookID)
	return nil
}

// AddMember adds a new member to the library.
func (l *Library) AddMember(m models.Member) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.members[m.ID]; exists {
		return errors.New("member with this ID already exists")
	}
	m.BorrowedBooks = []models.Book{}
	l.members[m.ID] = m
	return nil
}

// GetMember returns a pointer to a member if exists.
func (l *Library) GetMember(memberID int) (*models.Member, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	m, ok := l.members[memberID]
	if !ok {
		return nil, errors.New("member not found")
	}
	// Return a copy to avoid exposing internal state
	copyM := m
	return &copyM, nil
}

// BorrowBook allows a member to borrow a book if it is available or reserved by them.
func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	// if already borrowed
	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}

	// if reserved by someone else
	if reserver, isReserved := l.reservations[bookID]; isReserved && reserver != memberID {
		return errors.New("book reserved by another member")
	}

	// mark book as borrowed
	book.Status = "Borrowed"
	book.ReservedBy = 0
	book.ReservedAt = time.Time{}
	l.books[bookID] = book

	// cancel auto-cancel timer if exists
	if timer, exists := l.timers[bookID]; exists {
		timer.Stop()
		delete(l.timers, bookID)
	}
	// remove reservation entry
	delete(l.reservations, bookID)

	// attach to member
	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member

	return nil
}

// ReturnBook allows a member to return a borrowed book.
func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	// check that member has borrowed the book
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

	// remove from member
	member.BorrowedBooks = append(member.BorrowedBooks[:foundIdx], member.BorrowedBooks[foundIdx+1:]...)
	l.members[memberID] = member

	// update book to available (note: not reserved)
	book.Status = "Available"
	l.books[bookID] = book

	return nil
}

// ListAvailableBooks lists all available books.
func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	list := make([]models.Book, 0, len(l.books))
	for _, b := range l.books {
		if b.Status == "Available" {
			list = append(list, b)
		}
	}
	return list
}

// ListBorrowedBooks lists all books borrowed by a specific member.
func (l *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	member, ok := l.members[memberID]
	if !ok {
		return nil, errors.New("member not found")
	}
	// return a copy
	copySlice := make([]models.Book, len(member.BorrowedBooks))
	copy(copySlice, member.BorrowedBooks)
	return copySlice, nil
}

// ReserveBook reserves a book for a member. If reserved, it returns error.
// A timer is scheduled to auto-cancel the reservation after 5 seconds if not borrowed.
func (l *Library) ReserveBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	// cannot reserve if already borrowed
	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}
	// cannot reserve if already reserved
	if reserver, exists := l.reservations[bookID]; exists {
		return fmt.Errorf("book already reserved by member %d", reserver)
	}

	// mark reservation
	l.reservations[bookID] = memberID
	book.ReservedBy = memberID
	book.ReservedAt = time.Now()
	l.books[bookID] = book

	// schedule auto-cancel in 5 seconds
	timer := time.AfterFunc(5*time.Second, func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		// only cancel if still reserved and not borrowed
		if reserver, exists := l.reservations[bookID]; exists && reserver == memberID {
			// double-check book status
			if b, bok := l.books[bookID]; bok && b.Status != "Borrowed" {
				delete(l.reservations, bookID)
				if t, ok := l.timers[bookID]; ok {
					t.Stop()
				}
				delete(l.timers, bookID)
				// clear reservation metadata
				b.ReservedBy = 0
				b.ReservedAt = time.Time{}
				l.books[bookID] = b
				fmt.Printf("[AUTO-CANCEL] Reservation for book %d auto-cancelled (member %d)\n", bookID, memberID)
			}
		}
	})

	// store timer and return
	l.timers[bookID] = timer
	return nil
}

// SeedSampleData seeds the library with sample data.
func (l *Library) SeedSampleData() {
	l.AddBook(models.Book{ID: 1, Title: "1984", Author: "George Orwell"})
	l.AddBook(models.Book{ID: 2, Title: "The Hobbit", Author: "J.R.R. Tolkien"})
	l.AddBook(models.Book{ID: 3, Title: "Clean Code", Author: "Robert C. Martin"})
	_ = l.AddMember(models.Member{ID: 1, Name: "Alice"})
	_ = l.AddMember(models.Member{ID: 2, Name: "Bob"})
	_ = l.AddMember(models.Member{ID: 3, Name: "Carol"})
}
