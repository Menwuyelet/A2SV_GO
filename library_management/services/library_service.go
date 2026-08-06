package services

import "library_management/models"

func AddBookToLibrary(ID int, Title string, Author string, Status bool) {
	book := models.Book{
		ID: ID,
		Title: Title,
		Author: Author,
		Status: "Available",
	}

	models.LibraryInstance.AddBook(book)
}

func RemoveBookFromLibrary(bookId int) {
	models.LibraryInstance.RemoveBook(bookId)
}

func BorrowBookForMember(bookId, memberId int) string{
	borrow := models.LibraryInstance.BorrowBook(bookId, memberId)

	if borrow != nil {
		return borrow.Error()
	}

	return "Successfully borrowed the book."
}


func ReturnBookForMember(bookId, memberID int) string{
	ret := models.LibraryInstance.ReturnBook(bookId, memberID)

	if ret != nil{
		return ret.Error()
	}

	return "Successfully returned the book."
}

func ListAllAvailableBooks() []models.Book{
	return models.LibraryInstance.ListAvailableBooks()
}


