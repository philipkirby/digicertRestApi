package db

import (
	"dockerrestapi/lib"
)

// RestDbInterface built for book library.
type RestDbInterface interface {
	Disconnect()
	GetAllBooks() ([]lib.BookIdentifier, error)
	CreateNewBook(book *lib.Book) error
	GetOneBook(bookIdentifier *lib.BookIdentifier) (*lib.Book, error)
	UpdateExistingBook(book *lib.Book) error
	DeleteBook(bookIdentifier *lib.BookIdentifier) error
}
