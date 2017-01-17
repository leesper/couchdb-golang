package couchdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Row represents a row returned by database views.
type Row struct {
	ID  string
	Key interface{}
	Val interface{}
	Doc interface{}
	Err error
}

// String returns a string representation for Row
func (r Row) String() string {
	id := fmt.Sprintf("%s=%s", "id", r.ID)
	key := fmt.Sprintf("%s=%v", "key", r.Key)
	doc := fmt.Sprintf("%s=%v", "doc", r.Doc)
	estr := fmt.Sprintf("%s=%v", "err", r.Err)
	val := fmt.Sprintf("%s=%v", "val", r.Val)
	return fmt.Sprintf("<%s %s>", "Row", strings.Join([]string{id, key, doc, estr, val}, ", "))
}

// ViewResults represents the results produced by design document views.
type ViewResults struct {
	resource  *Resource
	designDoc string
	options   map[string]interface{}
	wrapper   func(Row) Row

	offset    int
	totalRows int
	updateSeq int
	rows      []Row
	err       error
}

// newViewResults returns a newly-allocated *ViewResults
func newViewResults(r *Resource, ddoc string, opt map[string]interface{}, wr func(Row) Row) *ViewResults {
	return &ViewResults{
		resource:  r,
		designDoc: ddoc,
		options:   opt,
		wrapper:   wr,
		offset:    -1,
		totalRows: -1,
		updateSeq: -1,
	}
}

// Offset returns offset of ViewResults
func (vr *ViewResults) Offset() (int, error) {
	if vr.rows == nil {
		vr.rows, vr.err = vr.fetch()
	}
	return vr.offset, vr.err
}

// TotalRows returns total rows of ViewResults
func (vr *ViewResults) TotalRows() (int, error) {
	if vr.rows == nil {
		vr.rows, vr.err = vr.fetch()
	}
	return vr.totalRows, vr.err
}

// UpdateSeq returns update sequence of ViewResults
func (vr *ViewResults) UpdateSeq() (int, error) {
	if vr.rows == nil {
		vr.rows, vr.err = vr.fetch()
	}
	return vr.updateSeq, vr.err
}

// Rows returns a slice of rows mapped (and reduced) by the view.
func (vr *ViewResults) Rows() ([]Row, error) {
	if vr.rows == nil {
		vr.rows, vr.err = vr.fetch()
	}
	return vr.rows, vr.err
}

func (vr *ViewResults) fetch() ([]Row, error) {
	var data []byte
	var err error
	body := map[string]interface{}{}
	params := url.Values{}
	for key, val := range vr.options {
		switch key {
		case "keys": // json-array, put it in body and send POST request
			body[key] = val
		case "key", "startkey", "start_key", "endkey", "end_key": // json
			data, err = json.Marshal(val)
			if err != nil {
				return nil, err
			}
			params.Add(key, string(data))
		case "conflicts", "descending", "group", "include_docs", "attachments", "att_encoding_info", "inclusive_end", "reduce", "sorted", "update_seq": // boolean
			if val.(bool) {
				params.Add(key, "true")
			} else {
				params.Add(key, "false")
			}
		case "endkey_docid", "end_key_doc_id", "stale", "startkey_docid", "start_key_doc_id": // string
			params.Add(key, val.(string))
		case "group_level", "limit", "skip": // number
			params.Add(key, fmt.Sprintf("%d", val.(int)))
		}
	}

	if len(body) > 0 {
		_, data, err = vr.resource.PostJSON(vr.designDoc, nil, body, params)
	} else {
		_, data, err = vr.resource.GetJSON(vr.designDoc, nil, params)
	}
	if err != nil {
		return nil, err
	}

	var jsonMap map[string]*json.RawMessage
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return nil, err
	}

	var totalRows float64
	json.Unmarshal(*jsonMap["total_rows"], &totalRows)
	vr.totalRows = int(totalRows)

	if offsetRaw, ok := jsonMap["offset"]; ok {
		var offset float64
		json.Unmarshal(*offsetRaw, &offset)
		vr.offset = int(offset)
	}

	if updateSeqRaw, ok := jsonMap["update_seq"]; ok {
		var updateSeq float64
		json.Unmarshal(*updateSeqRaw, &updateSeq)
		vr.updateSeq = int(updateSeq)
	}

	var rowsRaw []*json.RawMessage
	json.Unmarshal(*jsonMap["rows"], &rowsRaw)

	rows := make([]Row, len(rowsRaw))
	var rowMap map[string]interface{}
	for idx, raw := range rowsRaw {
		json.Unmarshal(*raw, &rowMap)
		row := Row{}
		if id, ok := rowMap["id"]; ok {
			row.ID = id.(string)
		}

		if key, ok := rowMap["key"]; ok {
			row.Key = key
		}

		if val, ok := rowMap["value"]; ok {
			row.Val = val
		}

		if errmsg, ok := rowMap["error"]; ok {
			row.Err = errors.New(errmsg.(string))
		}

		if doc, ok := rowMap["doc"]; ok {
			row.Doc = doc
		}

		if vr.wrapper != nil {
			row = vr.wrapper(row)
		}
		rows[idx] = row
	}
	return rows, nil
}

// ViewDefinition is a definition of view stored in a specific design document.
type ViewDefinition struct{}

// NewViewDefinition returns a newly-created *ViewDefinition.
// design: the name of the design document.
//
// name: the name of the view.
//
// mapFun: the map function code.
//
// reduceFun: the reduce function code(optional).
//
// language: the name of the programming language used.
//
// wrapper: an optional function for processing the result rows after retrieved.
//
// options: view specific options.
func NewViewDefinition(design, name, mapFun, reduceFun, language string, wrapper func(Row) Row, options map[string]interface{}) {
}

// GetDoc retrieves the design document corresponding to this view definition from
// the given database.
func (vd *ViewDefinition) GetDoc(db *Database) {}

// Sync ensures that the view stored in the database matches the view defined by this instance.
func (vd *ViewDefinition) Sync(db *Database) {}

// SyncMany ensures that the views stored in the database match the views defined
// by the corresponding view definitions. This function might update more than
// one design document. This is done using CouchDB's bulk update to ensure atomicity of the opeation.
// db: the corresponding database.
//
// viewDefns: a sequence of *ViewDefinition instances.
//
// removeMissing: whether to remove views found in a design document that are not
// found in the list of ViewDefinition instances.
//
// callback: a callback function invoked when a design document gets updated;
// it is called before the doc has actually been saved back to the database.
func SyncMany(db *Database, viewDefns []*ViewDefinition, removeMissing bool, callback func(map[string]interface{})) {
}
