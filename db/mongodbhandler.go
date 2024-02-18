package db

import (
	"context"
	"dockerrestapi/restlib"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	dataBaseName   = "WanShiTong"
	collectionName = "library"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// CreateMongoDBHandler returns a db interface to a mongo handler given access dsn.
func CreateMongoDBHandler(dsn string) (RestDbInterface, error) {
	clientOptions := options.Client().ApplyURI(dsn)

	log.Println("connecting...")

	// Connect to MockDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("cant connect")
		return nil, err
	}

	log.Println("pinging...")

	//Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Println("cant ping")
		return nil, err
	}

	log.Println("Connected to mongoDB!")
	return &MongoDB{
		client:     client,
		collection: client.Database(dataBaseName).Collection(collectionName),
	}, nil
}

func (m *MongoDB) Disconnect() {
	err := m.client.Disconnect(context.Background())
	if err != nil {
		log.Println("disconnect error:", err.Error())
	}
}

// GetAllBooks , retrieves all books in mongo, returns only identifiers
func (m *MongoDB) GetAllBooks() ([]restlib.BookIdentifier, error) {
	// Find documents with specific fields

	// Specify the fields to include (1) or exclude (0)
	projection := bson.M{restlib.JsonBsonTagName: 1, restlib.JsonBsonTagAuthor: 1}

	cursor, err := m.collection.Find(context.Background(), bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var result []restlib.BookIdentifier
	err = cursor.All(context.Background(), &result)
	if err != nil {
		log.Println(result)
		log.Println(err)
		return nil, err
	}
	return result, nil
}

// GetOneBook retrieves single book given a book Identifier
func (m *MongoDB) GetOneBook(bookIdentifier *restlib.BookIdentifier) (*restlib.Book, error) {
	// Get a single person by ID from the database

	match := bson.M{restlib.JsonBsonTagName: bookIdentifier.Name, restlib.JsonBsonTagAuthor: bookIdentifier.Author}
	receivedBook := &restlib.Book{}
	err := m.collection.FindOne(context.Background(), match).Decode(receivedBook)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, restlib.NoMatchingBook
		}
		return nil, err
	}

	return receivedBook, nil
}

// CreateNewBook stores a new book in db
func (m *MongoDB) CreateNewBook(book *restlib.Book) error {
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

	_, err = m.collection.InsertOne(context.Background(), book)
	if err != nil {
		log.Println("cant insert document", err.Error())
		return err
	}
	return nil
}

// UpdateExistingBook updates existing  book in db
func (m *MongoDB) UpdateExistingBook(book *restlib.Book) error {

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
	book.UpdatedDate = time.Now().Format(restlib.DbTimeFormat)

	_, err = m.collection.UpdateOne(
		context.Background(),
		bson.M{restlib.JsonBsonTagName: book.Name, restlib.JsonBsonTagAuthor: book.Author},
		bson.M{"$set": book},
	)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return restlib.NoMatchingBook
		}
		return err
	}
	return nil
}

// DeleteBook deletes existing book given Identifier
func (m *MongoDB) DeleteBook(bookIdentifier *restlib.BookIdentifier) error {
	// Delete a person by ID from the database

	inDb, err := m.isBookInDb(restlib.BookIdentifier{
		Name:   bookIdentifier.Name,
		Author: bookIdentifier.Author,
	})
	if err != nil {
		return err
	}
	if !inDb {
		return restlib.NoMatchingBook
	}

	match := bson.M{restlib.JsonBsonTagName: bookIdentifier.Name, restlib.JsonBsonTagAuthor: bookIdentifier.Author}
	result, err := m.collection.DeleteOne(context.Background(), match)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return restlib.NoMatchingBook
		}
		return err
	}
	log.Println("deleted no.", result.DeletedCount)

	return nil
}

// checks to see if book Is in DB
func (m *MongoDB) isBookInDb(book restlib.BookIdentifier) (bool, error) {
	projection := bson.M{restlib.JsonBsonTagName: 1, restlib.JsonBsonTagAuthor: 1}
	match := bson.M{restlib.JsonBsonTagName: book.Name, restlib.JsonBsonTagAuthor: book.Author}
	cursor := m.collection.FindOne(context.Background(), match, options.FindOne().SetProjection(projection))
	if cursor.Err() != nil {
		if errors.Is(cursor.Err(), mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, cursor.Err()
	}

	return true, nil
}
