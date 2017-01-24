package couchdb

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	// DefaultBaseURL is the default address of CouchDB server.
	DefaultBaseURL = "http://localhost:5984"
)

var (
	// ErrBatchValue for invalid batch parameter of IterView
	ErrBatchValue = errors.New("batch must be 1 or more")
	// ErrLimitValue for invalid limit parameter of IterView
	ErrLimitValue = errors.New("limit must be 1 or more")
)

// getDefaultCouchDBURL returns the default CouchDB server url.
func getDefaultCouchDBURL() string {
	var couchdbURLEnviron string
	for _, couchdbURLEnviron = range os.Environ() {
		if strings.HasPrefix(couchdbURLEnviron, "COUCHDB_URL") {
			break
		}
	}
	if len(couchdbURLEnviron) == 0 {
		couchdbURLEnviron = DefaultBaseURL
	} else {
		couchdbURLEnviron = strings.Split(couchdbURLEnviron, "=")[1]
	}
	return couchdbURLEnviron
}

// Database represents a CouchDB database instance.
type Database struct {
	resource *Resource
}

// NewDatabase returns a CouchDB database instance.
func NewDatabase(urlStr string) (*Database, error) {
	var dbURLStr string
	if !strings.HasPrefix(urlStr, "http") {
		base, err := url.Parse(getDefaultCouchDBURL())
		if err != nil {
			return nil, err
		}
		dbURL, err := base.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		dbURLStr = dbURL.String()
	} else {
		dbURLStr = urlStr
	}

	res, err := NewResource(dbURLStr, nil)
	if err != nil {
		return nil, err
	}

	return newDatabase(res)
}

// NewDatabaseWithResource returns a CouchDB database instance with resource obj.
func NewDatabaseWithResource(res *Resource) (*Database, error) {
	return newDatabase(res)
}

func newDatabase(res *Resource) (*Database, error) {
	return &Database{
		resource: res,
	}, nil
}

// Available returns error if the database is not good to go.
func (d *Database) Available() error {
	_, _, err := d.resource.Head("", nil, nil)
	return err
}

// Save creates a new document or update an existing document.
// If doc has no _id the server will generate a random UUID and a new document will be created.
// Otherwise the doc's _id will be used to identify the document to create or update.
// Trying to update an existing document with an incorrect _rev will cause failure.
// *NOTE* It is recommended to avoid saving doc without _id and instead generate document ID on client side.
// To avoid such problems you can generate a UUID on the client side.
// GenerateUUID provides a simple, platform-independent implementation.
// You can also use other third-party packages instead.
// doc: the document to create or update.
func (d *Database) Save(doc map[string]interface{}, options url.Values) (string, string, error) {
	var id, rev string

	var httpFunc func(string, http.Header, map[string]interface{}, url.Values) (http.Header, []byte, error)
	if v, ok := doc["_id"]; ok {
		httpFunc = docResource(d.resource, v.(string)).PutJSON
	} else {
		httpFunc = d.resource.PostJSON
	}

	_, data, err := httpFunc("", nil, doc, options)
	if err != nil {
		return id, rev, err
	}

	var jsonMap map[string]interface{}
	jsonMap, err = parseData(data)
	if err != nil {
		return id, rev, err
	}

	if v, ok := jsonMap["id"]; ok {
		id = v.(string)
		doc["_id"] = id
	}

	if v, ok := jsonMap["rev"]; ok {
		rev = v.(string)
		doc["_rev"] = rev
	}

	return id, rev, nil
}

// Get returns the document with the specified ID.
func (d *Database) Get(docid string, options url.Values) (map[string]interface{}, error) {
	docRes := docResource(d.resource, docid)
	_, data, err := docRes.GetJSON("", nil, options)
	if err != nil {
		return nil, err
	}
	var doc map[string]interface{}
	doc, err = parseData(data)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Delete deletes the document with the specified ID.
func (d *Database) Delete(docid string) error {
	docRes := docResource(d.resource, docid)
	header, _, err := docRes.Head("", nil, nil)
	if err != nil {
		return err
	}
	rev := strings.Trim(header.Get("ETag"), `"`)
	return deleteDoc(docRes, rev)
}

// DeleteDoc deletes the specified document
func (d *Database) DeleteDoc(doc map[string]interface{}) error {
	id, ok := doc["_id"]
	if !ok || id == nil {
		return errors.New("document ID not existed")
	}

	rev, ok := doc["_rev"]
	if !ok || rev == nil {
		return errors.New("document rev not existed")
	}

	docRes := docResource(d.resource, id.(string))
	return deleteDoc(docRes, rev.(string))
}

func deleteDoc(docRes *Resource, rev string) error {
	_, _, err := docRes.DeleteJSON("", nil, url.Values{"rev": []string{rev}})
	return err
}

// Set creates or updates a document with the specified ID.
func (d *Database) Set(docid string, doc map[string]interface{}) error {
	docRes := docResource(d.resource, docid)
	_, data, err := docRes.PutJSON("", nil, doc, nil)
	if err != nil {
		return err
	}

	result, err := parseData(data)
	if err != nil {
		return err
	}

	doc["_id"] = result["id"].(string)
	doc["_rev"] = result["rev"].(string)
	return nil
}

// Contains returns true if the database contains a document with the specified ID.
func (d *Database) Contains(docid string) error {
	docRes := docResource(d.resource, docid)
	_, _, err := docRes.Head("", nil, nil)
	return err
}

// UpdateResult represents result of an update.
type UpdateResult struct {
	ID, Rev string
	Err     error
}

// Update performs a bulk update or creation of the given documents in a single HTTP request.
// It returns a 3-tuple (id, rev, error)
func (d *Database) Update(docs []map[string]interface{}, options map[string]interface{}) ([]UpdateResult, error) {
	results := make([]UpdateResult, len(docs))
	body := map[string]interface{}{}
	if options != nil {
		for k, v := range options {
			body[k] = v
		}
	}
	body["docs"] = docs

	_, data, err := d.resource.PostJSON("_bulk_docs", nil, body, nil)
	if err != nil {
		return nil, err
	}
	var jsonArr []map[string]interface{}
	err = json.Unmarshal(data, &jsonArr)
	if err != nil {
		return nil, err
	}

	for i, v := range jsonArr {
		var retErr error
		var result UpdateResult
		if val, ok := v["error"]; ok {
			errMsg := val.(string)
			switch errMsg {
			case "conflict":
				retErr = ErrConflict
			case "forbidden":
				retErr = ErrForbidden
			default:
				retErr = ErrInternalServerError
			}
			result = UpdateResult{
				ID:  v["id"].(string),
				Rev: "",
				Err: retErr,
			}
		} else {
			id, rev := v["id"].(string), v["rev"].(string)
			result = UpdateResult{
				ID:  id,
				Rev: rev,
				Err: retErr,
			}
			doc := docs[i]
			doc["_id"] = id
			doc["_rev"] = rev
		}
		results[i] = result
	}
	return results, nil
}

// DocIDs returns the IDs of all documents in database.
func (d *Database) DocIDs() ([]string, error) {
	docRes := docResource(d.resource, "_all_docs")
	_, data, err := docRes.GetJSON("", nil, nil)
	if err != nil {
		return nil, err
	}
	var jsonMap map[string]*json.RawMessage
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return nil, err
	}
	var jsonArr []*json.RawMessage
	json.Unmarshal(*jsonMap["rows"], &jsonArr)
	ids := make([]string, len(jsonArr))
	for i, v := range jsonArr {
		var row map[string]interface{}
		err = json.Unmarshal(*v, &row)
		if err != nil {
			return ids, err
		}
		ids[i] = row["id"].(string)
	}
	return ids, nil
}

// Name returns the name of database.
func (d *Database) Name() (string, error) {
	var name string
	info, err := d.Info("")
	if err != nil {
		return name, err
	}
	return info["db_name"].(string), nil
}

// Info returns the information about the database or design document
func (d *Database) Info(ddoc string) (map[string]interface{}, error) {
	var data []byte
	var err error
	if ddoc == "" {
		_, data, err = d.resource.GetJSON("", nil, url.Values{})
		if err != nil {
			return nil, err
		}
	} else {
		_, data, err = d.resource.GetJSON(fmt.Sprintf("_design/%s/_info", ddoc), nil, nil)
		if err != nil {
			return nil, err
		}
	}

	var info map[string]interface{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (d *Database) String() string {
	return fmt.Sprintf("Database %s", d.resource.base)
}

// Commit flushes any recent changes to the specified database to disk.
// If the server is configured to delay commits or previous requests use the special
// "X-Couch-Full-Commit: false" header to disable immediate commits, this method
// can be used to ensure that non-commited changes are commited to physical storage.
func (d *Database) Commit() error {
	_, _, err := d.resource.PostJSON("_ensure_full_commit", nil, nil, nil)
	return err
}

// Compact compacts the database by compressing the disk database file.
func (d *Database) Compact() error {
	_, _, err := d.resource.PostJSON("_compact", nil, nil, nil)
	return err
}

// Revisions returns all available revisions of the given document in reverse
// order, e.g. latest first.
func (d *Database) Revisions(docid string, options url.Values) ([]map[string]interface{}, error) {
	docRes := docResource(d.resource, docid)
	_, data, err := docRes.GetJSON("", nil, url.Values{"revs": []string{"true"}})
	if err != nil {
		return nil, err
	}
	var jsonMap map[string]*json.RawMessage
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return nil, err
	}
	var revsMap map[string]interface{}
	err = json.Unmarshal(*jsonMap["_revisions"], &revsMap)
	if err != nil {
		return nil, err
	}
	startRev := int(revsMap["start"].(float64))
	val := reflect.ValueOf(revsMap["ids"])
	if options == nil {
		options = url.Values{}
	}
	docs := make([]map[string]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		rev := fmt.Sprintf("%d-%s", startRev-i, val.Index(i).Interface().(string))
		options.Set("rev", rev)
		doc, err := d.Get(docid, options)
		if err != nil {
			return nil, err
		}
		docs[i] = doc
	}
	return docs, nil
}

// GetAttachment returns the file attachment associated with the document.
// The raw data is returned as a []byte.
func (d *Database) GetAttachment(doc map[string]interface{}, name string) ([]byte, error) {
	docid, ok := doc["_id"]
	if !ok {
		return nil, errors.New("doc _id not existed")
	}
	return d.getAttachment(docid.(string), name)
}

// GetAttachmentID returns the file attachment associated with the document ID.
// The raw data is returned as []byte.
func (d *Database) GetAttachmentID(docid, name string) ([]byte, error) {
	return d.getAttachment(docid, name)
}

func (d *Database) getAttachment(docid, name string) ([]byte, error) {
	docRes := docResource(docResource(d.resource, docid), name)
	_, data, err := docRes.Get("", nil, nil)
	return data, err
}

// PutAttachment uploads the supplied []byte as an attachment to the specified document.
// doc: the document that the attachment belongs to. Must have _id and _rev inside.
// content: the data to be attached to doc.
// name: name of attachment.
// mimeType: MIME type of content.
func (d *Database) PutAttachment(doc map[string]interface{}, content []byte, name, mimeType string) error {
	if id, ok := doc["_id"]; !ok || id.(string) == "" {
		return errors.New("doc _id not existed")
	}
	if rev, ok := doc["_rev"]; !ok || rev.(string) == "" {
		return errors.New("doc _rev not extisted")
	}

	id, rev := doc["_id"].(string), doc["_rev"].(string)

	docRes := docResource(docResource(d.resource, id), name)
	header := http.Header{}
	header.Set("Content-Type", mimeType)
	params := url.Values{}
	params.Set("rev", rev)

	_, data, err := docRes.Put("", header, content, params)
	if err != nil {
		return err
	}

	result, err := parseData(data)
	if err != nil {
		return err
	}

	doc["_rev"] = result["rev"].(string)
	return nil
}

// DeleteAttachment deletes the specified attachment
func (d *Database) DeleteAttachment(doc map[string]interface{}, name string) error {
	if id, ok := doc["_id"]; !ok || id.(string) == "" {
		return errors.New("doc _id not existed")
	}
	if rev, ok := doc["_rev"]; !ok || rev.(string) == "" {
		return errors.New("doc _rev not extisted")
	}

	id, rev := doc["_id"].(string), doc["_rev"].(string)

	params := url.Values{}
	params.Set("rev", rev)
	docRes := docResource(docResource(d.resource, id), name)
	_, data, err := docRes.DeleteJSON("", nil, params)
	if err != nil {
		return err
	}

	result, err := parseData(data)
	if err != nil {
		return err
	}
	doc["_rev"] = result["rev"]

	return nil
}

// Copy copies an existing document to a new or existing document.
func (d *Database) Copy(srcID, destID, destRev string) (string, error) {
	docRes := docResource(d.resource, srcID)
	var destination string
	if destRev != "" {
		destination = fmt.Sprintf("%s?rev=%s", destID, destRev)
	} else {
		destination = destID
	}
	header := http.Header{
		"Destination": []string{destination},
	}
	_, data, err := request("COPY", docRes.base, header, nil, nil)
	var rev string
	if err != nil {
		return rev, err
	}
	result, err := parseData(data)
	if err != nil {
		return rev, err
	}
	rev = result["rev"].(string)
	return rev, nil
}

// Changes returns a sorted list of changes feed made to documents in the database.
func (d *Database) Changes(options url.Values) (map[string]interface{}, error) {
	_, data, err := d.resource.GetJSON("_changes", nil, options)
	if err != nil {
		return nil, err
	}
	result, err := parseData(data)
	return result, err
}

// Purge performs complete removing of the given documents.
func (d *Database) Purge(docs []map[string]interface{}) (map[string]interface{}, error) {
	revs := map[string][]string{}
	for _, doc := range docs {
		id, rev := doc["_id"].(string), doc["_rev"].(string)
		if _, ok := revs[id]; !ok {
			revs[id] = []string{}
		}
		revs[id] = append(revs[id], rev)
	}

	body := map[string]interface{}{}
	for k, v := range revs {
		body[k] = v
	}
	_, data, err := d.resource.PostJSON("_purge", nil, body, nil)
	if err != nil {
		return nil, err
	}

	return parseData(data)
}

func parseData(data []byte) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	if _, ok := result["error"]; ok {
		reason := result["reason"].(string)
		return result, errors.New(reason)
	}
	return result, nil
}

func parseRaw(data []byte) (map[string]*json.RawMessage, error) {
	result := map[string]*json.RawMessage{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	if _, ok := result["error"]; ok {
		var reason string
		json.Unmarshal(*result["reason"], &reason)
		return result, errors.New(reason)
	}
	return result, nil
}

// GenerateUUID returns a random 128-bit UUID
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

// SetSecurity sets the security object for the given database.
func (d *Database) SetSecurity(securityDoc map[string]interface{}) error {
	_, _, err := d.resource.PutJSON("_security", nil, securityDoc, nil)
	return err
}

// GetSecurity returns the current security object from the given database.
func (d *Database) GetSecurity() (map[string]interface{}, error) {
	_, data, err := d.resource.GetJSON("_security", nil, nil)
	if err != nil {
		return nil, err
	}
	return parseData(data)
}

// Len returns the number of documents stored in it.
func (d *Database) Len() (int, error) {
	info, err := d.Info("")
	if err != nil {
		return 0, err
	}
	return int(info["doc_count"].(float64)), nil
}

// GetRevsLimit gets the current revs_limit(revision limit) setting.
func (d *Database) GetRevsLimit() (int, error) {
	_, data, err := d.resource.Get("_revs_limit", nil, nil)
	if err != nil {
		return 0, err
	}
	limit, err := strconv.Atoi(strings.Trim(string(data), "\n"))
	if err != nil {
		return limit, err
	}
	return limit, nil
}

// SetRevsLimit sets the maximum number of document revisions that will be
// tracked by CouchDB.
func (d *Database) SetRevsLimit(limit int) error {
	_, _, err := d.resource.Put("_revs_limit", nil, []byte(strconv.Itoa(limit)), nil)
	return err
}

// docResource returns a Resource instance for docID
func docResource(res *Resource, docID string) *Resource {
	if len(docID) == 0 {
		return res
	}

	docRes := res
	if docID[:1] == "_" {
		paths := strings.SplitN(docID, "/", 2)
		for _, p := range paths {
			docRes, _ = docRes.NewResourceWithURL(p)
		}
		return docRes
	}

	docRes, _ = res.NewResourceWithURL(url.QueryEscape(docID))
	return docRes
}

// Cleanup removes all view index files no longer required by CouchDB.
func (d *Database) Cleanup() error {
	_, _, err := d.resource.PostJSON("_view_cleanup", nil, nil, nil)
	return err
}

// Query returns documents using a conditional selector statement in Golang.
//
// selector: A filter string declaring which documents to return, formatted as a Golang statement.
//
// fields: Specifying which fields to be returned, if passing nil the entire
// is returned, no automatic inclusion of _id or other metadata fields.
//
// sorts: How to order the documents returned, formatted as ["desc(fieldName1)", "desc(fieldName2)"]
// or ["fieldNameA", "fieldNameB"] of which "asc" is used by default, passing nil to disable ordering.
//
// limit: Maximum number of results returned, passing nil to use default value(25).
//
// skip: Skip the first 'n' results, where 'n' is the number specified, passing nil for no-skip.
//
// index: Instruct a query to use a specific index, specified either as "<design_document>" or
// ["<design_document>", "<index_name>"], passing nil to use primary index(_all_docs) by default.
//
// Inner functions for selector syntax
//
// nor(condexprs...) matches if none of the conditions in condexprs match($nor).
//
// For example: nor(year == 1990, year == 1989, year == 1997) returns all documents
// whose year field not in 1989, 1990 and 1997.
//
// all(field, array) matches an array value if it contains all the elements of the argument array($all).
//
// For example: all(genre, []string{"Comedy", "Short"} returns all documents whose
// genre field contains "Comedy" and "Short".
//
// any(field, condexpr) matches an array field with at least one element meets the specified condition($elemMatch).
//
// For example: any(genre, genre == "Short" || genre == "Horror") returns all documents whose
// genre field contains "Short" or "Horror" or both.
//
// exists(field, boolean) checks whether the field exists or not, regardless of its value($exists).
//
// For example: exists(director, false) returns all documents who does not have a director field.
//
// typeof(field, type) checks the document field's type, valid types are
// "null", "boolean", "number", "string", "array", "object"($type).
//
// For example: typeof(genre, "array") returns all documents whose genre field is of array type.
//
// in(field, array) the field must exist in the array provided($in).
//
// For example: in(director, []string{"Mike Portnoy", "Vitali Kanevsky"}) returns all documents
// whose director field is "Mike Portnoy" or "Vitali Kanevsky".
//
// nin(field, array) the document field must not exist in the array provided($nin).
//
// For example: nin(year, []int{1990, 1992, 1998}) returns all documents whose year field is not
// in 1990, 1992 or 1998.
//
// size(field, int) matches the length of an array field in a document($size).
//
// For example: size(genre, 2) returns all documents whose genre field is of length 2.
//
// mod(field, divisor, remainder) matches documents where field % divisor == remainder($mod).
//
// For example: mod(year, 2, 1) returns all documents whose year field is an odd number.
//
// regex(field, regexstr) a regular expression pattern to match against the document field.
//
// For example: regex(title, "^A") returns all documents whose title is begin with an "A".
//
// Inner functions for sort syntax
//
// asc(field) sorts the field in ascending order, this is the default option while
// desc(field) sorts the field in descending order.
func (d *Database) Query(fields []string, selector string, sorts []string, limit, skip, index interface{}) ([]map[string]interface{}, error) {
	selectorJSON, err := parseSelectorSyntax(selector)
	if err != nil {
		return nil, err
	}
	find := map[string]interface{}{
		"selector": selectorJSON,
	}

	if limitVal, ok := limit.(int); ok {
		find["limit"] = limitVal
	}

	if skipVal, ok := skip.(int); ok {
		find["skip"] = skipVal
	}

	if sorts != nil {
		sortsJSON, err := parseSortSyntax(sorts)
		if err != nil {
			return nil, err
		}
		find["sort"] = sortsJSON
	}

	if fields != nil {
		find["fields"] = fields
	}

	if index != nil {
		find["use_index"] = index
	}

	return d.queryJSON(find)
}

// QueryJSON returns documents using a declarative JSON querying syntax.
func (d *Database) QueryJSON(query string) ([]map[string]interface{}, error) {
	queryMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(query), &queryMap)
	if err != nil {
		return nil, err
	}
	return d.queryJSON(queryMap)
}

func (d *Database) queryJSON(queryMap map[string]interface{}) ([]map[string]interface{}, error) {
	_, data, err := d.resource.PostJSON("_find", nil, queryMap, nil)
	if err != nil {
		return nil, err
	}

	result, err := parseRaw(data)
	if err != nil {
		return nil, err
	}

	docs := []map[string]interface{}{}
	err = json.Unmarshal(*result["docs"], &docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

// parseSelectorSyntax returns a map representing the selector JSON struct.
func parseSelectorSyntax(selector string) (interface{}, error) {
	// protect selector against query selector injection attacks
	if strings.Contains(selector, "$") {
		return nil, fmt.Errorf("no $s are allowed in selector: %s", selector)
	}

	// parse selector into abstract syntax tree (ast)
	expr, err := parser.ParseExpr(selector)
	if err != nil {
		return nil, err
	}

	// recursively processing ast into json object
	selectObj, err := parseAST(expr)
	if err != nil {
		return nil, err
	}

	return selectObj, nil
}

// parseSortSyntax returns a slice of sort JSON struct.
func parseSortSyntax(sorts []string) (interface{}, error) {
	if sorts == nil {
		return nil, nil
	}

	sortObjs := []interface{}{}
	for _, sort := range sorts {
		sortExpr, err := parser.ParseExpr(sort)
		if err != nil {
			return nil, err
		}

		sortObj, err := parseAST(sortExpr)
		if err != nil {
			return nil, err
		}
		sortObjs = append(sortObjs, sortObj)
	}

	return sortObjs, nil
}

// parseAST converts and returns a JSON struct according to
// CouchDB mango query syntax for the abstract syntax tree represented by expr.
func parseAST(expr ast.Expr) (interface{}, error) {
	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		// fmt.Println("BinaryExpr", expr)
		return parseBinary(expr.Op, expr.X, expr.Y)
	case *ast.UnaryExpr:
		// fmt.Println("UnaryExpr", expr)
		return parseUnary(expr.Op, expr.X)
	case *ast.CallExpr:
		// fmt.Println("CallExpr", expr, expr.Fun, expr.Args)
		return parseFuncCall(expr.Fun, expr.Args)
	case *ast.Ident:
		// fmt.Println("Ident", expr)
		switch expr.Name {
		case "nil": // for nil value such as _id > nil
			return nil, nil
		case "true": // for boolean value true
			return true, nil
		case "false":
			return false, nil // for boolean value false
		default:
			return expr.Name, nil
		}
	case *ast.BasicLit:
		// fmt.Println("BasicLit", expr)
		switch expr.Kind {
		case token.INT:
			intVal, err := strconv.Atoi(expr.Value)
			if err != nil {
				return nil, err
			}
			return intVal, nil
		case token.FLOAT:
			floatVal, err := strconv.ParseFloat(expr.Value, 64)
			if err != nil {
				return nil, err
			}
			return floatVal, nil
		case token.STRING:
			return strings.Trim(expr.Value, `"`), nil
		default:
			return nil, fmt.Errorf("token type %s not supported", expr.Kind.String())
		}
	case *ast.SelectorExpr:
		// fmt.Println("SelectorExpr", expr.X, expr.Sel)
		xExpr, err := parseAST(expr.X)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%s.%s", xExpr, expr.Sel.Name), nil
	case *ast.ParenExpr:
		pExpr, err := parseAST(expr.X)
		if err != nil {
			return nil, err
		}
		return pExpr, nil
	case *ast.CompositeLit:
		if _, ok := expr.Type.(*ast.ArrayType); !ok {
			return nil, fmt.Errorf("not an ArrayType for a composite literal %v", expr.Type)
		}
		elements := make([]interface{}, len(expr.Elts))
		for idx, elt := range expr.Elts {
			e, err := parseAST(elt)
			if err != nil {
				return nil, err
			}
			elements[idx] = e
		}
		return elements, nil
	default:
		return nil, fmt.Errorf("expressions other than unary, binary and function call are not allowed %v", expr)
	}
}

// parseBinary parses and returns a JSON struct according to
// CouchDB mango query syntax for the supported binary operators.
func parseBinary(operator token.Token, leftOperand, rightOperand ast.Expr) (interface{}, error) {
	left, err := parseAST(leftOperand)
	if err != nil {
		return nil, err
	}
	right, err := parseAST(rightOperand)
	if err != nil {
		return nil, err
	}

	// <, <=, ==, !=, >=, >, &&, ||
	switch operator {
	case token.LSS:
		return map[string]interface{}{
			left.(string): map[string]interface{}{"$lt": right},
		}, nil
	case token.LEQ:
		return map[string]interface{}{
			left.(string): map[string]interface{}{"$lte": right},
		}, nil
	case token.EQL:
		return map[string]interface{}{
			left.(string): map[string]interface{}{"$eq": right},
		}, nil
	case token.NEQ:
		return map[string]interface{}{
			left.(string): map[string]interface{}{"$ne": right},
		}, nil
	case token.GEQ:
		return map[string]interface{}{
			left.(string): map[string]interface{}{"$gte": right},
		}, nil
	case token.GTR:
		return map[string]interface{}{
			left.(string): map[string]interface{}{"$gt": right},
		}, nil
	case token.LAND:
		return map[string]interface{}{
			"$and": []interface{}{left, right},
		}, nil
	case token.LOR:
		return map[string]interface{}{
			"$or": []interface{}{left, right},
		}, nil
	}
	return nil, fmt.Errorf("binary operator %s not supported", operator)
}

// parseUnary parses and returns a JSON struct according to
// CouchDB mango query syntax for supported unary operators.
func parseUnary(operator token.Token, operandExpr ast.Expr) (interface{}, error) {
	operand, err := parseAST(operandExpr)
	if err != nil {
		return nil, err
	}

	switch operator {
	case token.NOT:
		return map[string]interface{}{
			"$not": operand,
		}, nil
	}
	return nil, fmt.Errorf("unary operator %s not supported", operator)
}

// parseFuncCall parses and returns a JSON struct according to
// CouchDB mango query syntax for supported meta functions.
func parseFuncCall(funcExpr ast.Expr, args []ast.Expr) (interface{}, error) {
	funcIdent := funcExpr.(*ast.Ident)
	functionName := funcIdent.Name
	switch functionName {
	case "nor":
		if len(args) < 1 {
			return nil, fmt.Errorf("function nor(exprs...) need at least 1 arguments, not %d", len(args))
		}

		selectors := make([]interface{}, len(args))
		for idx, arg := range args {
			selector, err := parseAST(arg)
			if err != nil {
				return nil, err
			}
			selectors[idx] = selector
		}

		return map[string]interface{}{
			"$nor": selectors,
		}, nil
	case "all":
		if len(args) != 2 {
			return nil, fmt.Errorf("function all(field, array) need 2 arguments, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		arrayExpr, err := parseAST(args[1])
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			fieldExpr.(string): map[string]interface{}{"$all": arrayExpr},
		}, nil
	case "any":
		if len(args) != 2 {
			return nil, fmt.Errorf("function any(field, condition) need 2 arguments, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		anyExpr, err := parseAST(args[1])
		if err != nil {
			return nil, err
		}
		anyExpr, err = removeFieldKey(fieldExpr.(string), anyExpr)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			fieldExpr.(string): map[string]interface{}{"$elemMatch": anyExpr},
		}, nil
	case "exists":
		if len(args) != 2 {
			return nil, fmt.Errorf("function exists(field, boolean) need 2 arguments, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		boolExpr, err := parseAST(args[1])
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			fieldExpr.(string): map[string]interface{}{"$exists": boolExpr},
		}, nil
	case "typeof":
		if len(args) != 2 {
			return nil, fmt.Errorf("function typeof(field, type) need 2 arguments, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		typeStr, err := parseAST(args[1])
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			fieldExpr.(string): map[string]interface{}{"$type": typeStr},
		}, nil
	case "in":
		if len(args) != 2 {
			return nil, fmt.Errorf("function in(field, array) need 2 arguments, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		arrExpr, err := parseAST(args[1])
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			fieldExpr.(string): map[string]interface{}{"$in": arrExpr},
		}, nil
	case "nin":
		if len(args) != 2 {
			return nil, fmt.Errorf("function nin(field, array) need 2 arguments, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		arrExpr, err := parseAST(args[1])
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			fieldExpr.(string): map[string]interface{}{"$nin": arrExpr},
		}, nil
	case "size":
		if len(args) != 2 {
			return nil, fmt.Errorf("function size(field, int) need 2 arguments, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		intExpr, err := parseAST(args[1])
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			fieldExpr.(string): map[string]interface{}{"$size": intExpr},
		}, nil
	case "mod":
		if len(args) != 3 {
			return nil, fmt.Errorf("function mod(field, divisor, remainder) need 3 arguments, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		divisorExpr, err := parseAST(args[1])
		if err != nil {
			return nil, err
		}
		divisor, ok := divisorExpr.(int)
		if !ok {
			return nil, fmt.Errorf("invalid divisor %s", divisorExpr)
		}

		remainderExpr, err := parseAST(args[2])
		if err != nil {
			return nil, err
		}
		remainder, ok := remainderExpr.(int)
		if !ok {
			return nil, fmt.Errorf("invalid remainder %s", remainderExpr)
		}

		expr, err := parser.ParseExpr(fmt.Sprintf("%#v", []int{divisor, remainder}))
		if err != nil {
			return nil, err
		}
		arrExpr, err := parseAST(expr)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			fieldExpr.(string): map[string]interface{}{"$mod": arrExpr},
		}, nil
	case "regex":
		if len(args) != 2 {
			return nil, fmt.Errorf("function regex(field, regexstr) need 2 arguments, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		regexExpr, err := parseAST(args[1])
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			fieldExpr.(string): map[string]interface{}{"$regex": regexExpr},
		}, nil
	case "asc": // for sort syntax
		if len(args) != 1 {
			return nil, fmt.Errorf("function asc(field) need 1 argument, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		return map[string]interface{}{
			fieldExpr.(string): "asc",
		}, nil
	case "desc": // for sort syntax
		if len(args) != 1 {
			return nil, fmt.Errorf("function desc(field) need 1 argument, not %d", len(args))
		}

		fieldExpr, err := parseAST(args[0])
		if err != nil {
			return nil, err
		}
		if _, ok := fieldExpr.(string); !ok {
			return nil, fmt.Errorf("invalid field expression type %s", fieldExpr)
		}

		return map[string]interface{}{
			fieldExpr.(string): "desc",
		}, nil
	}
	return nil, fmt.Errorf("function %s() not supported", functionName)
}

// removeFieldKey removes the key which equals to fieldName,
// moves its value one level up in the map.
func removeFieldKey(fieldName string, exprMap interface{}) (interface{}, error) {
	mapValue := reflect.ValueOf(exprMap)
	if mapValue.Kind() != reflect.Map {
		return nil, errors.New("not a map type")
	}
	mapKeys := mapValue.MapKeys()
	for _, mapKey := range mapKeys {
		// exprMap is a interface type contains map[string]interface{}
		// so MapIndex returns a value whose Kind is Interface, so we
		// have to call its Interface() methods then pass to ValueOf()
		// to get the underline map type.
		value := reflect.ValueOf(mapValue.MapIndex(mapKey).Interface())
		if value.Kind() == reflect.Slice {
			for idx := 0; idx < value.Len(); idx++ {
				elemVal := value.Index(idx)
				processed, err := removeFieldKey(fieldName, elemVal.Interface())
				if err != nil {
					return nil, err
				}
				elemVal.Set(reflect.ValueOf(processed))
			}
			mapValue.SetMapIndex(mapKey, value)
		} else if value.Kind() == reflect.Map {
			if mapKey.Interface().(string) == fieldName { // found
				if value.Len() != 1 {
					return nil, fmt.Errorf("field map length %d, not 1", value.Len())
				}
				// setting to empty value deletes the key
				mapValue.SetMapIndex(mapKey, reflect.Value{})
				keys := value.MapKeys()
				// moves the value one level up
				for _, key := range keys {
					val := value.MapIndex(key)
					mapValue.SetMapIndex(key, val)
				}
			} else {
				processed, err := removeFieldKey(fieldName, value.Interface())
				if err != nil {
					return nil, err
				}
				mapValue.SetMapIndex(mapKey, reflect.ValueOf(processed))
			}
		}
	}
	return mapValue.Interface(), nil
}

// beautifulJSONString returns a beautified string representing the JSON struct.
func beautifulJSONString(jsonable interface{}) (string, error) {
	b, err := json.Marshal(jsonable)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// PutIndex creates a new index in database.
//
// indexFields: a JSON array of field names following the sort syntax.
//
// ddoc: optional, name of the design document in which the index will be created.
// By default each index will be created in its own design document. Indexes can be
// grouped into design documents for efficiency. However a change to one index
// in a design document will invalidate all other indexes in the same document.
//
// name: optional, name of the index. A name generated automatically if not provided.
func (d *Database) PutIndex(indexFields []string, ddoc, name string) (string, string, error) {
	var design, index string
	if len(indexFields) == 0 {
		return design, index, errors.New("index fields cannot be empty")
	}

	indexObjs, err := parseSortSyntax(indexFields)
	if err != nil {
		return design, index, err
	}

	indexJSON := map[string]interface{}{}
	indexJSON["index"] = map[string]interface{}{
		"fields": indexObjs,
	}

	if len(ddoc) > 0 {
		indexJSON["ddoc"] = ddoc
	}

	if len(name) > 0 {
		indexJSON["name"] = name
	}

	_, data, err := d.resource.PostJSON("_index", nil, indexJSON, nil)
	if err != nil {
		return design, index, err
	}

	result, err := parseData(data)
	if err != nil {
		return design, index, err
	}
	design = result["id"].(string)
	index = result["name"].(string)

	return design, index, nil
}

// GetIndex gets all indexes created in database.
func (d *Database) GetIndex() (map[string]*json.RawMessage, error) {
	_, data, err := d.resource.GetJSON("_index", nil, nil)
	if err != nil {
		return nil, err
	}
	return parseRaw(data)
}

// DeleteIndex deletes index in database.
func (d *Database) DeleteIndex(ddoc, name string) error {
	indexRes := docResource(d.resource, fmt.Sprintf("_index/%s/json/%s", ddoc, name))
	_, _, err := indexRes.DeleteJSON("", nil, nil)
	return err
}

// designPath resturns a make-up design path based on designDoc and designType
// for example designPath("design/foo", "_view") returns "_design/design/_view/foo"
func designPath(designDoc, designType string) string {
	if strings.HasPrefix(designDoc, "_") {
		return designDoc
	}
	parts := strings.SplitN(designDoc, "/", 2)
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join([]string{"_design", parts[0], designType, parts[1]}, "/")
}

// View executes a predefined design document view and returns the results.
//
// name: the name of the view, for user-defined views use the format "design_docid/viewname",
// that is, the document ID of the design document and the name of the view, separated by a /.
//
// wrapper: an optional function for processing the result rows after retrieved.
//
// options: optional query parameters.
func (d *Database) View(name string, wrapper func(Row) Row, options map[string]interface{}) (*ViewResults, error) {
	designDocPath := designPath(name, "_view")
	return newViewResults(d.resource, designDocPath, options, wrapper), nil
}

// IterView returns a channel fetching rows in batches which iterates a row at a time(pagination).
//
// name: the name of the view, for user-defined views use the format "design_docid/viewname",
// that is, the document ID of the design document and the name of the view, separated by a /.
//
// wrapper: an optional function for processing the result rows after retrieved.
//
// options: optional query parameters.
func (d *Database) IterView(name string, batch int, wrapper func(Row) Row, options map[string]interface{}) (<-chan Row, error) {
	if batch <= 0 {
		return nil, ErrBatchValue
	}

	if options == nil {
		options = map[string]interface{}{}
	}

	_, ok := options["limit"]
	var limit int
	if ok {
		if options["limit"].(int) <= 0 {
			return nil, ErrLimitValue
		}
		limit = options["limit"].(int)
	}

	// Row generator
	rchan := make(chan Row)
	var err error
	go func() {
		defer close(rchan)
		for {
			loopLimit := batch
			if ok {
				loopLimit = min(batch, limit)
			}
			// get rows in batch with one extra for start of next batch
			options["limit"] = loopLimit + 1
			var results *ViewResults
			results, err = d.View(name, wrapper, options)
			if err != nil {
				break
			}
			var rows []Row
			rows, err = results.Rows()
			if err != nil {
				break
			}

			// send all rows to channel except the last extra one
			for _, row := range rows[:min(len(rows), loopLimit)] {
				rchan <- row
			}

			if ok {
				limit -= min(len(rows), batch)
			}

			if len(rows) <= batch || (ok && limit == 0) {
				break
			}
			options["startkey"] = rows[len(rows)-1].Key
			options["startkey_docid"] = rows[len(rows)-1].ID
			options["skip"] = 0
		}
	}()
	return rchan, nil
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

// Show calls a server-side 'show' function.
//
// name: the name of the show function in the format "designdoc/showname"
//
// docID: optional document ID to pass to the show function
//
// params: optional query parameters
func (d *Database) Show(name, docID string, params url.Values) (http.Header, []byte, error) {
	designDocPath := designPath(name, "_show")
	if docID != "" {
		designDocPath = fmt.Sprintf("%s/%s", designDocPath, docID)
	}
	return d.resource.Get(designDocPath, nil, params)
}

// List formats a view using a server-side 'list' function.
//
// name: the name of the list function in the format "designdoc/listname"
//
// view: the name of the view in the format "designdoc/viewname"
//
// options: optional query parameters
func (d *Database) List(name, view string, options map[string]interface{}) (http.Header, []byte, error) {
	designDocPath := designPath(name, "_list")
	res := docResource(d.resource, fmt.Sprintf("%s/%s", designDocPath, strings.Split(view, "/")[1]))
	return viewLikeResourceRequest(res, options)
}

// UpdateDoc calls server-side update handler.
//
// name: the name of the update handler function in the format "designdoc/updatename".
//
// docID: optional document ID to pass to the show function
//
// params: optional query parameters
func (d *Database) UpdateDoc(name, docID string, params url.Values) (http.Header, []byte, error) {
	designDocPath := designPath(name, "_update")
	if docID == "" {
		return d.resource.Post(designDocPath, nil, nil, params)
	}

	designDocPath = fmt.Sprintf("%s/%s", designDocPath, docID)
	return d.resource.Put(designDocPath, nil, nil, params)
}
