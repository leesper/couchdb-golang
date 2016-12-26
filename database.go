package couchdb

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	DEFAULT_BASE_URL = "http://localhost:5984"
)

// getDefaultCouchDBURL returns the default CouchDB server url.
func getDefaultCouchDBURL() string {
	var couchdbUrlEnviron string
	for _, couchdbUrlEnviron = range os.Environ() {
		if strings.HasPrefix(couchdbUrlEnviron, "COUCHDB_URL") {
			break
		}
	}
	if len(couchdbUrlEnviron) == 0 {
		couchdbUrlEnviron = DEFAULT_BASE_URL
	} else {
		couchdbUrlEnviron = strings.Split(couchdbUrlEnviron, "=")[1]
	}
	return couchdbUrlEnviron
}

// Database represents a CouchDB database instance.
type Database struct {
	resource *Resource
}

// NewDatabase returns a CouchDB database instance.
func NewDatabase(urlStr string) (*Database, error) {
	var dbUrlStr string
	if !strings.HasPrefix(urlStr, "http") {
		base, err := url.Parse(getDefaultCouchDBURL())
		if err != nil {
			return nil, err
		}
		dbUrl, err := base.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		dbUrlStr = dbUrl.String()
	} else {
		dbUrlStr = urlStr
	}

	res, err := NewResource(dbUrlStr, nil)
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

// Aavailable returns true if the database is good to go.
func (d *Database) Available() bool {
	_, _, err := d.resource.Head("", nil, nil)
	return err == nil
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

// Update performs a bulk update or creation of the given documents in a single HTTP request.
// It returns a 3-tuple (id, rev, error)
type UpdateResult struct {
	id, rev string
	err     error
}

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
				id:  v["id"].(string),
				rev: "",
				err: retErr,
			}
		} else {
			id, rev := v["id"].(string), v["rev"].(string)
			result = UpdateResult{
				id:  id,
				rev: rev,
				err: retErr,
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
	info, err := d.Info()
	if err != nil {
		return name, err
	}
	return info["db_name"].(string), nil
}

// Info returns the information about the database
func (d *Database) Info() (map[string]interface{}, error) {
	_, data, err := d.resource.GetJSON("", nil, url.Values{})

	if err != nil {
		return nil, err
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

func (d *Database) SetSecurity(securityDoc map[string]interface{}) error {
	_, _, err := d.resource.PutJSON("_security", nil, securityDoc, nil)
	return err
}

func (d *Database) GetSecurity() (map[string]interface{}, error) {
	_, data, err := d.resource.GetJSON("_security", nil, nil)
	if err != nil {
		return nil, err
	}
	return parseData(data)
}

// Len returns the number of documents stored in it.
func (d *Database) Len() (int, error) {
	info, err := d.Info()
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

	var docRes *Resource
	if docID[:1] == "_" {
		paths := strings.SplitN(docID, "/", 2)
		for _, p := range paths {
			docRes, _ = res.NewResourceWithURL(p)
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

func (d *Database) Query( /* fields, selector, skip, sort, limit, use_index */ ) {}
