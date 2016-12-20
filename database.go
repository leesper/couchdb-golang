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
func NewDatabase(urlStr string) *Database {
  var dbUrlStr string
  if !strings.HasPrefix(urlStr, "http") {
    base, err := url.Parse(getDefaultCouchDBURL())
    if err != nil {
      return nil
    }
    dbUrl, err := base.Parse(urlStr)
    if err != nil {
      return nil
    }
    dbUrlStr = dbUrl.String()
  } else {
    dbUrlStr = urlStr
  }
  res := NewResource(dbUrlStr, nil)

  if res == nil {
    return nil
  }

  return &Database{
    resource: res,
  }
}

// NewDatabaseWithResource returns a CouchDB database instance with resource obj.
func NewDatabaseWithResource(res *Resource) *Database {
  return &Database{
    resource: res,
  }
}

// Name returns the name of database.
func (d *Database)Name() string {
  info := d.databaseInfo()
  if _, ok := info["db_name"]; !ok {
    return ""
  }

  return info["db_name"].(string)
}

func (d *Database)databaseInfo() map[string]interface{} {
  _, _, jsonData := d.resource.GetJSON("", nil, url.Values{})

  var jsonMap map[string]interface{}

  if jsonData == nil {
    return jsonMap
  }

  json.Unmarshal(*jsonData, &jsonMap)

  return jsonMap
}

// Aavailable returns true if the database is good to go.
func (d *Database)Available() bool {
  status, _, _ := d.resource.Head("", nil, nil)
  return status == OK
}

// Contains returns true if the database contains a document with the specified ID.
func (d *Database)Contains(docid string) bool {
  docRes := docResource(d.resource, docid)
  status, _, _ := docRes.Head("", nil, nil)
  return status == OK
}

// Get returns the document with the specified ID.
func (d *Database)Get(docid string) map[string]interface{} {
  docRes := docResource(d.resource, docid)
  status, _, data := docRes.GetJSON("", nil, nil)
  if status != OK {
    return nil
  }
  var doc map[string]interface{}
  json.Unmarshal(*data, &doc)
  return doc
}

// Delete deletes the document with the specified ID.
func (d *Database)Delete(docid string) bool {
  docRes := docResource(d.resource, docid)
  status, header, _ := docRes.Head("", nil, nil)
  if status != OK {
    return false
  }
  rev := strings.Trim(header.Get("ETag"), `"`)
  params := url.Values{}
  params.Set("rev", rev)
  status, _, _ = docRes.DeleteJSON("", nil, params)
  return status == OK
}

// Set creates or updates a document with the specified ID.
func (d *Database)Set(docid string, doc map[string]interface{}) bool {
  if doc == nil {
    return false
  }

  docRes := docResource(d.resource, docid)
  status, _, data := docRes.PutJSON("", nil, doc, nil)
  if status != Created {
    return false
  }

  var jsonMap map[string]interface{}
  json.Unmarshal(*data, &jsonMap)
  doc["_id"] = jsonMap["id"].(string)
  doc["_rev"] = jsonMap["rev"].(string)
  return true
}

// DocIDs returns the IDs of all documents in database.
func (d *Database)DocIDs() []string {
  docRes := docResource(d.resource, "_all_docs")
  status, _, data := docRes.GetJSON("", nil, nil)
  if status != OK {
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

// Len returns the number of documents stored in it.
func (d *Database)Len() int {
  info := d.databaseInfo()
  if count, ok := info["doc_count"]; ok {
    return int(count.(float64))
  }
  return -1
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
func (d *Database)Save(doc map[string]interface{}) (string, string) {

  var id, rev string
  if doc == nil {
    return id, rev
  }

  var httpFunc func(string, *http.Header, map[string]interface{}, url.Values) (int, http.Header, *json.RawMessage)
  if v, ok := doc["_id"]; ok {
    httpFunc = docResource(d.resource, v.(string)).PutJSON
  } else {
    httpFunc = d.resource.PostJSON
  }

  _, _, data := httpFunc("", nil, doc, nil)
  var jsonMap map[string]interface{}
  json.Unmarshal(*data, &jsonMap)

  if v, ok := jsonMap["id"]; ok {
    id = v.(string)
    doc["_id"] = id
  }

  if v, ok := jsonMap["rev"]; ok {
    rev = v.(string)
    doc["_rev"] = rev
  }

  return id, rev
}

// docResource returns a Resource instance for docID
func docResource(res *Resource, docID string) *Resource {
  if docID[:1] == "_" {
    paths := strings.SplitN(docID, "/", 2)
    for _, p := range paths {
      res = res.NewResourceWithURL(p)
    }
    return res
  }

  return res.NewResourceWithURL(docID)
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
func (d *Database)Commit() bool {
  status, _, _ := d.resource.PostJSON("_ensure_full_commit", nil, nil, nil)
  return status == Created
}

// GetAttachment returns the file attachment associated with the document.
// The raw data of the associated attachment is returned as a []byte.
func (d *Database)GetAttachment(docid, fileName string) ([]byte, bool) {
  // defensive check
  if len(docid) == 0 || len(fileName) == 0 {
    return nil, false
  }

  docRes := docResource(docResource(d.resource, docid), fileName)
  status, _, data := docRes.Get("", nil, nil)
  return data, status == OK
}

// PutAttachment uploads the supplied *os.File as an attachment to the specified document.
// doc: the document that the attachment belongs to. Must have _id and _rev inside.
func (d *Database)PutAttachment(doc map[string]interface{}, file *os.File, mimeType string) bool {
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

  status, _, data := docRes.Put("", &header, contents, params)
  if status == Created {
    var jsonMap map[string]interface{}
    json.Unmarshal(data, &jsonMap)
    doc["_rev"] = jsonMap["rev"].(string)
  }

  return status == Created
}

// DeleteAttachment deletes the specified attachment
func (d *Database)DeleteAttachment(doc map[string]interface{}, fileName string) bool {
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
  status, _, data := docRes.DeleteJSON("", nil, params)
  if status == OK {
    var jsonMap map[string]interface{}
    json.Unmarshal(*data, &jsonMap)
    doc["_rev"] = jsonMap["rev"]
  }
  return status == OK
}

type IDRev struct {
  Id string
  Rev string
}

// UpdateDocuments performs a bulk update or creation of the given documents in a single HTTP request.
func (d *Database)UpdateDocuments(docs []map[string]interface{}, options map[string]interface{}) ([]IDRev, bool) {
  results := []IDRev{}

  if docs == nil {
    return results, false
  }

  body := map[string]interface{}{}
  if options != nil {
    for k, v := range options {
      body[k] = v
    }
  }
  body["docs"] = docs

  status, _, data := d.resource.PostJSON("_bulk_docs", nil, body, nil)
  if status == Created {
    var jsonArr []map[string]interface{}
    json.Unmarshal(*data, &jsonArr)
    for _, ele := range jsonArr {
      id, rev := ele["id"].(string), ele["rev"].(string)
      results = append(results, IDRev{Id: id, Rev: rev})
    }
  }
  return results, status == Created
}

// GetRevsLimit gets the current revs_limit(revision limit) setting.
func (d *Database)GetRevsLimit() (int, bool) {
  status, _, data := d.resource.Get("_revs_limit", nil, nil)
  limit, err := strconv.Atoi(strings.Trim(string(data), "\n"))
  return limit, status == OK && err == nil
}

// SetRevsLimit sets the maximum number of document revisions that will be
// tracked by CouchDB.
func (d *Database)SetRevsLimit(limit int) bool {
  status, _, _ := d.resource.Put("_revs_limit", nil, []byte(strconv.Itoa(limit)), nil)
  return status == OK
}

// Changes returns a sorted list of changes feed made to documents in the database.
func (d *Database)Changes(options url.Values) (map[string]interface{}, bool) {
  status, _, data := d.resource.GetJSON("_changes", nil, options)
  if status != OK {
    return nil, false
  }
  var changes map[string]interface{}
  json.Unmarshal(*data, &changes)
  return changes, status == OK
}

// Cleanup removes all view index files no longer required by CouchDB.
func (d *Database)Cleanup() bool {
  status, _, _ := d.resource.PostJSON("_view_cleanup", nil, nil, nil)
  return status == Accepted
}

// Compact compacts the database by compressing the disk database file.
func (d *Database)Compact() bool {
  status, _, _ := d.resource.PostJSON("_compact", nil, nil, nil)
  return status == Accepted
}

// Copy copies an existing document to a new or existing document.
func (d *Database)Copy(srcID, destID string) (string, bool) {
  docRes := docResource(d.resource, srcID)
  header := &http.Header{
    "Destination": []string{destID},
  }
  status, _, data := request("COPY", docRes.base, header, nil, nil)
  var rev string
  if status == Created {
    var jsonMap map[string]interface{}
    json.Unmarshal(data, &jsonMap)
    rev = jsonMap["rev"].(string)
  }

  return rev, status == Created
}

// Purge performs complete removing of the given documents.
func (d *Database)Purge(docIDs []string) bool {
  // TODO
  return false
}

func (d *Database)SetSecurity(securityDoc map[string]interface{}) bool {
  status, _, _ := d.resource.PutJSON("_security", nil, securityDoc, nil)
  return status == OK
}

func (d *Database)GetSecurity() (map[string]interface{}, bool) {
  status, _, data := d.resource.GetJSON("_security", nil, nil)
  var secDoc map[string]interface{}
  if status == OK {
    json.Unmarshal(*data, &secDoc)
  }
  return secDoc, status == OK
}
