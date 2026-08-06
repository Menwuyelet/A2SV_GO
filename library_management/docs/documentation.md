# Library Management System — Documentation

## Overview

A simple console-based library management system written in Go. It demonstrates core
Go concepts: **structs**, **interfaces**, **methods**, **slices**, and **maps**.

The application lets a user:
- Add and remove books
- Borrow and return books
- List available books
- List all books borrowed by a member

## Folder Structure

```
library_management/
├── main.go                              # Entry point
├── controllers/
│   └── library_controller.go            # Console menu + input handling
├── interfaces/
│   └── library.go                       # LibraryManager interface
├── models/
│   ├── book.go                          # Book struct
│   ├── member.go                        # Member struct + helpers
│   └── library.go                       # Library struct (implements interface)
├── services/
│   ├── library_service.go               # Book/borrow/return business logic
│   └── member_service.go                # Member-specific queries
├── docs/
│   └── documentation.md                 # This file
└── go.mod
```

## How to Run

```bash
go run .
```

You will see a numbered menu:

```
1. List available books.
2. Remove an existing book.
3. Add new book.
4. Borrow a book.
5. Return a book.
6. List all borrowed books for a member.
0. Exit.
```

## Data Models

### Book (`models/book.go`)

| Field  | Type   | Description                   |
|--------|--------|-------------------------------|
| ID     | int    | Unique book identifier        |
| Title  | string | Book title                    |
| Author | string | Book author                   |
| Status | string | `"Available"` or `"Borrowed"` |

### Member (`models/member.go`)

| Field         | Type     | Description                     |
|---------------|----------|---------------------------------|
| ID            | int      | Unique member identifier        |
| Name          | string   | Member name                     |
| BorrowedBooks | `[]Book` | Slice of books the member holds |

Helper methods:
- `CheckIfBookIsBorrowedByMember(bookID int) bool`
- `RemoveBookFromMemberBorrowList(bookID int)`

## LibraryManager Interface (`interfaces/library.go`)

```go
type LibraryManager interface {
    AddBook(book models.Book)
    RemoveBook(bookID int)
    BorrowBook(bookID int, memberID int) error
    ReturnBook(bookID int, memberID int) error
    ListAvailableBooks() []models.Book
    ListBorrowedBooks(memberID int) []models.Book
}
```

## Library Implementation (`models/library.go`)

`Library` stores books and members in maps keyed by ID and implements `LibraryManager`:

- `AddBook` — inserts the book into the `books` map.
- `RemoveBook` — deletes the book by ID.
- `BorrowBook` — returns an error if the book is missing/not available or the member
  does not exist; otherwise marks the book `Borrowed` and appends it to the member's list.
- `ReturnBook` — returns an error if the member/book is missing or the book was not
  borrowed by that member; otherwise marks the book `Available` and removes it from the
  member's list.
- `ListAvailableBooks` — returns only books with `Status == "Available"`.
- `ListBorrowedBooks` — returns the member's borrowed books; prints an error for an
  unknown member.

A package-level singleton `LibraryInstance` holds the shared state.

## Error Handling

Errors are surfaced in these scenarios:
- Book not found or already borrowed when borrowing
- Member not found when borrowing, returning, or listing borrowed books
- Book not borrowed by the member when returning
- Invalid console input (non-numeric IDs)

## Services Layer

`services/library_service.go` and `services/member_service.go` wrap the model methods
and translate errors into user-friendly messages returned to the controller.
