package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"library_management/services"
)

// Our user interaction controller 
func StartMenu() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("==================== Welcome to our Library. ====================")

	for {
		fmt.Println()
		fmt.Println("1. List available books.")
		fmt.Println("2. Remove an existing book.")
		fmt.Println("3. Add new book.")
		fmt.Println("4. Borrow a book.")
		fmt.Println("5. Return a book.")
		fmt.Println("6. List all borrowed books for a member.")
		fmt.Println("0. Exit.")
		fmt.Println()
		fmt.Print("=> Choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		// the switch block that controls the flow of the program
		switch choice {
		case "0":
			fmt.Println("Goodbye!")

			return

		case "1":
			// list available books

			books := services.ListAllAvailableBooks()

			if len(books) == 0 {
				fmt.Println("There is no available books in the library.")
			} else {
				for _, book := range books {
					fmt.Printf("=> ID: %v Title: %v Author: %v Status: %v\n", book.ID, book.Title, book.Author, book.Status)
				}
			}

		case "2":
			// remove specific book

			fmt.Print("Enter the book ID to remove: ")

			bookID, err := readInt(reader)
			if err != nil {
				fmt.Println("Invalid book ID.")
				return
			}

			services.RemoveBookFromLibrary(bookID)

			fmt.Println("Book removed successfully.")

		case "3":
			// Add new book

			var ID int
			fmt.Print("Enter Book ID: ")
			ID, err := readInt(reader)

			if err != nil {
				fmt.Println("Invalid book ID.")
				return
			}

			fmt.Print("Enter Book Title: ")
			Title, _ := reader.ReadString('\n')
			Title = strings.TrimSpace(Title)

			fmt.Print("Enter Book Author: ")
			Author, _ := reader.ReadString('\n')
			Author = strings.TrimSpace(Author)

			services.AddBookToLibrary(ID, Title, Author)

			fmt.Println("Book added successfully.")

		case "4":
			// barrow a book
			fmt.Print("Enter Book ID: ")
			bookID, err := readInt(reader)
			if err != nil {
				fmt.Println("Invalid input.")
				return
			}

			fmt.Print("Enter Member ID: ")
			memberID, err := readInt(reader)
			if err != nil {
				fmt.Println("Invalid input.")
				return
			}

			fmt.Println(services.BorrowBookForMember(bookID, memberID))

		case "5":
			// return a book
			fmt.Print("Enter Book ID: ")
			bookID, err := readInt(reader)
			if err != nil {
				fmt.Println("Invalid input.")
				return
			}

			fmt.Print(("Enter Member ID: "))
			memberID, err := readInt(reader)
			if err != nil {
				fmt.Println("Invalid input.")
				return
			}

			fmt.Println(services.ReturnBookForMember(bookID, memberID))

		case "6":
			// list borrowed books
			fmt.Print("Enter Member ID: ")
			memberID, err := readInt(reader)

			if err != nil {
				fmt.Println("Invalid member ID.")
				return
			}

			books, err := services.ListBorrowedBooksByMember(memberID)

			if err != nil {
				fmt.Println(err)
				return
			}

			if len(books) == 0 {
				fmt.Println("This member has not borrowed any books.")
			} else {
				fmt.Println("Borrowed books:")
				for _, book := range books {
					fmt.Printf("=> ID: %v Title: %v Author: %v Status: %v\n", book.ID, book.Title, book.Author, book.Status)
				}
			}

		default:
			fmt.Println("Invalid choice.")
		}
	}
}


// helper function to read integer input from terminal
func readInt(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return 0, err
	}

	return value, nil
}
