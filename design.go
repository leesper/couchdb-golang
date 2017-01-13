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
	wrapper   func(Row) Row
	Offset    int
	TotalRows int
}

// Rows returns a slice of rows mapped (and reduced) by the view.
func (vr *ViewResults) Rows() ([]Row, error) {
	keys, ok := vr.options["keys"]
	var data []byte
	var err error
	if ok {
		body := map[string]interface{}{"keys": keys}
		options := url.Values{}
		for k, v := range vr.options {
			if k != "keys" {
				options[k] = v
			}
		}
		_, data, err = vr.resource.PostJSON(vr.designDoc, nil, body, options)
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
	vr.TotalRows = int(totalRows)

	var offset float64
	json.Unmarshal(jsonMap["offset"], &offset)
	vr.Offset = int(offset)

	rowsRaw := []json.RawMessage{}
	json.Unmarshal(jsonMap["rows"], &rowsRaw)

	rows := make([]Row, len(rowsRaw))
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

		if vr.wrapper != nil {
			row = vr.wrapper(row)
		}
		rows[idx] = row
	}
	return rows, nil
}
