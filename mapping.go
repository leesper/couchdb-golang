package couchdb

import "errors"

var (
	// ErrSetID for setting ID to document which already has one.
	ErrSetID = errors.New("id can only be set on new documents")
)

// Document represents a document object in database.
type Document struct {
	id  string
	rev string
	ID  string `json:"_id"`  // for json only, call SetID/GetID instead
	Rev string `json:"_rev"` // for json only, call GetRev instead
}

// DocumentWithID returns a new Document with ID.
func DocumentWithID(id string) Document {
	return Document{
		id: id,
	}
}

// SetID sets ID for new document or return error.
func (d *Document) SetID(id string) error {
	if d.id != "" {
		return ErrSetID
	}
	d.id = id
	return nil
}

// GetID returns the document ID.
func (d *Document) GetID() string {
	return d.id
}

// GetRev returns the document revision.
func (d *Document) GetRev() string {
	return d.rev
}

// Store stores the document in specified database.
func Store(db *Database, obj interface{}) error {
	return errors.New("not implemented")
}

// Load loads the document in specified database.
func Load(db *Database, docID string, obj interface{}) error {
	return errors.New("not implemented")
}

// FromMap converts a map into struct.
func FromMap(obj interface{}, docMap map[string]interface{}) error {
	return errors.New("not implemented")
}

// ToMap converts a struct into map.
func ToMap(obj interface{}) map[string]interface{} {
	return nil
}

// ViewField represents a view definition value bound to Document.
type ViewField struct{}
