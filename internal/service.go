package internal

import (
	"dockerrestapi/db"
	"dockerrestapi/lib"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

const (
	BasePath       = "/api/library"
	getBooksPath   = BasePath + "/getlist"
	getBookPath    = BasePath + "/get/{" + paramName + "}/{" + paramAuthor + "}"
	createBookPath = BasePath + "/create"
	updateBookPath = BasePath + "/update"
	deleteBookPath = BasePath + "/delete/{" + paramName + "}/{" + paramAuthor + "}"
	paramAuthor    = "author"
	paramName      = "name"
)

type RestService struct {
	db     db.RestDbInterface
	router *mux.Router
	port   string
}

// Start starts rest api
func (r *RestService) Start() {
	go func() {
		err := http.ListenAndServe(":"+r.port, r.router)
		if err != nil {
			panic(err)
		}
	}()
	log.Printf("rest started on port %s\n", r.port)
}

// Stop stops rest api
func (r *RestService) Stop() {
	r.db.Disconnect()
	stdInfo("stopped restapi")
}

// CreateRestApiService creates a restapi given a db interface
func CreateRestApiService(db db.RestDbInterface, port string) (*RestService, error) {

	stdInfo("creating rest api")
	router := mux.NewRouter()
	restAPi := &RestService{
		db:     db,
		router: router,
		port:   port,
	}

	// Define endpoints
	router.HandleFunc(getBooksPath, restAPi.getBooks).Methods(http.MethodGet)
	router.HandleFunc(createBookPath, restAPi.createBook).Methods(http.MethodPut)
	router.HandleFunc(getBookPath, restAPi.getBook).Methods(http.MethodGet)
	router.HandleFunc(updateBookPath, restAPi.updateBook).Methods(http.MethodPut)
	router.HandleFunc(deleteBookPath, restAPi.deleteBook).Methods(http.MethodDelete)
	return restAPi, nil
}

// getBooks Retrieves a full list (name and author) of every book stored in db
func (r *RestService) getBooks(writer http.ResponseWriter, request *http.Request) {
	stdInfo("received get Books request")
	books, err := r.db.GetAllBooks()
	if err != nil {
		stdError(err.Error())
		r.restResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	r.restResponse(writer, http.StatusOK, books)
}

// getBook Retrieves a single book from the db given the name and author in the path.
// eg : api/library/get/{name}/{author}
func (r *RestService) getBook(writer http.ResponseWriter, request *http.Request) {
	stdInfo("received get Book request")
	params := mux.Vars(request)

	bookIdentifier, err := r.createBookIdentifierFromParams(params)
	if err != nil {
		r.restResponse(writer, http.StatusBadRequest, lib.IncorrectParameters)
		return
	}

	returnedBook, err := r.db.GetOneBook(bookIdentifier)
	if err != nil {
		if errors.Is(err, lib.NoMatchingBook) {
			r.restResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		stdError(err.Error())
		r.restResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	r.restResponse(writer, http.StatusOK, *returnedBook)
}

// createBook Creates stores a new book into the db
func (r *RestService) createBook(writer http.ResponseWriter, request *http.Request) {
	stdInfo("received create Books request")

	book, err := r.unmarshalAndValidateStoreBookRequest(request.Body)
	if err != nil {
		r.restResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	err = r.db.CreateNewBook(book)
	if err != nil {
		if errors.Is(lib.BookAlreadyExists, err) {
			r.restResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		stdError(err.Error())
		r.restResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	r.restResponse(writer, http.StatusOK, nil)
}

// updateBook Updates an existing book in the db
func (r *RestService) updateBook(writer http.ResponseWriter, request *http.Request) {
	stdInfo("received update Book request")
	book, err := r.unmarshalAndValidateStoreBookRequest(request.Body)
	if err != nil {
		stdError(err.Error())
		r.restResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	err = r.db.UpdateExistingBook(book)
	if err != nil {
		if errors.Is(err, lib.NoMatchingBook) {
			r.restResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		stdError(err.Error())
		r.restResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	r.restResponse(writer, http.StatusOK, nil)
}

// deleteBook deletes an existing book in the db given the name and author in the path.
// eg : api/library/get/{name}/{author}
func (r *RestService) deleteBook(writer http.ResponseWriter, request *http.Request) {
	stdInfo("received delete Book request")
	params := mux.Vars(request)
	bookIdentifier, err := r.createBookIdentifierFromParams(params)
	if err != nil {
		r.restResponse(writer, http.StatusBadRequest, lib.IncorrectParameters)
		return
	}

	err = r.db.DeleteBook(bookIdentifier)
	if err != nil {
		if errors.Is(err, lib.NoMatchingBook) {
			r.restResponse(writer, http.StatusNotFound, err.Error())
			return
		}
		r.restResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	r.restResponse(writer, http.StatusOK, nil)
}

// restResponse reponds to a rest call given the Writer status and data to write.
// Assumes "Content-Type", "application/json"
func (r *RestService) restResponse(writer http.ResponseWriter, status int, data any) {
	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "application/json")

	if data == nil {
		_, err := writer.Write(nil)
		if err != nil {
			stdInfo("|Error| cant respond " + err.Error())
		}
		return
	}

	responseBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(writer).Encode(err.Error())
		if err != nil {
			stdInfo("|Error| cant respond " + err.Error())
		}
		return
	}
	_, err = writer.Write(responseBytes)
	if err != nil {
		stdInfo("|Error| cant respond " + err.Error())
	}
}

// unmarshalAndValidateStoreBookRequest unmarshal's json from an io body, validates it has no empty fields, and return a pointer to that book
func (r *RestService) unmarshalAndValidateStoreBookRequest(body io.ReadCloser) (*lib.Book, error) {
	book := &lib.Book{}
	err := json.NewDecoder(body).Decode(book)
	if err != nil {
		return nil, err
	}
	if book.Name == "" || book.Author == "" || book.Contents == "" {
		return nil, errors.New("not enough information to store book")
	}
	return book, nil
}

// createBookIdentifierFromParams returns a bookIdentifier object, given a map of parameters that must contain "name" and "author" keys with non-empty values.
func (r *RestService) createBookIdentifierFromParams(params map[string]string) (*lib.BookIdentifier, error) {
	var name, author string
	var exists bool
	if name, exists = params[paramName]; !exists {
		return nil, errors.New("not enough arguments")
	}
	if author, exists = params[paramAuthor]; !exists {
		return nil, errors.New("not enough arguments")
	}
	if name == "" || author == "" {
		return nil, errors.New("empty argument")
	}
	return &lib.BookIdentifier{
		Name:   name,
		Author: author,
	}, nil
}

// stdInfo prints string to standard out , prefix |Info|
func stdInfo(s string) {
	log.Println("|Info|", s)
}

// stdError prints string to standard out , prefix |Error|
func stdError(s string) {
	log.Println("|Error|", s)
}
