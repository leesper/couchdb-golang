package couchdb

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

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

func TestViewCleanup(t *testing.T) {
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

	ddoc, err := designDB.Get("_design/test", nil)
	if err != nil {
		t.Error("db get error", err)
	}
	ddoc["views"] = map[string]interface{}{
		"ids": map[string]string{"map": "function(doc) { emit(doc._id, null); }"},
	}
	_, err = designDB.Update([]map[string]interface{}{ddoc}, nil)
	if err != nil {
		t.Error("db update error", err)
	}

	designDB.View("test/ids", nil, nil)
	err = designDB.Cleanup()
	if err != nil {
		t.Error("db cleanup error", err)
	}
}

func TestViewWrapperFunction(t *testing.T) {
	ddoc, err := designDB.Get("_design/test", nil)
	if err != nil {
		t.Error("db get error", err)
	}

	ddoc["views"] = map[string]interface{}{
		"ids":       map[string]string{"map": "function(doc) { emit(doc._id, null); }"},
		"multi_key": map[string]string{"map": "function(doc) { emit(doc.i, null); }"},
	}
	_, err = designDB.Update([]map[string]interface{}{ddoc}, nil)
	if err != nil {
		t.Error("db set error", err)
	}

	results, err := designDB.View("test/multi_key", func(row Row) Row {
		key := row.Key.(float64)
		key *= key
		row.Key = int(key)
		return row
	}, nil)

	if err != nil {
		t.Error("db view error", err)
	}

	rows, err := results.Rows()
	if err != nil {
		t.Error("rows error", err)
	}

	for idx, i := range []int{1, 4, 9, 16, 25} {
		if i != rows[idx].Key.(int) {
			t.Errorf("key = %d want %d", rows[idx].Key.(int), i)
		}
	}
}

func TestUpdateSeq(t *testing.T) {
	err := designDB.Set("foo", map[string]interface{}{})
	if err != nil {
		t.Error("db set error", err)
	}

	results, err := designDB.View("_all_docs", nil, map[string]interface{}{"update_seq": true})
	if err != nil {
		t.Error("db view error", err)
	}

	_, err = results.Rows()
	if err != nil {
		t.Error("rows error", err)
	}

	updateSeq, err := results.UpdateSeq()
	if err != nil {
		t.Errorf("update seq error %s %d", err, updateSeq)
	}

}

func TestProperties(t *testing.T) {
	results, err := designDB.View("_all_docs", nil, nil)
	if err != nil {
		t.Error("db view error", err)
	}

	rows, err := results.Rows()
	if err != nil {
		t.Error("rows error", err)
	}

	if rows == nil {
		t.Error("rows nil")
	}

	totalRows, _ := results.TotalRows()
	if totalRows == -1 {
		t.Error("total rows invalid")
	}

	offset, _ := results.Offset()
	if offset == -1 {
		t.Error("offset invalid")
	}
}

func TestRowRepr(t *testing.T) {
	results, err := designDB.View("_all_docs", nil, nil)
	if err != nil {
		t.Error("db view error", err)
	}

	rows, err := results.Rows()
	if err != nil {
		t.Error("rows error", err)
	}

	if !strings.Contains(rows[0].String(), "id") {
		t.Errorf("row %s not contains id", rows[0])
	}

	if !strings.Contains(rows[0].String(), "Row") {
		t.Errorf("row %s not contains Row", rows[0])
	}

	results, err = designDB.View("test/multi_key", nil, nil)
	if err != nil {
		t.Error("db view error", err)
	}

	rows, err = results.Rows()
	if err != nil {
		t.Error("rows error", err)
	}

	if !strings.Contains(rows[0].String(), "id") {
		t.Errorf("row %s not contains id", rows[0])
	}

	if !strings.Contains(rows[0].String(), "Row") {
		t.Errorf("row %s not contains Row", rows[0])
	}
}

func TestAllRows(t *testing.T) {
	rch, err := iterDB.IterView("test/nums", 10, nil, nil)
	if err != nil {
		t.Fatal("db iter view error", err)
	}

	err = testViewResults(rch, 0, NumDocs, 1)
	if err != nil {
		t.Error("test view results error", err)
	}
}

func testViewResults(rch <-chan Row, begin, end, incr int) error {
	rowsCollected := []Row{}
	for row := range rch {
		rowsCollected = append(rowsCollected, row)
	}

	nums := iterateSlice(begin, end, incr)
	if len(rowsCollected) != min(len(nums), NumDocs) {
		return fmt.Errorf("number of docs %d want %d", len(rowsCollected), len(nums))
	}

	docsLeft := make([]map[string]interface{}, len(nums))
	for idx, row := range rowsCollected {
		docsLeft[idx] = docFromRow(row)
	}

	docsRight := make([]map[string]interface{}, len(nums))
	for idx, num := range nums[:min(len(rowsCollected), NumDocs)] {
		docsRight[idx] = docFromNum(num)
	}

	if !reflect.DeepEqual(docsLeft, docsRight) {
		return errors.New("doc from row not equal to doc from num")
	}
	return nil
}

func iterateSlice(begin, end, incr int) []int {
	s := []int{}
	if begin <= end {
		for i := begin; i < end; i += incr {
			s = append(s, i)
		}
	} else {
		for i := begin; i > end; i += incr {
			s = append(s, i)
		}
	}
	return s
}

func TestBatchSizes(t *testing.T) {
	_, err := iterDB.IterView("test/nums", 0, nil, nil)
	if err != ErrBatchValue {
		t.Fatalf("db iter view %s want %s", err, ErrBatchValue)
	}

	_, err = iterDB.IterView("test/nums", -1, nil, nil)
	if err != ErrBatchValue {
		t.Fatalf("db iter view %s want %s", err, ErrBatchValue)
	}

	rch, err := iterDB.IterView("test/nums", 1, nil, nil)
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResultsLength(rch, NumDocs)
	if err != nil {
		t.Fatal("test view results length error", err)
	}

	rch, err = iterDB.IterView("test/nums", NumDocs/2, nil, nil)
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResultsLength(rch, NumDocs)
	if err != nil {
		t.Fatal("test view results length error", err)
	}

	rch, err = iterDB.IterView("test/nums", NumDocs*2, nil, nil)
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResultsLength(rch, NumDocs)
	if err != nil {
		t.Fatal("test view results length error", err)
	}

	rch, err = iterDB.IterView("test/nums", NumDocs-1, nil, nil)
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResultsLength(rch, NumDocs)
	if err != nil {
		t.Fatal("test view results length error", err)
	}

	rch, err = iterDB.IterView("test/nums", NumDocs, nil, nil)
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResultsLength(rch, NumDocs)
	if err != nil {
		t.Fatal("test view results length error", err)
	}

	rch, err = iterDB.IterView("test/nums", NumDocs+1, nil, nil)
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResultsLength(rch, NumDocs)
	if err != nil {
		t.Fatal("test view results length error", err)
	}
}

func testViewResultsLength(rch <-chan Row, length int) error {
	total := 0
	for _ = range rch {
		total++
	}
	if total != length {
		return fmt.Errorf("length %d want %d", total, length)
	}
	return nil
}

func TestBatchSizesWithSkip(t *testing.T) {
	rch, err := iterDB.IterView("test/nums", NumDocs/10, nil, map[string]interface{}{
		"skip": NumDocs / 2,
	})
	if err != nil {
		t.Fatal("db iter view error", err)
	}

	err = testViewResultsLength(rch, NumDocs/2)
	if err != nil {
		t.Error("test batch sizes with skip error", err)
	}
}

func TestLimit(t *testing.T) {
	var limit int
	var err error
	var rch <-chan Row
	_, err = iterDB.IterView("test/nums", 10, nil, map[string]interface{}{
		"limit": limit,
	})
	if err != ErrLimitValue {
		t.Fatalf("db iter view %s want %s", err, ErrLimitValue)
	}

	for _, limit = range []int{1, NumDocs / 4, NumDocs - 1, NumDocs, NumDocs + 1} {
		rch, err = iterDB.IterView("test/nums", 10, nil, map[string]interface{}{
			"limit": limit,
		})
		if err != nil {
			t.Fatal("db iter view error", err)
		}
		err = testViewResults(rch, 0, limit, 1)
		if err != nil {
			t.Fatal("test view results error", err)
		}
	}

	limit = NumDocs / 4
	rch, err = iterDB.IterView("test/nums", limit, nil, map[string]interface{}{
		"limit": limit,
	})
	if err != nil {
		t.Fatal("db iter view error", err)
	}

	err = testViewResults(rch, 0, limit, 1)
	if err != nil {
		t.Error("test view results error", err)
	}
}

func TestDescending(t *testing.T) {
	rch, err := iterDB.IterView("test/nums", 10, nil, map[string]interface{}{"descending": true})
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResults(rch, NumDocs-1, -1, -1)
	if err != nil {
		t.Error("test view results error", err)
	}

	rch, err = iterDB.IterView("test/nums", 10, nil, map[string]interface{}{
		"descending": true,
		"limit":      NumDocs / 4,
	})
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResults(rch, NumDocs-1, NumDocs*3/4-1, -1)
	if err != nil {
		t.Error("test view results error", err)
	}
}

func TestStartKey(t *testing.T) {
	vch, err := iterDB.IterView("test/nums", 10, nil, map[string]interface{}{"startkey": NumDocs/2 - 1})
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResults(vch, NumDocs-2, NumDocs, 1)
	if err != nil {
		t.Fatal("test view results error", err)
	}

	vch, err = iterDB.IterView("test/nums", 10, nil, map[string]interface{}{"startkey": 1, "descending": true})
	if err != nil {
		t.Error("db iter view error", err)
	}
	err = testViewResults(vch, 3, -1, -1)
	if err != nil {
		t.Error("teset view results error", err)
	}
}

func TestNullKeys(t *testing.T) {
	vch, err := iterDB.IterView("test/nulls", 10, nil, nil)
	if err != nil {
		t.Fatal("db iter view error", err)
	}
	err = testViewResultsLength(vch, NumDocs)
	if err != nil {
		t.Error("test view results length error", err)
	}
}

func TestViewDefinitionOptions(t *testing.T) {
	options := map[string]interface{}{"collation": "raw"}
	view, err := NewViewDefinition("bar", "baz", "function(doc) { emit(doc._id, doc._rev); }", "", "", nil, options)
	if err != nil {
		t.Fatal("create view definition error", err)
	}
	_, err = view.Sync(defnDB)
	if err != nil {
		t.Fatal("view definition sync error", err)
	}

	designDoc, err := defnDB.Get("_design/bar", nil)
	if err != nil {
		t.Fatal("db get error", err)
	}

	designOptions := designDoc["views"].(map[string]interface{})["baz"].(map[string]interface{})["options"]

	if !reflect.DeepEqual(options, designOptions) {
		t.Error("options not identical")
	}
}

func TestRetrieveViewDefinition(t *testing.T) {
	version, _ := server.Version()

	view, err := NewViewDefinition("foo", "bar", "baz", "", "", nil, nil)
	if err != nil {
		t.Fatal("create view definition error", err)
	}
	results, err := view.Sync(defnDB)
	if err != nil {
		t.Fatal("view definition sync error", err)
	}

	if len(results) <= 0 {
		t.Fatal("no results returned")
	}

	if strings.HasPrefix(version, "1") {
		if results[0].Err != nil {
			t.Error("update result error", results[0].Err)
		}
	} else if strings.HasPrefix(version, "2") {
		if results[0].Err != ErrInternalServerError {
			t.Errorf("update result error %v want %v", results[0].Err, ErrInternalServerError)
		}
	}

	if results[0].ID != "_design/foo" {
		t.Fatalf("results ID %s want _design/foo", results[0].ID)
	}

	_, err = defnDB.Get(results[0].ID, nil)
	if strings.HasPrefix(version, "1") {
		if err != nil {
			t.Fatal("db get error", err)
		}
	} else if strings.HasPrefix(version, "2") {
		if err != ErrNotFound {
			t.Fatalf("db get error %s want %s", err, ErrNotFound)
		}
	}

}

func TestSyncMany(t *testing.T) {
	jsFunc := "function(doc) { emit(doc._id, doc._rev); }"
	firstView, _ := NewViewDefinition("design_doc", "view_one", jsFunc, "", "", nil, nil)
	secondView, _ := NewViewDefinition("design_doc_two", "view_one", jsFunc, "", "", nil, nil)
	thirdView, _ := NewViewDefinition("design_doc", "view_two", jsFunc, "", "", nil, nil)
	results, err := SyncMany(defnDB, []*ViewDefinition{firstView, secondView, thirdView}, false, nil)
	if err != nil {
		t.Fatal("sync many error", err)
	}
	valid := 0
	for _, res := range results {
		if res.Err == nil {
			valid++
		}
	}
	if valid != 2 {
		t.Errorf("returned %d results, there should be only 2 design documents", len(results))
	}
}

func TestShowUrls(t *testing.T) {
	_, data, err := showListDB.Show("_design/foo/_show/bar", "", nil)
	if err != nil {
		t.Fatal("db show error", err)
	}

	if string(data) != "null:<default>" {
		t.Errorf("db show returns %s want null:<default>", string(data))
	}

	_, data, err = showListDB.Show("foo/bar", "", nil)
	if err != nil {
		t.Fatal("db show error", err)
	}

	if string(data) != "null:<default>" {
		t.Errorf("db show returns %s want null:<default>", string(data))
	}
}

func TestShowDocID(t *testing.T) {
	_, data, err := showListDB.Show("foo/bar", "", nil)
	if err != nil {
		t.Fatal("db show error", err)
	}

	if string(data) != "null:<default>" {
		t.Errorf("db show returns %s want null:<default>", string(data))
	}

	_, data, err = showListDB.Show("foo/bar", "1", nil)
	if err != nil {
		t.Fatal("db show error", err)
	}

	if string(data) != "1:<default>" {
		t.Errorf("db show returns %s want 1:<default>", string(data))
	}

	_, data, err = showListDB.Show("foo/bar", "2", nil)
	if err != nil {
		t.Fatal("db show error", err)
	}

	if string(data) != "2:<default>" {
		t.Errorf("db show returns %s want 2:<default>", string(data))
	}
}

func TestShowParams(t *testing.T) {
	_, data, err := showListDB.Show("foo/bar", "", url.Values{"r": []string{"abc"}})
	if err != nil {
		t.Fatal("db show error", err)
	}

	if string(data) != "null:abc" {
		t.Errorf("db show returns %s want null:abc", string(data))
	}
}

func TestList(t *testing.T) {
	_, data, err := showListDB.List("foo/list", "foo/by_id", nil)
	if err != nil {
		t.Fatal("db list error", err)
	}

	if string(data) != `1\r\n2\r\n` {
		t.Errorf("db list returns %s want `1\r\n2\r\n`", string(data))
	}

	_, data, err = showListDB.List("foo/list", "foo/by_id", map[string]interface{}{"include_header": true})
	if err != nil {
		t.Fatal("db list error", err)
	}

	if string(data) != `id\r\n1\r\n2\r\n` {
		t.Errorf("db list returns %s want `id\r\n1\r\n2\r\n`", string(data))
	}
}

func TestListKeys(t *testing.T) {
	_, data, err := showListDB.List("foo/list", "foo/by_id", map[string]interface{}{"keys": []string{"1"}})
	if err != nil {
		t.Fatal("db list error", err)
	}

	if string(data) != `1\r\n` {
		t.Errorf("db list returns %s want `1\r\n`", string(data))
	}
}

func TestListViewParams(t *testing.T) {
	_, data, err := showListDB.List("foo/list", "foo/by_name", map[string]interface{}{"startkey": "o", "endkey": "p"})
	if err != nil {
		t.Fatal("db list error", err)
	}

	if string(data) != `1\r\n` {
		t.Errorf("db list returns %s want `1\r\n`", string(data))
	}

	_, data, err = showListDB.List("foo/list", "foo/by_name", map[string]interface{}{"descending": true})
	if err != nil {
		t.Fatal("db list error", err)
	}

	if string(data) != `2\r\n1\r\n` {
		t.Errorf("db list returns %s want `2\r\n1\r\n`", string(data))
	}
}

func TestEmptyDoc(t *testing.T) {
	_, data, err := updateDB.UpdateDoc("foo/bar", "", nil)
	if err != nil {
		t.Fatal("db updatedoc error", err)
	}

	if string(data) != "empty doc" {
		t.Errorf("db list returns %s want empty doc", string(data))
	}
}

func TestNewDoc(t *testing.T) {
	_, data, err := updateDB.UpdateDoc("foo/bar", "new", nil)
	if err != nil {
		t.Fatal("db updatedoc error", err)
	}

	if string(data) != "new doc" {
		t.Errorf("db list returns %s want new doc", string(data))
	}
}

func TestUpdateDoc(t *testing.T) {
	_, data, err := updateDB.UpdateDoc("foo/bar", "existed", nil)
	if err != nil {
		t.Fatal("db updatedoc error", err)
	}

	if string(data) != "hello doc" {
		t.Errorf("db list returns %s want hello doc", string(data))
	}
}
