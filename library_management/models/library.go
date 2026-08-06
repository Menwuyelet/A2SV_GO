package models

import (
	"errors"
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

func (l *Library) BorrowBook(bookID int, memberID int) error {
	if _, ok := l.books[bookID]; !ok || l.books[bookID].Status != "Available" {
		return errors.New("The book you'r looking for is not available.")
	}

	book := l.books[bookID]
	book.Status = "Borrowed"

	member, ok := l.RetrieveMember(memberID)

	if ok == nil{
		return errors.New("There is no member with the provided id")
	}

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	
	return nil
}

func (l *Library) ReturnBook(bookID, memberID int) error {
	member, ok := l.RetrieveMember(memberID)

	if ok == nil {
		return errors.New("There is no member with provided id")
	}

	_, existsInLibrary := l.books[bookID]
	borrowedByMember := member.CheckIfBookIsBorrowedByMember(bookID)

	if !existsInLibrary{
		return errors.New("The book you're trying to return is not borrowed from this library.")
	}

	if borrowedByMember {
		return errors.New("The book you're trying to return is not borrowed by the member.")
	}

	book := l.books[bookID]
	book.Status = "Available"

	member.RemoveBookFromMemberBorrowList(bookID)

	return nil
}

func (l Library) ListAvailableBooks() []Book {
	books := []Book{}

	for _, book := range l.books {
		books = append(books, book)
	}
	return books
}

func (l Library) ListBorrowedBooks(memberID int) ([]Book, error) {

	member, ok := l.RetrieveMember(memberID)
	
	if ok != nil {
		return []Book{}, errors.New("There is no member with given id")
	}

	borrowedBooks := l.ListBorrowedBooksForMember(member)

	books := []Book{}
	for _,book := range borrowedBooks {
		books = append(books, book)
	}
	return books, nil
}

func (l Library) RetrieveMember(memberID int) (Member, error) {
	members := []Member{}
	
	for _, member := range l.members{
		if memberID == member.ID{
			return member, nil
		}
		members = append(members, member)
	}
	member := Member{}
	return member, errors.New("There is no member with this id.")
}

func (l Library) ListBorrowedBooksForMember(member Member)[]Book {
	return member.BorrowedBooks
}

var LibraryInstance = &Library{}

