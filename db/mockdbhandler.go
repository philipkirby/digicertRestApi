package db

import (
	"dockerrestapi/restlib"
	"log"
	"time"
)

type MockDB struct {
	db map[restlib.BookIdentifier]restlib.Book
}

func CreateMockDBHandler() (RestDbInterface, error) {
	log.Println("Connected to MockDB!")
	return &MockDB{
		db: map[restlib.BookIdentifier]restlib.Book{},
	}, nil
}

func (m *MockDB) Disconnect() {
}

func (m *MockDB) GetAllBooks() ([]restlib.BookIdentifier, error) {

	var results []restlib.BookIdentifier
	for identifier, _ := range m.db {
		results = append(results, identifier)
	}
	return results, nil
}

func (m *MockDB) GetOneBook(bookIdentifier *restlib.BookIdentifier) (*restlib.Book, error) {
	for dbIdentifier, foundBook := range m.db {
		if dbIdentifier.Name == bookIdentifier.Name && dbIdentifier.Author == bookIdentifier.Author {
			return &foundBook, nil
		}
	}
	return nil, restlib.NoMatchingBook
}

func (m *MockDB) CreateNewBook(book *restlib.Book) error {
	// Create a new person and insert into the database

	inDb, err := m.isBookInDb(restlib.BookIdentifier{
		Name:   book.Name,
		Author: book.Author,
	})
	if err != nil {
		return err
	}
	if inDb {
		return restlib.BookAlreadyExists
	}

	book.UpdatedDate = time.Now().Format(restlib.DbTimeFormat)
	m.db[restlib.BookIdentifier{
		Name:   book.Name,
		Author: book.Author,
	}] = *book

	return nil
}

func (m *MockDB) UpdateExistingBook(book *restlib.Book) error {
	// Update a person by ID in the database

	inDb, err := m.isBookInDb(restlib.BookIdentifier{
		Name:   book.Name,
		Author: book.Author,
	})
	if err != nil {
		return err
	}
	if !inDb {
		return restlib.NoMatchingBook
	}

	m.db[restlib.BookIdentifier{
		Name:   book.Name,
		Author: book.Author,
	}] = *book

	return nil
}

func (m *MockDB) DeleteBook(bookIdentifier *restlib.BookIdentifier) error {
	// Delete a person by ID from the database
	inDb, err := m.isBookInDb(*bookIdentifier)
	if err != nil {
		return restlib.NoMatchingBook
	}
	if !inDb {
		return restlib.NoMatchingBook
	}

	delete(m.db, *bookIdentifier)

	return nil
}

func (m *MockDB) isBookInDb(bookIdentifier restlib.BookIdentifier) (bool, error) {
	_, exists := m.db[bookIdentifier]
	return exists, nil
}
