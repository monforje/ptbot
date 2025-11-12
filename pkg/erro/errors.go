package erro

import "errors"

var (
	ErrDocumentExists   = errors.New("document already exists")
	ErrDocumentNotFound = errors.New("document not found")
)
