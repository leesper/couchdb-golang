package couchdb

// MapField represents a field for map values in Document.
type MapField map[string]interface{}

// TextField represents a field for string values in Document.
type TextField string

// DateTimeField represents a field for date time values in Document.
type DateTimeField struct{}

// ListField represents a field for list values in Document.
type ListField struct{}

// DecimalField represents a field for decimal values in Document.
type DecimalField float64

// Document represents a document object in database.
type Document struct {
	outer interface{} // the outer struct that contains Document
}

// Store stores the document in specified database.
func (d *Document) Store(db *Database) {}

// Load returns the document in specified database.
func (d *Document) Load(db *Database, docID string) {}

// Wrap converts a map into Document object
func (d *Document) Wrap(docMap map[string]interface{}) {}

func store(doc interface{}, db *Database) {}
func load(db *Database, docID string)     {}

// ViewField represents a view definition value bound to Document.
type ViewField struct{}
