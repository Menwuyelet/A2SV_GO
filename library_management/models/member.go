package models

type Member struct {
	ID int
	Name string
	BorrowedBooks []Book
}

func (m Member) CheckIfBookIsBorrowedByMember(bookID int) bool {

	for borrowedBook := range m.BorrowedBooks {
		if borrowedBook == bookID{
			return true
		}
	}

	return false
}

func (m *Member) RemoveBookFromMemberBorrowList(bookID int){
	idx := 0

	for borrowedBook := range m.BorrowedBooks {
		if borrowedBook == bookID{
			break
		}
		idx += 1
	}
	
	m.BorrowedBooks[idx] = m.BorrowedBooks[len(m.BorrowedBooks)-1]
	m.BorrowedBooks = m.BorrowedBooks[:len(m.BorrowedBooks)-1]

}