# Library Management System (Go)

## 🧩 Overview

This is a simple **console-based Library Management System** implemented in **Go**.  
It demonstrates the use of structs, interfaces, slices, maps, and methods in Go.

---

## 📁 Folder Structure

library_management/
├── main.go
├── controllers/
│ └── library_controller.go
├── models/
│ ├── book.go
│ └── member.go
├── services/
│ └── library_service.go
├── docs/
│ └── documentation.md
└── go.mod

---

## 🏗️ Components Description

| Folder           | Description                                                                    |
| ---------------- | ------------------------------------------------------------------------------ |
| **controllers/** | Handles user input/output from the console and calls service layer methods.    |
| **models/**      | Contains data structures such as `Book` and `Member`.                          |
| **services/**    | Contains the business logic and the `LibraryManager` interface implementation. |
| **docs/**        | Documentation files for the project.                                           |
| **main.go**      | Entry point of the application.                                                |

---

## 📘 Features

- Add a new book to the library
- Remove an existing book
- Borrow a book (if available)
- Return a borrowed book
- List all available books
- List all borrowed books by a member

---
