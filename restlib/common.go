package restlib

import (
	"errors"
	"time"
)

const (
	JsonBsonTagName        = "name"
	JsonBsonTagAuthor      = "author"
	JsonBsonTagContents    = "contents"
	JsonBsonTagUpdatedTime = "updateDate"
	DbTimeFormat           = time.UnixDate
)

type BookIdentifier struct {
	Name   string `bson:"name" json:"name,omitempty" `
	Author string `bson:"author" json:"author,omitempty"`
}

type Book struct {
	Name        string `bson:"name" json:"name,omitempty" `
	Author      string `bson:"author" json:"author,omitempty"`
	Contents    string `bson:"contents" json:"contents,omitempty"`
	UpdatedDate string `bson:"updateDate" json:"updateDate,omitempty"`
}

var (
	NoMatchingBook      = errors.New("no matching book in library")
	BookAlreadyExists   = errors.New("book already exists library")
	IncorrectParameters = errors.New("incorrect request parameter")
)

type ErrorResponse struct {
	ReturnCode  int    `json:"returnCode"`
	ErrorDetail string `json:"errorDetail"`
}
