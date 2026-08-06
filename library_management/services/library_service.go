package services

import (
	"library_management/models"
)

func AddBookToLibrary(ID int, Title string, Author string) {
	book := models.Book{
		ID:     ID,
		Title:  Title,
		Author: Author,
		Status: "Available",
	}

	models.LibraryInstance.AddBook(book)
}

func RemoveBookFromLibrary(bookId int) {
	models.LibraryInstance.RemoveBook(bookId)
}

func BorrowBookForMember(bookId, memberId int) string {
	borrow := models.LibraryInstance.BorrowBook(bookId, memberId)

	if borrow != nil {
		return borrow.Error()
	}

	return "Successfully borrowed the book."
}

func ReturnBookForMember(bookId, memberID int) string {
	ret := models.LibraryInstance.ReturnBook(bookId, memberID)

	if ret != nil {
		return ret.Error()
	}

	return "Successfully returned the book."
}

func ListAllAvailableBooks() []models.Book {
	books := models.LibraryInstance.ListAvailableBooks()
	return books
}

func ListBorrowedBooksByMember(memberID int) ([]models.Book, error) {
	member, err := models.LibraryInstance.RetrieveMember(memberID)

	if err != nil {
		return []models.Book{}, err
	}

	books := models.LibraryInstance.ListBorrowedBooksForMember(member)

	return books, nil
}
