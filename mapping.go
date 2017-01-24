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
	zero                   = reflect.Value{}
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
// obj: a Document-embedded struct value, its id and rev will be updated after stored,
// so caller must pass a pointer value.
func Store(db *Database, obj interface{}) error {
	ptrValue := reflect.ValueOf(obj)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}

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
	jsonRevField.SetString(rev)

	return nil
}

// Load loads the document in specified database.
func Load(db *Database, docID string, obj interface{}) error {
	ptrValue := reflect.ValueOf(obj)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}

	if ptrValue.Elem().FieldByName("Document") == zero {
		return ErrNotDocumentEmbedded
	}

	doc, err := db.Get(docID, nil)
	if err != nil {
		return err
	}

	err = FromJSONCompatibleMap(obj, doc)
	if err != nil {
		return err
	}

	if id, ok := doc["_id"]; ok {
		setIDMethod := ptrValue.MethodByName("SetID")
		setIDMethod.Call([]reflect.Value{reflect.ValueOf(id)})
	}

	if rev, ok := doc["_rev"]; ok {
		setRevMethod := ptrValue.MethodByName("SetRev")
		setRevMethod.Call([]reflect.Value{reflect.ValueOf(rev)})
	}

	return nil
}

// FromJSONCompatibleMap constructs a Document-embedded struct from a JSON-compatible map.
func FromJSONCompatibleMap(obj interface{}, docMap map[string]interface{}) error {
	ptrValue := reflect.ValueOf(obj)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}

	if ptrValue.Elem().FieldByName("Document") == zero {
		return ErrNotDocumentEmbedded
	}

	data, err := json.Marshal(docMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	if id, ok := docMap["_id"]; ok {
		setIDMethod := ptrValue.MethodByName("SetID")
		setIDMethod.Call([]reflect.Value{reflect.ValueOf(id)})
	}

	if rev, ok := docMap["_rev"]; ok {
		setRevMethod := ptrValue.MethodByName("SetRev")
		setRevMethod.Call([]reflect.Value{reflect.ValueOf(rev)})
	}

	return nil
}

// ToJSONCompatibleMap converts a Document-embedded struct into a JSON-compatible map,
// e.g. anything that cannot be jsonified will be ignored silently.
func ToJSONCompatibleMap(obj interface{}) (map[string]interface{}, error) {
	structValue := reflect.ValueOf(obj)
	if structValue.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	zero := reflect.Value{}
	if structValue.FieldByName("Document") == zero {
		return nil, ErrNotDocumentEmbedded
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
type ViewField func() (*ViewDefinition, error)

// NewViewField returns a ViewField function.
// design: the name of the design document.
//
// name: the name of the view.
//
// mapFun: the map function code.
//
// reduceFun: the reduce function code(optional).
//
// language: the name of the programming language used, default is javascript.
//
// wrapper: an optional function for processing the result rows after retrieved.
//
// options: view specific options.
func NewViewField(design, name, mapFun, reduceFun, language string, wrapper func(Row) Row, options map[string]interface{}) ViewField {
	f := func() (*ViewDefinition, error) {
		return NewViewDefinition(design, name, mapFun, reduceFun, language, wrapper, options)
	}
	return ViewField(f)
}
