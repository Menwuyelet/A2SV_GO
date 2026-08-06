package services

import (
	"errors"
	"library_management/models"
)

func ListAllBorrowedBooksByMember(memberID int) ([]models.Book, error) {
	member, ok := models.LibraryInstance.RetrieveMember(memberID)
	
	if ok != nil {
		return []models.Book{}, errors.New("There is no member with given id")
	}
	books := models.LibraryInstance.ListBorrowedBooksForMember(member)

	return books, nil
}