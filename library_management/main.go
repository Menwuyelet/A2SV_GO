package main

import (
	"library_management/controllers"
	"library_management/models"
)

func main() {
	seedLibrary()
	controllers.StartMenu()
}

// Seed function to seed initial data for testing
func seedLibrary() {
	models.LibraryInstance.AddBook(models.Book{ID: 1, Title: "Sliding Window Tutorial", Author: "Nafiad", Status: "Available"})
	models.LibraryInstance.AddBook(models.Book{ID: 2, Title: "GO Lang", Author: "Mnu", Status: "Available"})
	models.LibraryInstance.AddBook(models.Book{ID: 3, Title: "The 100", Author: "Jon Doe", Status: "Available"})

	models.LibraryInstance.AddMember(models.Member{ID: 1, Name: "Nafiad"})
	models.LibraryInstance.AddMember(models.Member{ID: 2, Name: "Mnu"})
}
