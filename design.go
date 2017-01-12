package couchdb

import (
	"encoding/json"
	"errors"
	"net/url"
)

// Row represents a row returned by database views.
type Row struct {
	ID  string
	Key interface{}
	Val interface{}
	Doc interface{}
	Err error
}

// ViewResults represents the results produced by design document views.
type ViewResults struct {
	resource  *Resource
	designDoc string
	options   url.Values
}

func (vr ViewResults) Rows() ([]Row, error) {
	keys, ok := vr.options["keys"]
	var data []byte
	var err error
	if ok {
		body := map[string]interface{}{"keys": keys}
		_, data, err = vr.resource.PostJSON(vr.designDoc, nil, body, vr.options)
		if err != nil {
			return nil, err
		}
	} else {
		_, data, err = vr.resource.GetJSON(vr.designDoc, nil, vr.options)
		if err != nil {
			return nil, err
		}
	}

	var jsonMap map[string]json.RawMessage
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return nil, err
	}

	var totalRows float64
	json.Unmarshal(jsonMap["total_rows"], &totalRows)

	rowsRaw := []json.RawMessage{}
	json.Unmarshal(jsonMap["rows"], &rowsRaw)

	rows := make([]Row, int(totalRows))
	var rowMap map[string]interface{}
	for idx, raw := range rowsRaw {
		json.Unmarshal(raw, &rowMap)
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
		rows[idx] = row
	}
}
