// CouchDB resource
//
// This is the low-level wrapper functions of HTTP methods
//
// Used for communicating with CouchDB Server
package couchdb

import (
  "http"
)

// Head is a wrapper around http.Head
func Head() {}

// Get is a wrapper around http.Get
func Get() {}

// Post is a wrapper around http.Post
func Post() {}

// Delete is a wrapper around http.Delete
func Delete() {}

// Put is a wrapper around http.Put
func Put() {}

// GetJSON issues a GET to the specified URL, with data returned as json
func GetJSON() {}

// PostJSON issues a POST to the specified URL, with data returned as json
func PostJSON() {}

// DeleteJSON issues a DELETE to the specified URL, with data returned as json
func DeleteJSON() {}

// PutJSON issues a PUT to the specified URL, with data returned as json
func PutJSON() {}
