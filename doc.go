// Package couchdb provides components to work with CouchDB 2.x with Go.
//
// Resource is the low-level wrapper functions of HTTP methods
// used for communicating with CouchDB Server.
//
// Server contains all the functions to work with CouchDB server, including some
// basic functions to facilitate the basic user management provided by it.
//
// Database contains all the functions to work with CouchDB database, such as
// documents manipulating and querying.
//
// ViewResults represents the results produced by design document views. When calling
// any of its functions like Offset(), TotalRows(), UpdateSeq() or Rows(), it will
// perform a query on views on server side, and returns results as slice of Row
//
// ViewDefinition is a definition of view stored in a specific design document,
// you can define your own map-reduce functions and Sync with the database.
//
// Document represents a document object in database. All struct that can be mapped
// into CouchDB document must have it embedded. For example:
//
//  type User struct {
//    Name string `json:"name"`
//    Age int `json:"age"`
//    Document
//  }
//  user := User{"Mike", 18}
//  anotherUser := User{}
//
// Then you can call Store(db, &user) to store it into CouchDB or Load(db, user.GetID(), &anotherUser)
// to get the data from database.
//
// ViewField represents a view definition value bound to Document.
package couchdb
