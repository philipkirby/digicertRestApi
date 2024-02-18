package internal

import (
	"bytes"
	"dockerrestapi/db"
	"dockerrestapi/restlib"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var api *RestService

var (
	defaultBook1 = restlib.Book{
		Name:     "book1",
		Author:   "philip",
		Contents: "A bad read",
	}
	defaultBook1Updated = restlib.Book{
		Name:     "book1",
		Author:   "philip",
		Contents: "A good read",
	}
	defaultBook2 = restlib.Book{
		Name:     "book2",
		Author:   "Gino",
		Contents: "A wild read",
	}
	defaultBook3 = restlib.Book{
		Name:     "book3",
		Author:   "Sheldon",
		Contents: "A bizarre read",
	}
	mulFormedBook1 = restlib.Book{
		Name:     "",
		Author:   paramAuthor,
		Contents: "contents",
	}
	mulFormedBook2 = restlib.Book{
		Name:     paramName,
		Author:   "",
		Contents: "contents",
	}
	mulFormedBook3 = restlib.Book{
		Name:     paramName,
		Author:   paramAuthor,
		Contents: "",
	}
)

func TestInitialize(t *testing.T) {
	var err error
	api, err = createMockApi()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateBook(t *testing.T) {
	//// Insert Book1
	marshalDefaultBook1, err := json.Marshal(defaultBook1)
	if err != nil {
		return
	}
	response, err := testResponse(http.MethodPut, createBookPath, api.createBook, marshalDefaultBook1, http.StatusOK, nil)
	if err != nil {
		t.Error(err)
	}
	defaultBook1.UpdatedDate = time.Now().Format(restlib.DbTimeFormat)
	if response != "" {
		t.Error("expecting", "", "got", response)
	}

	//// Fail on inserting book1 again
	response, err = testResponse(http.MethodPut, createBookPath, api.createBook, marshalDefaultBook1, http.StatusBadRequest, nil)
	if err != nil {
		t.Error(err)
	}
	if response != "\""+restlib.BookAlreadyExists.Error()+"\"" {
		t.Error("expecting", "\""+restlib.BookAlreadyExists.Error()+"\"", "got", response)
	}

	//// Insert books 2 and 3
	marshalDefaultBook2, err := json.Marshal(defaultBook2)
	if err != nil {
		return
	}
	response, err = testResponse(http.MethodPut, createBookPath, api.createBook, marshalDefaultBook2, http.StatusOK, nil)
	if err != nil {
		t.Error()
	}
	defaultBook2.UpdatedDate = time.Now().Format(restlib.DbTimeFormat)
	if response != "" {
		t.Error("expecting", "", "got", response)
	}

	marshalDefaultBook3, err := json.Marshal(defaultBook3)
	if err != nil {
		return
	}
	response, err = testResponse(http.MethodPut, createBookPath, api.createBook, marshalDefaultBook3, http.StatusOK, nil)
	if err != nil {
		t.Error()
	}
	defaultBook3.UpdatedDate = time.Now().Format(restlib.DbTimeFormat)
	if response != "" {
		t.Error("expecting", "", "got", response)
	}
}

func TestUpdateBook(t *testing.T) {

	////update book1 with book1 updated
	marshalDefaultBook1, err := json.Marshal(defaultBook1Updated)
	if err != nil {
		return
	}
	response, err := testResponse(http.MethodPut, updateBookPath, api.updateBook, marshalDefaultBook1, http.StatusOK, nil)
	if err != nil {
		t.Error()
	}
	defaultBook1.UpdatedDate = time.Now().Format(restlib.DbTimeFormat)
	if response != "" {
		t.Error("expecting", "", "got", response)
	}

	//// fail on updating book that doesn't exist
	fakeBook := restlib.Book{
		Name:     "NameFake",
		Author:   "AuthorFake",
		Contents: "ContentsFake",
	}
	marshalFakeBook, err := json.Marshal(fakeBook)
	if err != nil {
		return
	}
	response, err = testResponse(http.MethodPut, updateBookPath, api.updateBook, marshalFakeBook, http.StatusBadRequest, nil)
	if err != nil {
		t.Error()
	}
	if response != "\""+restlib.NoMatchingBook.Error()+"\"" {
		t.Error("expecting", "\""+restlib.NoMatchingBook.Error()+"\"", "got", response)
	}
}

func TestGetBook(t *testing.T) {

	//Test Book 1, should be updated
	paramMap := map[string]string{paramName: defaultBook1Updated.Name, paramAuthor: defaultBook1Updated.Author}
	response, err := testResponse(http.MethodGet, getBookPath, api.getBook, nil, http.StatusOK, paramMap)
	if err != nil {
		t.Error(err)
	}
	if !sameBookFromHttpRequest(t, response, defaultBook1Updated) {
		t.Error("not the same")
	}

	// Get book 2
	paramMap = map[string]string{paramName: defaultBook2.Name, paramAuthor: defaultBook2.Author}
	response, err = testResponse(http.MethodGet, getBookPath, api.getBook, nil, http.StatusOK, paramMap)
	if err != nil {
		t.Error(err)
	}
	if !sameBookFromHttpRequest(t, response, defaultBook2) {
		t.Fail()
	}

	// Get book 3
	paramMap = map[string]string{paramName: defaultBook3.Name, paramAuthor: defaultBook3.Author}
	response, err = testResponse(http.MethodGet, getBookPath, api.getBook, nil, http.StatusOK, paramMap)
	if err != nil {
		t.Error()
	}
	if !sameBookFromHttpRequest(t, response, defaultBook3) {
		t.Fail()
	}
}

func TestDeleteBook(t *testing.T) {
	// delete book 3
	paramMap := map[string]string{paramName: defaultBook3.Name, paramAuthor: defaultBook3.Author}
	_, err := testResponse(http.MethodDelete, deleteBookPath, api.deleteBook, nil, http.StatusOK, paramMap)
	if err != nil {
		t.Error()
	}

	// fail on trying to delete book 3 again
	response, err := testResponse(http.MethodDelete, deleteBookPath, api.deleteBook, nil, http.StatusNotFound, paramMap)
	if err != nil {
		t.Error()
	}
	if response != "\""+restlib.NoMatchingBook.Error()+"\"" {
		t.Error("expecting", restlib.NoMatchingBook.Error(), "got", response)
	}

}

func TestGetBooks(t *testing.T) {
	var listOfDefaultBooks []restlib.BookIdentifier

	listOfDefaultBooks = append(listOfDefaultBooks, restlib.BookIdentifier{
		Name:   defaultBook1Updated.Name,
		Author: defaultBook1Updated.Author,
	})
	listOfDefaultBooks = append(listOfDefaultBooks, restlib.BookIdentifier{
		Name:   defaultBook2.Name,
		Author: defaultBook2.Author,
	})

	response, err := testResponse(http.MethodGet, getBooksPath, api.getBooks, nil, http.StatusOK, nil)
	if err != nil {
		t.Error()
	}

	var responseList []restlib.BookIdentifier
	err = json.Unmarshal([]byte(response), &responseList)
	if err != nil {
		t.Error(err)
	}

	for _, defaultBook := range listOfDefaultBooks {
		found := false
		for _, responseBook := range responseList {
			if defaultBook.Name == responseBook.Name && defaultBook.Author == responseBook.Author {
				found = true
			}
		}
		if found {
			continue
		}
		t.Error("didnt find", defaultBook)
	}
}

func createMockApi() (*RestService, error) {
	mockConn, err := db.CreateMockDBHandler()
	if err != nil {
		return nil, err
	}
	service, err := CreateRestApiService(mockConn)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func TestUnmarshalAndValidateStoreBookRequest(t *testing.T) {
	marshal, err := json.Marshal(defaultBook1)
	if err != nil {
		return
	}
	request := httptest.NewRequest(http.MethodPost, "/store/book", bytes.NewBufferString(string(marshal)))
	book, err := api.unmarshalAndValidateStoreBookRequest(request.Body)
	if err != nil {
		t.Error(err)
	}
	if !isEqualBook(*book, defaultBook1) {
		t.Errorf("Expected book %+v, but got %+v", defaultBook1, book)
	}

	// Should fail on all 3 mul formed books
	// mul formed 1
	marshal, err = json.Marshal(mulFormedBook1)
	if err != nil {
		return
	}
	request = httptest.NewRequest(http.MethodPost, "/store/book", bytes.NewBufferString(string(marshal)))
	book, err = api.unmarshalAndValidateStoreBookRequest(request.Body)
	if err == nil {
		t.Error("should of failed here")
	}

	// mul formed 2
	marshal, err = json.Marshal(mulFormedBook2)
	if err != nil {
		return
	}
	request = httptest.NewRequest(http.MethodPost, "/store/book", bytes.NewBufferString(string(marshal)))
	book, err = api.unmarshalAndValidateStoreBookRequest(request.Body)
	if err == nil {
		t.Error("should of failed here")
	}

	// mul formed 3
	marshal, err = json.Marshal(mulFormedBook3)
	if err != nil {
		return
	}
	request = httptest.NewRequest(http.MethodPost, "/store/book", bytes.NewBufferString(string(marshal)))
	book, err = api.unmarshalAndValidateStoreBookRequest(request.Body)
	if err == nil {
		t.Error("should of failed here")
	}

}

func isEqualBook(book1, book2 restlib.Book) bool {
	return book1.Name == book2.Name && book1.Author == book2.Author && book1.Contents == book2.Contents
}

func TestCreateBookIdentifierFromParams(t *testing.T) {

	// should pass this
	bookIdentifier, err := api.createBookIdentifierFromParams(map[string]string{paramName: "The Okay Gatsby", paramAuthor: "philip"})
	if err != nil {
		t.Error(err)
	}
	if bookIdentifier.Name != "The Okay Gatsby" || bookIdentifier.Author != "philip" {
		t.Error("failed to create identifier")
	}

	// should fail here
	_, err = api.createBookIdentifierFromParams(map[string]string{paramAuthor: "philip"})
	if err == nil {
		t.Error("should of failed here")
	}

	_, err = api.createBookIdentifierFromParams(map[string]string{paramName: "The Okay Gatsby"})
	if err == nil {
		t.Error("should of failed here")
	}

	_, err = api.createBookIdentifierFromParams(map[string]string{paramName: "", paramAuthor: ""})
	if err == nil {
		t.Error("should of failed here")
	}

	_, err = api.createBookIdentifierFromParams(map[string]string{paramName: "The Okay Gatsby", paramAuthor: ""})
	if err == nil {
		t.Error("should of failed here")
	}

}

func testResponse(httpMethod, url string, funcCall func(writer http.ResponseWriter,
	request *http.Request), input []byte, expectedStatus int, params map[string]string) (string, error) {

	// Create a request with the desired URL and payload
	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(input))
	if err != nil {
		return "", err
	}
	req = mux.SetURLVars(req, params)

	// Create a response recorder to capture the response
	responseWriter := httptest.NewRecorder()

	// Call the createBook handler directly, passing the response recorder and request
	funcCall(responseWriter, req)

	if status := responseWriter.Code; status != expectedStatus {
		return "", errors.New("wrong http status code:" + strconv.Itoa(status) + " expected:" + strconv.Itoa(expectedStatus))
	}

	return responseWriter.Body.String(), nil
}

func sameBookFromHttpRequest(t *testing.T, response string, matchedBook restlib.Book) bool {
	bookRetrieved := &restlib.Book{}
	err := json.Unmarshal([]byte(response), bookRetrieved)
	if err != nil {
		t.Fail()
		return false
	}

	if bookRetrieved.Name != matchedBook.Name ||
		bookRetrieved.Author != matchedBook.Author ||
		bookRetrieved.Contents != matchedBook.Contents ||
		bookRetrieved.UpdatedDate != matchedBook.UpdatedDate {

		expected, _ := fmt.Printf("%v", matchedBook)
		retrieved, _ := fmt.Printf("%v", bookRetrieved)
		t.Error("expected", expected, "got", retrieved)
		return false
	}
	return true
}
