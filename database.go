package couchdb

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

	var httpFunc func(string, http.Header, map[string]interface{}, url.Values) (http.Header, *json.RawMessage, error)
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
	err = json.Unmarshal(*data, &jsonMap)
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
func (d *Database) Get(docid string) (map[string]interface{}, error) {
	docRes := docResource(d.resource, docid)
	_, data, err := docRes.GetJSON("", nil, nil)
	if err != nil {
		return nil, err
	}
	var doc map[string]interface{}
	err = json.Unmarshal(*data, &doc)
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
	params := url.Values{}
	params.Set("rev", rev)
	_, _, err = docRes.DeleteJSON("", nil, params)
	if err != nil {
		return err
	}
	return nil
}

// Set creates or updates a document with the specified ID.
func (d *Database) Set(docid string, doc map[string]interface{}) bool {
	if doc == nil {
		return false
	}

	docRes := docResource(d.resource, docid)
	_, data, err := docRes.PutJSON("", nil, doc, nil)
	if err != nil {
		return false
	}

	var jsonMap map[string]interface{}
	json.Unmarshal(*data, &jsonMap)
	doc["_id"] = jsonMap["id"].(string)
	doc["_rev"] = jsonMap["rev"].(string)
	return true
}

// Contains returns true if the database contains a document with the specified ID.
func (d *Database) Contains(docid string) bool {
	docRes := docResource(d.resource, docid)
	_, _, err := docRes.Head("", nil, nil)
	return err == nil
}

// UpdateDocuments performs a bulk update or creation of the given documents in a single HTTP request.
func (d *Database) Update(docs []map[string]interface{}, options map[string]interface{}) ([]IDRev, error) {
	results := []IDRev{}

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
	err = json.Unmarshal(*data, &jsonArr)
	if err != nil {
		return nil, err
	}

	for _, ele := range jsonArr {
		id, rev := ele["id"].(string), ele["rev"].(string)
		results = append(results, IDRev{Id: id, Rev: rev})
	}
	return results, nil
}

// DocIDs returns the IDs of all documents in database.
func (d *Database) DocIDs() []string {
	docRes := docResource(d.resource, "_all_docs")
	_, data, err := docRes.GetJSON("", nil, nil)
	if err != nil {
		return nil
	}
	var jsonMap map[string]*json.RawMessage
	json.Unmarshal(*data, &jsonMap)
	if _, ok := jsonMap["rows"]; !ok {
		return nil
	}
	var jsonArr []*json.RawMessage
	json.Unmarshal(*jsonMap["rows"], &jsonArr)
	if len(jsonArr) == 0 {
		return nil
	}
	ids := make([]string, len(jsonArr))
	for i, v := range jsonArr {
		var row map[string]interface{}
		json.Unmarshal(*v, &row)
		ids[i] = row["id"].(string)
	}
	return ids
}

// Name returns the name of database.
func (d *Database) Name() (string, error) {
	var name string
	info, err := d.databaseInfo()
	if err != nil {
		return name, err
	}
	return info["db_name"].(string), nil
}

func (d *Database) databaseInfo() (map[string]interface{}, error) {
	_, data, err := d.resource.GetJSON("", nil, url.Values{})

	if err != nil {
		return nil, err
	}

	var info map[string]interface{}
	err = json.Unmarshal(*data, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (d *Database) String() string {
	return fmt.Sprintf("Database %s", d.resource.base)
}

///////////////////////////////////////////////////////

// Len returns the number of documents stored in it.
func (d *Database) Len() (int, error) {
	info, err := d.databaseInfo()
	if err != nil {
		return 0, err
	}
	return int(info["doc_count"].(float64)), nil
}

// docResource returns a Resource instance for docID
func docResource(res *Resource, docID string) *Resource {
	var docRes *Resource
	if docID[:1] == "_" {
		paths := strings.SplitN(docID, "/", 2)
		for _, p := range paths {
			docRes, _ = res.NewResourceWithURL(p)
		}
		return docRes
	}

	docRes, _ = res.NewResourceWithURL(docID)
	return docRes
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

// Commit flushes any recent changes to the specified database to disk.
// If the server is configured to delay commits or previous requests use the special
// "X-Couch-Full-Commit: false" header to disable immediate commits, this method
// can be used to ensure that non-commited changes are commited to physical storage.
func (d *Database) Commit() bool {
	_, _, err := d.resource.PostJSON("_ensure_full_commit", nil, nil, nil)
	return err == nil
}

// GetAttachment returns the file attachment associated with the document.
// The raw data of the associated attachment is returned as a []byte.
func (d *Database) GetAttachment(docid, fileName string) ([]byte, bool) {
	// defensive check
	if len(docid) == 0 || len(fileName) == 0 {
		return nil, false
	}

	docRes := docResource(docResource(d.resource, docid), fileName)
	_, data, err := docRes.Get("", nil, nil)
	return data, err == nil
}

// PutAttachment uploads the supplied *os.File as an attachment to the specified document.
// doc: the document that the attachment belongs to. Must have _id and _rev inside.
func (d *Database) PutAttachment(doc map[string]interface{}, file *os.File, mimeType string) bool {
	// defensive check
	if doc == nil || file == nil || len(mimeType) == 0 {
		return false
	}

	if _, ok := doc["_id"]; !ok {
		return false
	}
	if _, ok := doc["_rev"]; !ok {
		return false
	}

	id, rev := doc["_id"].(string), doc["_rev"].(string)

	if len(id) == 0 || len(rev) == 0 {
		return false
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}

	docRes := docResource(docResource(d.resource, id), fileInfo.Name())
	header := http.Header{}
	header.Set("Content-Type", mimeType)
	params := url.Values{}
	params.Set("rev", rev)

	_, data, err := docRes.Put("", header, contents, params)
	if err == nil {
		var jsonMap map[string]interface{}
		json.Unmarshal(data, &jsonMap)
		doc["_rev"] = jsonMap["rev"].(string)
	}

	return err == nil
}

// DeleteAttachment deletes the specified attachment
func (d *Database) DeleteAttachment(doc map[string]interface{}, fileName string) bool {
	// defensive check
	if doc == nil || len(fileName) == 0 {
		return false
	}

	if _, ok := doc["_id"]; !ok {
		return false
	}

	if _, ok := doc["_rev"]; !ok {
		return false
	}

	id, rev := doc["_id"].(string), doc["_rev"].(string)

	if len(id) == 0 || len(rev) == 0 {
		return false
	}

	params := url.Values{}
	params.Set("rev", rev)
	docRes := docResource(docResource(d.resource, id), fileName)
	_, data, err := docRes.DeleteJSON("", nil, params)
	if err == nil {
		var jsonMap map[string]interface{}
		json.Unmarshal(*data, &jsonMap)
		doc["_rev"] = jsonMap["rev"]
	}
	return err == nil
}

type IDRev struct {
	Id  string
	Rev string
}

// GetRevsLimit gets the current revs_limit(revision limit) setting.
func (d *Database) GetRevsLimit() (int, error) {
	limit := -1
	_, data, err := d.resource.Get("_revs_limit", nil, nil)
	if err != nil {
		return limit, err
	}
	limit, err = strconv.Atoi(strings.Trim(string(data), "\n"))
	if err != nil {
		return limit, err
	}
	return limit, nil
}

// SetRevsLimit sets the maximum number of document revisions that will be
// tracked by CouchDB.
func (d *Database) SetRevsLimit(limit int) bool {
	_, _, err := d.resource.Put("_revs_limit", nil, []byte(strconv.Itoa(limit)), nil)
	return err == nil
}

// Changes returns a sorted list of changes feed made to documents in the database.
func (d *Database) Changes(options url.Values) (map[string]interface{}, bool) {
	_, data, err := d.resource.GetJSON("_changes", nil, options)
	if err != nil {
		return nil, false
	}
	var changes map[string]interface{}
	json.Unmarshal(*data, &changes)
	return changes, err == nil
}

// Cleanup removes all view index files no longer required by CouchDB.
func (d *Database) Cleanup() bool {
	_, _, err := d.resource.PostJSON("_view_cleanup", nil, nil, nil)
	return err == nil
}

// Compact compacts the database by compressing the disk database file.
func (d *Database) Compact() bool {
	_, _, err := d.resource.PostJSON("_compact", nil, nil, nil)
	return err == nil
}

// Copy copies an existing document to a new or existing document.
func (d *Database) Copy(srcID, destID string) (string, bool) {
	docRes := docResource(d.resource, srcID)
	header := http.Header{
		"Destination": []string{destID},
	}
	_, data, err := request("COPY", docRes.base, header, nil, nil)
	var rev string
	if err == nil {
		var jsonMap map[string]interface{}
		json.Unmarshal(data, &jsonMap)
		rev = jsonMap["rev"].(string)
	}

	return rev, err == nil
}

// Purge performs complete removing of the given documents.
func (d *Database) Purge(docIDs []string) bool {
	// TODO
	return false
}

func (d *Database) SetSecurity(securityDoc map[string]interface{}) bool {
	_, _, err := d.resource.PutJSON("_security", nil, securityDoc, nil)
	return err == nil
}

func (d *Database) GetSecurity() (map[string]interface{}, bool) {
	_, data, err := d.resource.GetJSON("_security", nil, nil)
	var secDoc map[string]interface{}
	if err == nil {
		json.Unmarshal(*data, &secDoc)
	}
	return secDoc, err == nil
}

// GetRevisions returns all available revisions of the given document in reverse
// order, e.g. latest first.TODO
func (d *Database) GetRevisions() {}
