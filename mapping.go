package couchdb

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	// ErrSetID for setting ID to document which already has one.
	ErrSetID = errors.New("id can only be set on new documents")
	// ErrNotStruct for not a struct value
	ErrNotStruct = errors.New("value not of struct type")
	// ErrNotDocumentEmbedded for not a document-embedded value
	ErrNotDocumentEmbedded = errors.New("value not Document-embedded")
)

// Document represents a document object in database.
type Document struct {
	id  string
	rev string
	ID  string `json:"_id,omitempty"`  // for json only, call SetID/GetID instead
	Rev string `json:"_rev,omitempty"` // for json only, call GetRev instead
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

// SetRev sets revision for document.
func (d *Document) SetRev(rev string) {
	d.rev = rev
}

// GetRev returns the document revision.
func (d *Document) GetRev() string {
	return d.rev
}

// Store stores the document in specified database.
func Store(db *Database, obj interface{}) error {
	ptrValue := reflect.ValueOf(obj)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}

	zero := reflect.Value{}
	if ptrValue.Elem().FieldByName("Document") == zero {
		return ErrNotDocumentEmbedded
	}

	jsonIDField := ptrValue.Elem().FieldByName("ID")
	getIDMethod := ptrValue.MethodByName("GetID")

	idStr := getIDMethod.Call([]reflect.Value{})[0].Interface().(string)
	if idStr != "" {
		jsonIDField.SetString(idStr)
	}

	jsonRevField := ptrValue.Elem().FieldByName("Rev")
	getRevMethod := ptrValue.MethodByName("GetRev")
	revStr := getRevMethod.Call([]reflect.Value{})[0].Interface().(string)
	if revStr != "" {
		jsonRevField.SetString(revStr)
	}

	doc, err := ToJSONCompatibleMap(ptrValue.Elem().Interface())
	if err != nil {
		return err
	}

	id, rev, err := db.Save(doc, nil)
	if err != nil {
		return err
	}

	setIDMethod := ptrValue.MethodByName("SetID")
	setRevMethod := ptrValue.MethodByName("SetRev")

	if idStr == "" {
		setIDMethod.Call([]reflect.Value{reflect.ValueOf(id)})
	}

	setRevMethod.Call([]reflect.Value{reflect.ValueOf(rev)})

	return nil
}

// Load loads the document in specified database.
func Load(db *Database, docID string, obj interface{}) error {
	ptrValue := reflect.ValueOf(obj)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}

	zero := reflect.Value{}
	if ptrValue.Elem().FieldByName("Document") == zero {
		return ErrNotDocumentEmbedded
	}

	doc, err := db.Get(docID, nil)
	if err != nil {
		return err
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return nil
}

// FromJSONCompatibleMap constructs a struct from a JSON-compatible map.
func FromJSONCompatibleMap(obj interface{}, docMap map[string]interface{}) error {
	return errors.New("not implemented")
}

// ToJSONCompatibleMap converts a struct into a JSON-compatible map, e.g. anything that cannot
// be jsonified will be ignored silently.
func ToJSONCompatibleMap(obj interface{}) (map[string]interface{}, error) {
	structValue := reflect.ValueOf(obj)
	if structValue.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	doc := map[string]interface{}{}
	err = json.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// ViewField represents a view definition value bound to Document.
type ViewField struct{}
