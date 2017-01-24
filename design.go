package couchdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
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

func viewLikeResourceRequest(res *Resource, opts map[string]interface{}) (http.Header, []byte, error) {
	params := url.Values{}
	body := map[string]interface{}{}
	for key, val := range opts {
		switch key {
		case "keys": // json-array, put in body and send POST request
			body[key] = val
		case "key", "startkey", "start_key", "endkey", "end_key":
			data, err := json.Marshal(val)
			if err != nil {
				return nil, nil, err
			}
			params.Add(key, string(data))
		case "conflicts", "descending", "group", "include_docs", "attachments", "att_encoding_info", "inclusive_end", "reduce", "sorted", "update_seq":
			if val.(bool) {
				params.Add(key, "true")
			} else {
				params.Add(key, "false")
			}
		case "endkey_docid", "end_key_doc_id", "stale", "startkey_docid", "start_key_doc_id", "format": // format for _list request
			params.Add(key, val.(string))
		case "group_level", "limit", "skip":
			params.Add(key, fmt.Sprintf("%d", val))
		default:
			switch val := val.(type) {
			case bool:
				if val {
					params.Add(key, "true")
				} else {
					params.Add(key, "false")
				}
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				params.Add(key, fmt.Sprintf("%d", val))
			case float32, float64:
				params.Add(key, fmt.Sprintf("%f", val))
			default:
				return nil, nil, fmt.Errorf("value %v not supported", val)
			}
		}
	}

	if len(body) > 0 {
		return res.PostJSON("", nil, body, params)
	}

	return res.GetJSON("", nil, params)
}

func (vr *ViewResults) fetch() ([]Row, error) {
	res := docResource(vr.resource, vr.designDoc)
	_, data, err := viewLikeResourceRequest(res, vr.options)
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
type ViewDefinition struct {
	design, name, mapFun, reduceFun, language string
	wrapper                                   func(Row) Row
	options                                   map[string]interface{}
}

// NewViewDefinition returns a newly-created *ViewDefinition.
// design: the name of the design document.
//
// name: the name of the view.
//
// mapFun: the map function code.
//
// reduceFun: the reduce function code(optional).
//
// language: the name of the programming language used, default is javascript.
//
// wrapper: an optional function for processing the result rows after retrieved.
//
// options: view specific options.
func NewViewDefinition(design, name, mapFun, reduceFun, language string, wrapper func(Row) Row, options map[string]interface{}) (*ViewDefinition, error) {
	if language == "" {
		language = "javascript"
	}

	if mapFun == "" {
		return nil, errors.New("map function empty")
	}

	return &ViewDefinition{
		design:    design,
		name:      name,
		mapFun:    strings.TrimLeft(mapFun, "\n"),
		reduceFun: strings.TrimLeft(reduceFun, "\n"),
		language:  language,
		wrapper:   wrapper,
		options:   options,
	}, nil
}

// View executes the view definition in the given database.
func (vd *ViewDefinition) View(db *Database, options map[string]interface{}) (*ViewResults, error) {
	opts := deepCopy(options)
	for k, v := range vd.options {
		opts[k] = v
	}
	return db.View(fmt.Sprintf("%s/%s", vd.design, vd.name), nil, opts)
}

// GetDoc retrieves the design document corresponding to this view definition from
// the given database.
func (vd *ViewDefinition) GetDoc(db *Database) (map[string]interface{}, error) {
	if db == nil {
		return nil, errors.New("database nil")
	}
	return db.Get(fmt.Sprintf("_design/%s", vd.design), nil)
}

// Sync ensures that the view stored in the database matches the view defined by this instance.
func (vd *ViewDefinition) Sync(db *Database) ([]UpdateResult, error) {
	if db == nil {
		return nil, errors.New("database nil")
	}
	return SyncMany(db, []*ViewDefinition{vd}, false, nil)
}

// SyncMany ensures that the views stored in the database match the views defined
// by the corresponding view definitions. This function might update more than
// one design document. This is done using CouchDB's bulk update to ensure atomicity of the opeation.
// db: the corresponding database.
//
// viewDefns: a sequence of *ViewDefinition instances.
//
// removeMissing: whether to remove views found in a design document that are not
// found in the list of ViewDefinition instances, default false.
//
// callback: a callback function invoked when a design document gets updated;
// it is called before the doc has actually been saved back to the database.
func SyncMany(db *Database, viewDefns []*ViewDefinition, removeMissing bool, callback func(map[string]interface{})) ([]UpdateResult, error) {
	if db == nil {
		return nil, errors.New("database nil")
	}

	docs := []map[string]interface{}{}
	designs := map[string]bool{}
	defMap := map[string][]*ViewDefinition{}

	for _, dfn := range viewDefns {
		designs[dfn.design] = true
		if _, ok := defMap[dfn.design]; !ok {
			defMap[dfn.design] = []*ViewDefinition{}
		}
		defMap[dfn.design] = append(defMap[dfn.design], dfn)
	}

	orders := []string{}
	for k := range designs {
		orders = append(orders, k)
	}
	sort.Strings(orders)

	for _, design := range orders {
		docID := fmt.Sprintf("_design/%s", design)
		doc, err := db.Get(docID, nil)
		if err != nil {
			doc = map[string]interface{}{"_id": docID}
		}
		origDoc := deepCopy(doc)
		languages := map[string]bool{}

		missing := map[string]bool{}
		vs, ok := doc["views"]
		if ok {
			for k := range vs.(map[string]interface{}) {
				missing[k] = true
			}
		}

		for _, dfn := range defMap[design] {
			funcs := map[string]interface{}{"map": dfn.mapFun}
			if len(dfn.reduceFun) > 0 {
				funcs["reduce"] = dfn.reduceFun
			}
			if dfn.options != nil {
				funcs["options"] = dfn.options
			}
			_, ok = doc["views"]
			if ok {
				doc["views"].(map[string]interface{})[dfn.name] = funcs
			} else {
				doc["views"] = map[string]interface{}{dfn.name: funcs}
			}
			languages[dfn.language] = true
			if missing[dfn.name] {
				delete(missing, dfn.name)
			}
		}

		if removeMissing {
			for k := range missing {
				delete(doc["views"].(map[string]interface{}), k)
			}
		} else if _, ok := doc["language"]; ok {
			languages[doc["language"].(string)] = true
		}

		langs := []string{}
		for lang := range languages {
			langs = append(langs, lang)
		}

		if len(langs) > 1 {
			return nil, fmt.Errorf("found different language views in one design document %v", langs)
		}
		doc["language"] = langs[0]

		if !reflect.DeepEqual(doc, origDoc) {
			if callback != nil {
				callback(doc)
			}
			docs = append(docs, doc)
		}
	}

	return db.Update(docs, nil)
}

func deepCopy(src map[string]interface{}) map[string]interface{} {
	dst := map[string]interface{}{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
