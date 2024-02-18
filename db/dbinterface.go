package db

import (
	"dockerrestapi/restlib"
)

type RestDbInterface interface {
	Disconnect()
	GetAllBooks() ([]restlib.BookIdentifier, error)
	CreateNewBook(book *restlib.Book) error
	GetOneBook(bookIdentifier *restlib.BookIdentifier) (*restlib.Book, error)
	UpdateExistingBook(book *restlib.Book) error
	DeleteBook(bookIdentifier *restlib.BookIdentifier) error
}
