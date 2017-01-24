CouchDB-Golang Library v1.2
===========================================

A Golang library for working with CouchDB 2.x

supported Golang version:
* 1.7.x

[![Build Status](https://travis-ci.org/leesper/couchdb-golang.svg?branch=master)](https://travis-ci.org/leesper/couchdb-golang)

* Resource : a simple wrapper for HTTP requests and error handling
* Server : CouchDB server instance
* Database : CouchDB database instance
* ViewResults : a representation of the results produced by design document views
* ViewDefinition : a definition of view stored in a specific design document
* Document : a representation of document object in database

Inspired by [CouchDB-Python](https://github.com/djc/couchdb-python)

# Documentation
```go
package couchdb
import "github.com/leesper/couchdb"
```

## Constants
```go
const (
    // DefaultBaseURL is the default address of CouchDB server.
    DefaultBaseURL = "http://localhost:5984"
)
```

## Variables
```go
var (
    // ErrBatchValue for invalid batch parameter of IterView
    ErrBatchValue = errors.New("batch must be 1 or more")
    // ErrLimitValue for invalid limit parameter of IterView
    ErrLimitValue = errors.New("limit must be 1 or more")
)

var (
    // ErrSetID for setting ID to document which already has one.
    ErrSetID = errors.New("id can only be set on new documents")
    // ErrNotStruct for not a struct value
    ErrNotStruct = errors.New("value not of struct type")
    // ErrNotDocumentEmbedded for not a document-embedded value
    ErrNotDocumentEmbedded = errors.New("value not Document-embedded")
)

var (

    // ErrNotModified for HTTP status code 304
    ErrNotModified = errors.New("status 304 - not modified")
    // ErrBadRequest for HTTP status code 400
    ErrBadRequest = errors.New("status 400 - bad request")
    // ErrUnauthorized for HTTP status code 401
    ErrUnauthorized = errors.New("status 401 - unauthorized")
    // ErrForbidden for HTTP status code 403
    ErrForbidden = errors.New("status 403 - forbidden")
    // ErrNotFound for HTTP status code 404
    ErrNotFound = errors.New("status 404 - not found")
    // ErrResourceNotAllowed for HTTP status code 405
    ErrResourceNotAllowed = errors.New("status 405 - resource not allowed")
    // ErrNotAcceptable for HTTP status code 406
    ErrNotAcceptable = errors.New("status 406 - not acceptable")
    // ErrConflict for HTTP status code 409
    ErrConflict = errors.New("status 409 - conflict")
    // ErrPreconditionFailed for HTTP status code 412
    ErrPreconditionFailed = errors.New("status 412 - precondition failed")
    // ErrBadContentType for HTTP status code 415
    ErrBadContentType = errors.New("status 415 - bad content type")
    // ErrRequestRangeNotSatisfiable for HTTP status code 416
    ErrRequestRangeNotSatisfiable = errors.New("status 416 - requested range not satisfiable")
    // ErrExpectationFailed for HTTP status code 417
    ErrExpectationFailed = errors.New("status 417 - expectation failed")
    // ErrInternalServerError for HTTP status code 500
    ErrInternalServerError = errors.New("status 500 - internal server error")
)
```

## func FromJSONCompatibleMap
```go
func FromJSONCompatibleMap(obj interface{}, docMap map[string]interface{}) error
```
FromJSONCompatibleMap constructs a Document-embedded struct from a JSON-compatible map.

## func GenerateUUID
```go
func GenerateUUID() string
```
GenerateUUID returns a random 128-bit UUID

## func Load
```go
func Load(db *Database, docID string, obj interface{}) error
```
Load loads the document in specified database.

## func Store
```go
func Store(db *Database, obj interface{}) error
```
Store stores the document in specified database.

## func SyncMany
```go
func SyncMany(db *Database, viewDefns []*ViewDefinition, removeMissing bool, callback func(map[string]interface{})) ([]UpdateResult, error)
```
SyncMany ensures that the views stored in the database match the views defined by the corresponding view definitions. This function might update more than one design document. This is done using CouchDB's bulk update to ensure atomicity of the opeation. db: the corresponding database.

viewDefns: a sequence of \*ViewDefinition instances.

removeMissing: whether to remove views found in a design document that are not found in the list of ViewDefinition instances, default false.

callback: a callback function invoked when a design document gets updated; it is called before the doc has actually been saved back to the database.

## func ToJSONCompatibleMap
```go
func ToJSONCompatibleMap(obj interface{}) (map[string]interface{}, error)
```
ToJSONCompatibleMap converts a Document-embedded struct into a JSON-compatible map, e.g. anything that cannot be jsonified will be ignored silently.

## type Database
```go
type Database struct {
    // contains filtered or unexported fields
}
```
Database represents a CouchDB database instance.

### func NewDatabase
```go
func NewDatabase(urlStr string) (*Database, error)
```
NewDatabase returns a CouchDB database instance.

### func NewDatabaseWithResource
```go
func NewDatabaseWithResource(res *Resource) (*Database, error)
```
NewDatabaseWithResource returns a CouchDB database instance with resource obj.

### func (d \*Database) Available
```go
func (d *Database) Available() error
```
Available returns error if the database is not good to go.

### func (d \*Database) Changes
```go
func (d *Database) Changes(options url.Values) (map[string]interface{}, error)
```
Changes returns a sorted list of changes feed made to documents in the database.

### func (d \*Database) Cleanup
```go
func (d *Database) Cleanup() error
```
Cleanup removes all view index files no longer required by CouchDB.

### func (d \*Database) Commit
```go
func (d *Database) Commit() error
```
Commit flushes any recent changes to the specified database to disk. If the server is configured to delay commits or previous requests use the special "X-Couch-Full-Commit: false" header to disable immediate commits, this method can be used to ensure that non-commited changes are commited to physical storage.

### func (d \*Database) Compact
```go
func (d *Database) Compact() error
```
Compact compacts the database by compressing the disk database file.

### func (d \*Database) Contains
```go
func (d *Database) Contains(docid string) error
```
Contains returns true if the database contains a document with the specified ID.

### func (d \*Database) Copy
```go
func (d *Database) Copy(srcID, destID, destRev string) (string, error)
```
Copy copies an existing document to a new or existing document.

### func (d \*Database) Delete
```go
func (d *Database) Delete(docid string) error
```
Delete deletes the document with the specified ID.

### func (d \*Database) DeleteAttachment
```go
func (d *Database) DeleteAttachment(doc map[string]interface{}, name string) error
```
DeleteAttachment deletes the specified attachment.

### func (d \*Database) DeleteDoc
```go
func (d *Database) DeleteDoc(doc map[string]interface{}) error
```
DeleteDoc deletes the specified document

### func (d \*Database) DeleteIndex
```go
func (d *Database) DeleteIndex(ddoc, name string) error
```
DeleteIndex deletes index in database.

### func (d \*Database) DocIDs
```go
func (d *Database) DocIDs() ([]string, error)
```
DocIDs returns the IDs of all documents in database.

### func (d \*Database) Get
```go
func (d *Database) Get(docid string, options url.Values) (map[string]interface{}, error)
```
Get returns the document with the specified ID.

### func (d \*Database) GetAttachment
```go
func (d *Database) GetAttachment(doc map[string]interface{}, name string) ([]byte, error)
```
GetAttachment returns the file attachment associated with the document. The raw data is returned as a []byte.

### func (d \*Database) GetAttachmentID
```go
func (d *Database) GetAttachmentID(docid, name string) ([]byte, error)
```
GetAttachmentID returns the file attachment associated with the document ID. The raw data is returned as []byte.

### func (d \*Database) GetIndex
```go
func (d *Database) GetIndex() (map[string]*json.RawMessage, error)
```
GetIndex gets all indexes created in database.

### func (d \*Database) GetRevsLimit
```go
func (d *Database) GetRevsLimit() (int, error)
```
GetRevsLimit gets the current revs_limit(revision limit) setting.

### func (d \*Database) GetSecurity
```go
func (d *Database) GetSecurity() (map[string]interface{}, error)
```
GetSecurity returns the current security object from the given database.

### func (d \*Database) Info
```go
func (d *Database) Info() (map[string]interface{}, error)
```
Info returns the information about the database.

### func (\*Database) IterView
```go
func (d *Database) IterView(name string, batch int, wrapper func(Row) Row, options map[string]interface{}) (<-chan Row, error)
```
IterView returns a channel fetching rows in batches which iterates a row at a time(pagination).

name: the name of the view, for user-defined views use the format "design_docid/viewname", that is, the document ID of the design document and the name of the view, separated by a /.

wrapper: an optional function for processing the result rows after retrieved.

options: optional query parameters.

### func (d \*Database) Len
```go
func (d *Database) Len() (int, error)
```
Len returns the number of documents stored in it.

### func (\*Database) List
```go
func (d *Database) List(name, view string, options map[string]interface{}) (http.Header, []byte, error)
```
List formats a view using a server-side 'list' function.

name: the name of the list function in the format "designdoc/listname"

view: the name of the view in the format "designdoc/viewname"

options: optional query parameters

### func (d \*Database) Name
```go
func (d *Database) Name() (string, error)
```
Name returns the name of database.

### func (d \*Database) Purge
```go
func (d *Database) Purge(docs []map[string]interface{}) (map[string]interface{}, error)
```
Purge performs complete removing of the given documents.

### func (d \*Database) PutAttachment
```go
func (d *Database) PutAttachment(doc map[string]interface{}, content []byte, name, mimeType string) error
```
PutAttachment uploads the supplied []byte as an attachment to the specified document. doc: the document that the attachment belongs to. Must have \_id and \_rev inside. content: the data to be attached to doc. name: name of attachment. mimeType: MIME type of content.

### func (d \*Database) PutIndex
```go
func (d *Database) PutIndex(indexFields []string, ddoc, name string) (string, string, error)
```
PutIndex creates a new index in database. indexFields: a JSON array of field names following the sort syntax. ddoc: optional, name of the design document in which the index will be created. By default each index will be created in its own design document. Indexes can be grouped into design documents for efficiency. However a change to one index in a design document will invalidate all other indexes in the same document. name: optional, name of the index. A name generated automatically if not provided.

### func (d \*Database) Query
```go
func (d *Database) Query(fields []string, selector string, sorts []string, limit, skip, index interface{}) ([]map[string]interface{}, error)
```
Query returns documents using a conditional selector statement in Golang.

selector: A filter string declaring which documents to return, formatted as a Golang statement.

fields: Specifying which fields to be returned, if passing nil the entire is returned, no automatic inclusion of \_id or other metadata fields.

sorts: How to order the documents returned, formatted as ["desc(fieldName1)", "desc(fieldName2)"] or ["fieldNameA", "fieldNameB"] of which "asc" is used by default, passing nil to disable ordering.

limit: Maximum number of results returned, passing nil to use default value(25).

skip: Skip the first 'n' results, where 'n' is the number specified, passing nil for no-skip.

index: Instruct a query to use a specific index, specified either as "<design_document>" or ["<design_document>", "<index_name>"], passing nil to use primary index(\_all_docs) by default.

## Inner functions for selector syntax

*nor(condexprs...)* matches if none of the conditions in condexprs match($nor).

For example: nor(year == 1990, year == 1989, year == 1997) returns all documents whose year field not in 1989, 1990 and 1997.

*all(field, array)* matches an array value if it contains all the elements of the argument array($all).

For example: all(genre, []string{"Comedy", "Short"} returns all documents whose genre field contains "Comedy" and "Short".

*any(field, condexpr)* matches an array field with at least one element meets the specified condition($elemMatch).

For example: any(genre, genre == "Short" || genre == "Horror") returns all documents whose genre field contains "Short" or "Horror" or both.

*exists(field, boolean)* checks whether the field exists or not, regardless of its value($exists).

For example: exists(director, false) returns all documents who does not have a director field.

*typeof(field, type)* checks the document field's type, valid types are "null", "boolean", "number", "string", "array", "object"($type).

For example: typeof(genre, "array") returns all documents whose genre field is of array type.

*in(field, array)* the field must exist in the array provided($in).

For example: in(director, []string{"Mike Portnoy", "Vitali Kanevsky"}) returns all documents whose director field is "Mike Portnoy" or "Vitali Kanevsky".

*nin(field, array)* the document field must not exist in the array provided($nin).

For example: nin(year, []int{1990, 1992, 1998}) returns all documents whose year field is not in 1990, 1992 or 1998.

*size(field, int)* matches the length of an array field in a document($size).

For example: size(genre, 2) returns all documents whose genre field is of length 2.

*mod(field, divisor, remainder)* matches documents where field % divisor == remainder($mod).

For example: mod(year, 2, 1) returns all documents whose year field is an odd number.

*regex(field, regexstr)* a regular expression pattern to match against the document field.

For example: regex(title, "^A") returns all documents whose title is begin with an "A".

##Inner functions for sort syntax

*asc(field)* sorts the field in ascending order, this is the default option while desc(field) sorts the field in descending order.

### func (d \*Database) QueryJSON
```go
func (d *Database) QueryJSON(query string) ([]map[string]interface{}, error)
```
QueryJSON returns documents using a declarative JSON querying syntax.

### func (d \*Database) Revisions
```go
func (d *Database) Revisions(docid string, options url.Values) ([]map[string]interface{}, error)
```
Revisions returns all available revisions of the given document in reverse order, e.g. latest first.

### func (d \*Database) Save
```go
func (d *Database) Save(doc map[string]interface{}, options url.Values) (string, string, error)
```
Save creates a new document or update an existing document. If doc has no \_id the server will generate a random UUID and a new document will be created. Otherwise the doc's \_id will be used to identify the document to create or update. Trying to update an existing document with an incorrect \_rev will cause failure. *NOTE* It is recommended to avoid saving doc without \_id and instead generate document ID on client side. To avoid such problems you can generate a UUID on the client side. GenerateUUID provides a simple, platform-independent implementation. You can also use other third-party packages instead. doc: the document to create or update.

### func (d \*Database) Set
```go
func (d *Database) Set(docid string, doc map[string]interface{}) error
```
Set creates or updates a document with the specified ID.

### func (d \*Database) SetRevsLimit
```go
func (d *Database) SetRevsLimit(limit int) error
```
SetRevsLimit sets the maximum number of document revisions that will be tracked by CouchDB.

### func (d \*Database) SetSecurity
```go
func (d *Database) SetSecurity(securityDoc map[string]interface{}) error
```
SetSecurity sets the security object for the given database.

### func (\*Database) Show
```go
func (d *Database) Show(name, docID string, params url.Values) (http.Header, []byte, error)
```
Show calls a server-side 'show' function.

name: the name of the show function in the format "designdoc/showname"

docID: optional document ID to pass to the show function

params: optional query parameters

### func (d \*Database) String
```go
func (d *Database) String() string
```

### func (d \*Database) Update
```go
func (d *Database) Update(docs []map[string]interface{}, options map[string]interface{}) ([]UpdateResult, error)
```
Update performs a bulk update or creation of the given documents in a single HTTP request. It returns a 3-tuple (id, rev, error)

### func (\*Database) UpdateDoc
```go
func (d *Database) UpdateDoc(name, docID string, params url.Values) (http.Header, []byte, error)
```
UpdateDoc calls server-side update handler.

name: the name of the update handler function in the format "designdoc/updatename".

docID: optional document ID to pass to the show function

params: optional query parameters

### func (\*Database) View
```go
func (d *Database) View(name string, wrapper func(Row) Row, options map[string]interface{}) (*ViewResults, error)
```
View executes a predefined design document view and returns the results.

name: the name of the view, for user-defined views use the format "design_docid/viewname", that is, the document ID of the design document and the name of the view, separated by a /.

wrapper: an optional function for processing the result rows after retrieved.

options: optional query parameters.

## type Document
```go
type Document struct {
    ID  string `json:"_id,omitempty"`  // for json only, call SetID/GetID instead
    Rev string `json:"_rev,omitempty"` // for json only, call GetRev instead
    // contains filtered or unexported fields
}
```
Document represents a document object in database.

### func DocumentWithID
```go
func DocumentWithID(id string) Document
```
DocumentWithID returns a new Document with ID.

### func (\*Document) GetID
```go
func (d *Document) GetID() string
```
GetID returns the document ID.

### func (\*Document) GetRev
```go
func (d *Document) GetRev() string
```
GetRev returns the document revision.

### func (\*Document) SetID
```go
func (d *Document) SetID(id string) error
```
SetID sets ID for new document or return error.

### func (\*Document) SetRev
```go
func (d *Document) SetRev(rev string)
```
SetRev sets revision for document.

## type Resource
```go
type Resource struct {
    // contains filtered or unexported fields
}
```
Resource handles all requests to CouchDB.

### func NewResource
```go
func NewResource(urlStr string, header http.Header) (*Resource, error)
```
NewResource returns a newly-created Resource instance.

### func (r \*Resource) Delete
```go
func (r *Resource) Delete(path string, header http.Header, params url.Values) (http.Header, []byte, error)
```
Delete is a wrapper around http.Delete.

### func (r \*Resource) DeleteJSON
```go
func (r *Resource) DeleteJSON(path string, header http.Header, params url.Values) (http.Header, []byte, error)
```
DeleteJSON issues a DELETE to the specified URL, with data returned as json.

### func (r \*Resource) Get
```go
func (r *Resource) Get(path string, header http.Header, params url.Values) (http.Header, []byte, error)
```
Get is a wrapper around http.Get.

### func (r \*Resource) GetJSON
```go
func (r *Resource) GetJSON(path string, header http.Header, params url.Values) (http.Header, []byte, error)
```
GetJSON issues a GET to the specified URL, with data returned as json.

### func (r \*Resource) Head
```go
func (r *Resource) Head(path string, header http.Header, params url.Values) (http.Header, []byte, error)
```
Head is a wrapper around http.Head.

### func (r \*Resource) NewResourceWithURL
```go
func (r *Resource) NewResourceWithURL(resStr string) (*Resource, error)
```
NewResourceWithURL returns newly created \*Resource combined with resource string.

### func (r \*Resource) Post
```go
func (r *Resource) Post(path string, header http.Header, body []byte, params url.Values) (http.Header, []byte, error)
```
Post is a wrapper around http.Post.

### func (r \*Resource) PostJSON
```go
func (r *Resource) PostJSON(path string, header http.Header, body map[string]interface{}, params url.Values) (http.Header, []byte, error)
```
PostJSON issues a POST to the specified URL, with data returned as json.

### func (r \*Resource) Put
```go
func (r *Resource) Put(path string, header http.Header, body []byte, params url.Values) (http.Header, []byte, error)
```
Put is a wrapper around http.Put.

### func (r \*Resource) PutJSON
```go
func (r *Resource) PutJSON(path string, header http.Header, body map[string]interface{}, params url.Values) (http.Header, []byte, error)
```
PutJSON issues a PUT to the specified URL, with data returned as json.

## type Row
```go
type Row struct {
    ID  string
    Key interface{}
    Val interface{}
    Doc interface{}
    Err error
}
```
Row represents a row returned by database views.

### func (Row) String
```go
func (r Row) String() string
```
String returns a string representation for Row

## type Server
```go
type Server struct {
    // contains filtered or unexported fields
}
```
Server represents a CouchDB server instance.

### func NewServer
```go
func NewServer(urlStr string) (*Server, error)
```
NewServer creates a CouchDB server instance in address urlStr.

### func NewServerNoFullCommit
```go
func NewServerNoFullCommit(urlStr string) (*Server, error)
```
NewServerNoFullCommit creates a CouchDB server instance in address urlStr with X-Couch-Full-Commit disabled.

### func (s \*Server) ActiveTasks
```go
func (s *Server) ActiveTasks() ([]interface{}, error)
```
ActiveTasks lists of running tasks.

### func (s \*Server) AddUser
```go
func (s *Server) AddUser(name, password string, roles []string) (string, string, error)
```
AddUser adds regular user in authentication database. Returns id and rev of the registered user.

### func (s \*Server) Config
```go
func (s *Server) Config(node string) (map[string]map[string]string, error)
```
Config returns the entire CouchDB server configuration as JSON structure.

### func (s \*Server) Contains
```go
func (s *Server) Contains(name string) bool
```
Contains returns true if a db with given name exsited.

### func (s \*Server) Create
```go
func (s *Server) Create(name string) (*Database, error)
```
Create returns a database instance with the given name, returns true if created, if database already existed, returns false, \*Database will be nil if failed.

### func (s \*Server) DBs
```go
func (s *Server) DBs() ([]string, error)
```
DBs returns a list of all the databases in the CouchDB server instance.

### func (s \*Server) Delete
```go
func (s *Server) Delete(name string) error
```
Delete deletes a database with the given name. Return false if failed.

### func (s \*Server) Get
```go
func (s *Server) Get(name string) (*Database, error)
```
Get gets a database instance with the given name. Return nil if failed.

### func (s \*Server) Len
```go
func (s *Server) Len() (int, error)
```
Len returns the number of dbs in CouchDB server instance.

### func (s \*Server) Login
```go
func (s *Server) Login(name, password string) (string, error)
```
Login regular user in CouchDB, returns authentication token.

### func (s \*Server) Logout
```go
func (s *Server) Logout(token string) error
```
Logout regular user in CouchDB.

### func (s \*Server) Membership
```go
func (s *Server) Membership() ([]string, []string, error)
```
Membership displays the nodes that are part of the cluster as clusterNodes. The field allNodes displays all nodes this node knows about, including the ones that are part of cluster.

### func (s \*Server) RemoveUser
```go
func (s *Server) RemoveUser(name string) error
```
RemoveUser removes regular user in authentication database.

### func (s \*Server) Replicate
```go
func (s *Server) Replicate(source, target string, options map[string]interface{}) (map[string]interface{}, error)
```
Replicate requests, configure or stop a replication operation.

### func (s \*Server) Stats
```go
func (s *Server) Stats(node, entry string) (map[string]interface{}, error)
```
Stats returns a JSON object containing the statistics for the running server.

### func (s \*Server) String
```go
func (s *Server) String() string
```

### func (s \*Server) UUIDs
```go
func (s *Server) UUIDs(count int) ([]string, error)
```
UUIDs requests one or more Universally Unique Identifiers from the CouchDB instance. The response is a JSON object providing a list of UUIDs. count - Number of UUIDs to return. Default is 1.

### func (s \*Server) VerifyToken
```go
func (s *Server) VerifyToken(token string) error
```
VerifyToken returns error if user's token is invalid.

### func (s \*Server) Version
```go
func (s *Server) Version() (string, error)
```
Version returns the version info about CouchDB instance.

## type UpdateResult
```go
type UpdateResult struct {
    ID, Rev string
    Err     error
}
```
UpdateResult represents result of an update.

## type ViewDefinition
```go
type ViewDefinition struct {
    // contains filtered or unexported fields
}
```
ViewDefinition is a definition of view stored in a specific design document.

### func NewViewDefinition
```go
func NewViewDefinition(design, name, mapFun, reduceFun, language string, wrapper func(Row) Row, options map[string]interface{}) (*ViewDefinition, error)
```
NewViewDefinition returns a newly-created \*ViewDefinition. design: the name of the design document.

name: the name of the view.

mapFun: the map function code.

reduceFun: the reduce function code(optional).

language: the name of the programming language used, default is javascript.

wrapper: an optional function for processing the result rows after retrieved.

options: view specific options.

### func (\*ViewDefinition) GetDoc
```go
func (vd *ViewDefinition) GetDoc(db *Database) (map[string]interface{}, error)
```
GetDoc retrieves the design document corresponding to this view definition from the given database.

### func (\*ViewDefinition) Sync
```go
func (vd *ViewDefinition) Sync(db *Database) ([]UpdateResult, error)
```
Sync ensures that the view stored in the database matches the view defined by this instance.

### func (\*ViewDefinition) View
```go
func (vd *ViewDefinition) View(db *Database, options map[string]interface{}) (*ViewResults, error)
```
View executes the view definition in the given database.

## type ViewField
```go
type ViewField func() (*ViewDefinition, error)
```
ViewField represents a view definition value bound to Document.

### func NewViewField
```go
func NewViewField(design, name, mapFun, reduceFun, language string, wrapper func(Row) Row, options map[string]interface{}) ViewField
```
NewViewField returns a ViewField function. design: the name of the design document.

name: the name of the view.

mapFun: the map function code.

reduceFun: the reduce function code(optional).

language: the name of the programming language used, default is javascript.

wrapper: an optional function for processing the result rows after retrieved.

options: view specific options.

## type ViewResults
```go
type ViewResults struct {
    // contains filtered or unexported fields
}
```
ViewResults represents the results produced by design document views.

### func (\*ViewResults) Offset
```go
func (vr *ViewResults) Offset() (int, error)
```
Offset returns offset of ViewResults

### func (\*ViewResults) Rows
```go
func (vr *ViewResults) Rows() ([]Row, error)
```
Rows returns a slice of rows mapped (and reduced) by the view.

### func (\*ViewResults) TotalRows
```go
func (vr *ViewResults) TotalRows() (int, error)
```
TotalRows returns total rows of ViewResults

### func (\*ViewResults) UpdateSeq
```go
func (vr *ViewResults) UpdateSeq() (int, error)
```
UpdateSeq returns update sequence of ViewResults
