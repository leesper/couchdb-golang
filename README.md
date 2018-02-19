# CouchDB-Golang

[![Build Status](https://travis-ci.org/leesper/couchdb-golang.svg?branch=master)](https://travis-ci.org/leesper/couchdb-golang) [![GoDoc](https://godoc.org/github.com/leesper/couchdb-golang?status.svg)](http://godoc.org/github.com/leesper/couchdb-golang) [![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/leesper/couchdb-golang/master/LICENSE)

A Golang library for CouchDB 2.x, inspired by [CouchDB-Python](https://github.com/djc/couchdb-python).

## Features

* Resource : a simple wrapper for HTTP requests and error handling
* Server : CouchDB server instance
* Database : CouchDB database instance
* ViewResults : a representation of the results produced by design document views
* ViewDefinition : a definition of view stored in a specific design document
* Document : a representation of document object in database
* tools/replicate : a command-line tool for replicating

```go
func (d *Database) Query(fields []string, selector string, sorts []string, limit, skip, index interface{}) ([]map[string]interface{}, error)
```

You can query documents using a conditional selector statement in Golang. It will converts to the corresponding JSON query string.

* **selector**: A filter string declaring which documents to return, formatted as a Golang statement.
* **fields**: Specifying which fields to be returned, if passing nil the entire is returned, no automatic inclusion of \_id or other metadata fields.
* **sorts**: How to order the documents returned, formatted as ["desc(fieldName1)", "desc(fieldName2)"] or ["fieldNameA", "fieldNameB"] of which "asc" is used by default, passing nil to disable ordering.
* **limit**: Maximum number of results returned, passing nil to use default value(25).
* **skip**: Skip the first 'n' results, where 'n' is the number specified, passing nil for no-skip.
* **index**: Instruct a query to use a specific index, specified either as "<design_document>" or ["<design_document>", "<index_name>"], passing nil to use primary index(\_all_docs) by default.

For example:
```go
docsQuery, err := movieDB.Query(nil, `year == 1989 && (director == "Ademir Kenovic" || director == "Dezs Garas")`, nil, nil, nil, nil)
```
equals to:
```go
docsRaw, err := movieDB.QueryJSON(`
{
  "selector": {
    "year": 1989,
    "$or": [
      { "director": "Ademir Kenovic" },
      { "director": "Dezs Garas" }
    ]
  }
}`)
```

### Inner functions for selector syntax

* **nor(condexprs...)** matches if none of the conditions in condexprs match($nor).
For example: nor(year == 1990, year == 1989, year == 1997) returns all documents whose year field not in 1989, 1990 and 1997.

* **all(field, array)** matches an array value if it contains all the elements of the argument array($all).
For example: all(genre, []string{"Comedy", "Short"} returns all documents whose genre field contains "Comedy" and "Short".

* **any(field, condexpr)** matches an array field with at least one element meets the specified condition($elemMatch).
For example: any(genre, genre == "Short" || genre == "Horror") returns all documents whose genre field contains "Short" or "Horror" or both.

* **exists(field, boolean)** checks whether the field exists or not, regardless of its value($exists).
For example: exists(director, false) returns all documents who does not have a director field.

* **typeof(field, type)** checks the document field's type, valid types are "null", "boolean", "number", "string", "array", "object"($type).
For example: typeof(genre, "array") returns all documents whose genre field is of array type.

* **in(field, array)** the field must exist in the array provided($in).
For example: in(director, []string{"Mike Portnoy", "Vitali Kanevsky"}) returns all documents whose director field is "Mike Portnoy" or "Vitali Kanevsky".

* **nin(field, array)** the document field must not exist in the array provided($nin).
For example: nin(year, []int{1990, 1992, 1998}) returns all documents whose year field is not in 1990, 1992 or 1998.

* **size(field, int)** matches the length of an array field in a document($size).
For example: size(genre, 2) returns all documents whose genre field is of length 2.

* **mod(field, divisor, remainder)** matches documents where field % divisor == remainder($mod).
For example: mod(year, 2, 1) returns all documents whose year field is an odd number.

* **regex(field, regexstr)** a regular expression pattern to match against the document field.
For example: regex(title, "^A") returns all documents whose title is begin with an "A".

### Inner functions for sort syntax

**asc(field)** sorts the field in ascending order, this is the default option while desc(field) sorts the field in descending order.

## Requirements

* Golang 1.7.x and above

## Installation

`go get -u -v github.com/leesper/couchdb-golang`

## Authors and acknowledgment

* [Philipp Winter](https://github.com/philippwinter)
* [Serkan Sipahi](https://github.com/SerkanSipahi)
* [paduraru](https://github.com/paduraru)
* [Andrei Pavel](https://github.com/andreipavelQ)
* [jcantonio](https://github.com/jcantonio)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Please make sure to update unit tests as appropriate.
