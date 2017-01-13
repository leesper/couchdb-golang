package couchdb

import "testing"

func TestRowObject(t *testing.T) {
	results, err := designDB.View("_all_docs", nil, map[string]interface{}{"keys": []string{"blah"}})
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

	results, err = designDB.View("_all_docs", nil, map[string]interface{}{"keys": []string{"xyz"}})
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

	designDB.Delete(row.ID)
}

func TestViewMultiGet(t *testing.T) {
	for i := 1; i < 6; i++ {
		designDB.Save(map[string]interface{}{"i": i}, nil)
	}

	designDB.Set("_design/test", map[string]interface{}{
		"language": "javascript",
		"views": map[string]interface{}{
			"multi_key": map[string]string{
				"map": "function(doc) { emit(doc.i, null); }",
			},
		},
	})

	results, err := designDB.View("test/multi_key", nil, map[string]interface{}{"keys": []int{1, 3, 5}})
	if err != nil {
		t.Error("db view error", err)
	}

	rows, err := results.Rows()
	if err != nil {
		t.Error("rows error", err)
	}

	if len(rows) != 3 {
		t.Errorf("rows length %d want 3", len(rows))
	}

	for idx, i := range []int{1, 3, 5} {
		if i != int(rows[idx].Key.(float64)) {
			t.Errorf("key = %d want %d", int(rows[idx].Key.(float64)), i)
		}
	}
}

func TestDesignDocInfo(t *testing.T) {
	designDB.Set("_design/test", map[string]interface{}{
		"language": "javascript",
		"views": map[string]interface{}{
			"test": map[string]string{"map": "function(doc) { emit(doc.type, null); }"},
		},
	})
	info, _ := designDB.Info("test")
	compactRunning := info["view_index"].(map[string]interface{})["compact_running"].(bool)
	if compactRunning {
		t.Error("compact running true want false")
	}
}

func TestViewCompaction(t *testing.T) {
	designDB.Set("_design/test", map[string]interface{}{
		"language": "javascript",
		"views": map[string]interface{}{
			"multi_key": map[string]string{"map": "function(doc) { emit(doc.i, null); }"},
		},
	})

	_, err := designDB.View("test/multi_key", nil, nil)
	if err != nil {
		t.Error("db view error", err)
	}
	err = designDB.Compact()
	if err != nil {
		t.Error("db compact error", err)
	}
}
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
