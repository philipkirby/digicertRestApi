package db

import (
	"dockerrestapi/lib"
	"log"
	"time"
)

type MockDB struct {
	db map[lib.BookIdentifier]lib.Book
}

func CreateMockDBHandler() (RestDbInterface, error) {
	log.Println("Connected to MockDB!")
	return &MockDB{
		db: map[lib.BookIdentifier]lib.Book{},
	}, nil
}

func (m *MockDB) Disconnect() {
}

func (m *MockDB) GetAllBooks() ([]lib.BookIdentifier, error) {

	var results []lib.BookIdentifier
	for identifier, _ := range m.db {
		results = append(results, identifier)
	}
	return results, nil
}

func (m *MockDB) GetOneBook(bookIdentifier *lib.BookIdentifier) (*lib.Book, error) {
	for dbIdentifier, foundBook := range m.db {
		if dbIdentifier.Name == bookIdentifier.Name && dbIdentifier.Author == bookIdentifier.Author {
			return &foundBook, nil
		}
	}
	return nil, lib.NoMatchingBook
}

func (m *MockDB) CreateNewBook(book *lib.Book) error {
	// Create a new person and insert into the database

	inDb, err := m.isBookInDb(lib.BookIdentifier{
		Name:   book.Name,
		Author: book.Author,
	})
	if err != nil {
		return err
	}
	if inDb {
		return lib.BookAlreadyExists
	}

	book.UpdatedDate = time.Now().Format(lib.DbTimeFormat)
	m.db[lib.BookIdentifier{
		Name:   book.Name,
		Author: book.Author,
	}] = *book

	return nil
}

func (m *MockDB) UpdateExistingBook(book *lib.Book) error {
	// Update a person by ID in the database

	inDb, err := m.isBookInDb(lib.BookIdentifier{
		Name:   book.Name,
		Author: book.Author,
	})
	if err != nil {
		return err
	}
	if !inDb {
		return lib.NoMatchingBook
	}

	m.db[lib.BookIdentifier{
		Name:   book.Name,
		Author: book.Author,
	}] = *book

	return nil
}

func (m *MockDB) DeleteBook(bookIdentifier *lib.BookIdentifier) error {
	// Delete a person by ID from the database
	inDb, err := m.isBookInDb(*bookIdentifier)
	if err != nil {
		return lib.NoMatchingBook
	}
	if !inDb {
		return lib.NoMatchingBook
	}

	delete(m.db, *bookIdentifier)

	return nil
}

func (m *MockDB) isBookInDb(bookIdentifier lib.BookIdentifier) (bool, error) {
	_, exists := m.db[bookIdentifier]
	return exists, nil
}
