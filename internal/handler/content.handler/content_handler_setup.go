package contenthandler

import (
	"database/sql"

	contentservice "github.com/ArdiSasongko/go-forum-backend/internal/service/content.service"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/content"
)

type contentHandler struct {
	service contentservice.ContentService
}

func NewContentHandler(service contentservice.ContentService) *contentHandler {
	return &contentHandler{service: service}
}

var db *sql.DB
var queries = contentservice.Queries{
	ContentQueries: content.New(db),
}
