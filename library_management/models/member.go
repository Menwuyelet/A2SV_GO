package models

type Member struct {
	ID int
	Name string
	BorrowedBooks []Book
}

func (m Member) CheckIfBookIsBorrowedByMember(bookID int) bool {

	for _, book := range m.BorrowedBooks {
		if book.ID == bookID {
			return true
		}
	}

	return false
}

func (m *Member) RemoveBookFromMemberBorrowList(bookID int) {
	if len(m.BorrowedBooks) == 0 {
		return
	}

	idx := -1

	for i, book := range m.BorrowedBooks {
		if book.ID == bookID {
			idx = i
			break
		}
	}

	if idx == -1 {
		return
	}

	m.BorrowedBooks[idx] = m.BorrowedBooks[len(m.BorrowedBooks)-1]
	m.BorrowedBooks = m.BorrowedBooks[:len(m.BorrowedBooks)-1]

}
