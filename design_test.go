package couchdb

import (
	"net/url"
	"testing"
)

func TestRowObject(t *testing.T) {
	results, err := designDB.View("_all_docs", nil, url.Values{"keys": []string{"blah"}})
	if err != nil {
		t.Error("db view error", err)
	}

	rows, err := results.Rows()
	if err != nil {
		t.Error("rows error", err)
	}

	row := rows[0]
	if row.ID != "" {
		t.Error("row ID not empty", row.ID)
	}

	if row.Key.(string) != "blah" {
		t.Errorf("row key %s want blah", row.Key.(string))
	}

	if row.Val != nil {
		t.Error("row value not nil", row.Val)
	}

	if row.Err.Error() != "not_found" {
		t.Errorf("row error %s want not_found", row.Err)
	}

	_, _, err = designDB.Save(map[string]interface{}{"_id": "xyz", "foo": "bar"}, nil)
	if err != nil {
		t.Error("db save error", err)
	}

	results, err = designDB.View("_all_docs", nil, url.Values{"keys": []string{"xyz"}})
	if err != nil {
		t.Error("db view error", err)
	}

	rows, err = results.Rows()
	if err != nil {
		t.Error("rows error", err)
	}

	row = rows[0]
	if row.ID != "xyz" {
		t.Errorf("row ID %s want xyz", row.ID)
	}

	if row.Key.(string) != "xyz" {
		t.Errorf("row key %s want xyz", row.Key)
	}

	value := row.Val.(map[string]interface{})
	_, ok := value["rev"]
	if !(ok && len(value) == 1) {
		t.Error("row value not contains rev only")
	}

	if row.Err != nil {
		t.Error("row error not nil", row.Err)
	}
}

func TestViewMultiGet(t *testing.T)        {}
func TestDesignDocInfo(t *testing.T)       {}
func TestViewCompaction(t *testing.T)      {}
func TestViewCleanup(t *testing.T)         {}
func TestViewFunctionObjects(t *testing.T) {}
func TestInitWithResource(t *testing.T)    {}
func TestIterView(t *testing.T)            {}
func TestUpdateSeq(t *testing.T)           {}
func TestTmpviewRepr(t *testing.T)         {}
func TestWrapperIter(t *testing.T)         {}
func TestWrapperRows(t *testing.T)         {}
func TestProperties(t *testing.T)          {}
func TestRowRepr(t *testing.T)             {}

//
func TestAllRows(t *testing.T)            {}
func TestBatchSizes(t *testing.T)         {}
func TestBatchSizesWithSkip(t *testing.T) {}
func TestLimit(t *testing.T)              {}
func TestDescending(t *testing.T)         {}
func TestStartKey(t *testing.T)           {}
func TestNullKeys(t *testing.T)           {}
