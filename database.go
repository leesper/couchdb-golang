package couchdb

import (
  "crypto/rand"
  "encoding/json"
  "fmt"
  "net/http"
  "net/url"
  "os"
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
  _, _, jsonData := d.resource.GetJSON("", nil, url.Values{})

  if jsonData == nil {
    return ""
  }

  var jsonMap map[string]interface{}
  _ = json.Unmarshal(*jsonData, &jsonMap)
  if _, ok := jsonMap["db_name"]; !ok {
    return ""
  }

  return jsonMap["db_name"].(string)
}

// Save creates a new document or update an existing document.
// If doc has no _id the server will generate a random UUID and a new document will be created.
// Otherwise the doc's _id will be used to identify the document to create or update.
// Trying to update an existing document with an incorrect _rev will cause failure.
// *NOTE* It is recommended to avoid saving doc without _id and instead generate document ID on client side.
// To avoid such problems you can generate a UUID on the client side.
// GenerateUUID provides a simple, platform-independent implementation.
// You can also use other third-party packages instead.
// doc: the document to create or update
// options: optional args, such as batch='ok'
func (d *Database)Save(doc map[string]interface{}, options url.Values) (string, string) {
  var id, rev string
  var httpFunc func(string, *http.Header, map[string]interface{}, url.Values) (int, http.Header, *json.RawMessage)
  if v, ok := doc["_id"]; ok {
    httpFunc = docResource(d.resource, v.(string)).PutJSON
  } else {
    httpFunc = d.resource.PostJSON
  }

  _, _, data := httpFunc("", nil, doc, options)
  var jsonMap map[string]interface{}
  _ = json.Unmarshal(*data, &jsonMap)

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

  uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
  return uuid
}
