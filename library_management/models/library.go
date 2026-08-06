package models

import (
	"errors"
	"fmt"
)

type Library struct {
	books map[int]Book
	members map[int]Member
}

func (l *Library) AddBook(book Book) {
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	delete(l.books, bookID)
}

func (l *Library) AddMember(member Member) {
	l.members[member.ID] = member
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	if _, exists := l.books[bookID]; !exists || l.books[bookID].Status != "Available" {
		return errors.New("The requested book is not available.")
	}

	member, err := l.RetrieveMember(memberID)
	if err != nil {
		return err
	}

	book := l.books[bookID]
	book.Status = "Borrowed"
	l.books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member

	return nil
}

func (l *Library) ReturnBook(bookID, memberID int) error {
	member, err := l.RetrieveMember(memberID)
	if err != nil {
		return err
	}

	book, existsInLibrary := l.books[bookID]

	if !existsInLibrary {
		return errors.New("The book you are trying to return is not in this library.")
	}

	if !member.CheckIfBookIsBorrowedByMember(bookID) {
		return errors.New("The book you are trying to return is not borrowed by the member.")
	}

	book.Status = "Available"
	l.books[bookID] = book

	member.RemoveBookFromMemberBorrowList(bookID)
	l.members[memberID] = member

	return nil
}

func (l Library) ListAvailableBooks() []Book {
	books := []Book{}

	for _, book := range l.books {
		if book.Status == "Available" {
			books = append(books, book)
		}
	}

	return books
}

func (l Library) ListBorrowedBooks(memberID int) []Book {

	member, err := l.RetrieveMember(memberID)
	if err != nil {
		fmt.Println(err)
		return []Book{}
	}
	books := l.ListBorrowedBooksForMember(member)

	return books
}

func (l Library) RetrieveMember(memberID int) (Member, error) {
	member, ok := l.members[memberID]
	if !ok {
		return Member{}, errors.New("There is no member with this id.")
	}
	
	return member, nil
}

func (l Library) ListBorrowedBooksForMember(member Member) []Book {
	return member.BorrowedBooks
}

var LibraryInstance = &Library{
	books:   make(map[int]Book),
	members: make(map[int]Member),
}
