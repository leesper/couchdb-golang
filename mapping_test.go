package couchdb

import "testing"

type Post struct {
	Title string `json:"title"`
	Document
}

func TestAutomaticID(t *testing.T) {
	post := Post{Title: "Foo bar"}
	if post.GetID() != "" {
		t.Error("post ID not empty", post.ID)
	}

	err := Store(mappingDB, &post)
	if err != nil {
		t.Fatal("document store error", err)
	}

	if post.GetID() == "" {
		t.Error("post ID empty")
	}

	if post.GetRev() == "" {
		t.Error("post rev empty")
	}

	doc, err := mappingDB.Get(post.GetID(), nil)
	if err != nil {
		t.Fatal("db get error", err)
	}

	if doc["title"].(string) != "Foo bar" {
		t.Errorf("doc title %s want Foo bar", doc["title"].(string))
	}

	if doc["_rev"].(string) != post.GetRev() {
		t.Errorf("post rev %s want %s", post.GetRev(), doc["_rev"].(string))
	}
}

func TestExplicitIDByInit(t *testing.T) {
	post := Post{Document: DocumentWithID("foo_bar"), Title: "Foo bar"}
	if post.GetID() != "foo_bar" {
		t.Fatalf("post ID %s want foo_bar", post.GetID())
	}

	err := Store(mappingDB, &post)
	if err != nil {
		t.Fatal("document store error", err)
	}

	doc, err := mappingDB.Get(post.GetID(), nil)
	if err != nil {
		t.Fatal("db get error", err)
	}

	if doc["title"].(string) != "Foo bar" {
		t.Errorf("doc title %s want Foo bar", doc["title"].(string))
	}

	if doc["_id"].(string) != post.GetID() {
		t.Errorf("post id %s want %s", post.GetID(), doc["_id"].(string))
	}

	if doc["_rev"].(string) != post.GetRev() {
		t.Errorf("post rev %s want %s", post.GetRev(), doc["_rev"].(string))
	}
}

func TestExplicitIDBySetter(t *testing.T) {
	post := Post{Title: "Foo bar"}
	post.SetID("foo_baz")

	if post.GetID() != "foo_baz" {
		t.Errorf("post ID %s want foo_bar", post.GetID())
	}

	err := Store(mappingDB, &post)
	if err != nil {
		t.Fatal("document store error", err)
	}

	doc, err := mappingDB.Get(post.GetID(), nil)
	if err != nil {
		t.Fatal("db get error", err)
	}

	if doc["title"].(string) != "Foo bar" {
		t.Errorf("doc title %s want Foo bar", doc["title"].(string))
	}

	if doc["_id"].(string) != post.GetID() {
		t.Errorf("post id %s want %s", post.GetID(), doc["_id"].(string))
	}

	if doc["_rev"].(string) != post.GetRev() {
		t.Errorf("post rev %s want %s", post.GetRev(), doc["_rev"].(string))
	}
}

func TestChangeIDFailure(t *testing.T) {
	post := Post{Title: "Foo bar"}

	err := Store(mappingDB, &post)
	if err != nil {
		t.Fatal("document store error", err)
	}

	err = Load(mappingDB, post.GetID(), &post)
	if err != nil {
		t.Fatal("document load error", err)
	}

	err = post.SetID("foo_bar")
	if err != ErrSetID {
		t.Errorf("document set id error %v want %v", err, ErrSetID)
	}
}

type NotDocument struct{}

func TestNotADocument(t *testing.T) {
	notDocument := NotDocument{}
	err := Store(mappingDB, &notDocument)
	if err != ErrNotDocumentEmbedded {
		t.Fatalf("store error %v want %v", err, ErrNotDocumentEmbedded)
	}
}

func TestNestedStruct(t *testing.T) {}

////////////////////////////////////////////////////////////////////////////////

// func TestBatchUpdate(t *testing.T) {
// 	post1 := Post{Title: "Foo bar"}
// 	post2 := Post{Title: "Foo baz"}
// 	results, err := mappingDB.Update([]map[string]interface{}{ToMap(post1), ToMap(post2)}, nil)
// 	if err != nil {
// 		t.Fatal("db update error", err)
// 	}
//
// 	if len(results) != 2 {
// 		t.Fatalf("len(results) = %d want 2", len(results))
// 	}
//
// 	for idx, res := range results {
// 		if res.Err != nil {
// 			t.Errorf("result %d error %v", idx, res.Err)
// 		}
// 	}
// }
//
// func TestStoreExisting(t *testing.T) {
// 	post := Post{Title: "Foo bar"}
// 	err := Store(mappingDB, post)
// 	if err != nil {
// 		t.Fatal("document store error", err)
// 	}
//
// 	err = Store(mappingDB, post)
// 	if err != nil {
// 		t.Fatal("document store error", err)
// 	}
//
// 	results, err := mappingDB.View("_all_docs", nil, nil)
// 	if err != nil {
// 		t.Fatal("db view _all_docs error", err)
// 	}
//
// 	rows, err := results.Rows()
// 	if err != nil {
// 		t.Fatal("rows error", err)
// 	}
//
// 	if len(rows) != 1 {
// 		t.Errorf("len(rows) %d want 1", len(rows))
// 	}
// }
//
// type PostWithComment struct {
// 	Title    string
// 	Comments []map[string]string
// 	Document
// }
//
// func TestCompareDocWithObj(t *testing.T) {
// 	postDoc := map[string]interface{}{
// 		"comments": []map[string]string{
// 			{"author": "Joe", "content": "Hey"},
// 		},
// 	}
// 	err := mappingDB.Set("test", postDoc)
// 	if err != nil {
// 		t.Fatal("db set error", err)
// 	}
// 	post1 := PostWithComment{}
// 	err = Load(mappingDB, "test", &post1)
// 	if err != nil {
// 		t.Fatal("document load error", err)
// 	}
// 	post2 := PostWithComment{}
// 	err = FromMap(&post2, postDoc)
// 	if err != nil {
// 		t.Fatal("from map error", err)
// 	}
// 	if !reflect.DeepEqual(post1, post2) {
// 		t.Errorf("post1 %v != post2 %v", post1, post2)
// 	}
// }
//
// type Thing struct {
// 	Numbers []float64
// 	Document
// }
//
// func compareFloat(x, y float64) int {
// 	a := big.NewFloat(x)
// 	b := big.NewFloat(y)
// 	return a.Cmp(b)
// }
//
// func TestSliceFieldFloat(t *testing.T) {
// 	err := mappingDB.Set("test", map[string]interface{}{"numbers": []float64{1.0, 2.0}})
// 	if err != nil {
// 		t.Fatal("db set error", err)
// 	}
//
// 	thing := Thing{}
// 	err = Load(mappingDB, "test", &thing)
// 	if err != nil {
// 		t.Fatal("document load error", err)
// 	}
//
// 	if compareFloat(thing.Numbers[0], 1.0) != 0 {
// 		t.Errorf("thing numbers[0] %v want 1.0", thing.Numbers[0])
// 	}
// }
//
// func TestAutomaticLowerCase(t *testing.T) {}

type unitTestItem struct {
	withIncludeDocs    ViewField
	withoutIncludeDocs ViewField
	Document
}

var allMapFunc = `function(doc) { emit(doc._id, doc); }`

// func TestViewFieldProperty(t *testing.T) {
// 	item := unitTestItem{}
// 	err := Store(mappingDB, item)
// 	if err != nil {
// 		t.Fatal("document store error", err)
// 	}
// 	results, err := item.withIncludeDocs(mappingDB)
// 	if err != nil {
// 		t.Fatal("view with include docs error", err)
// 	}
// 	rows, err := results.Rows()
// 	if err != nil {
// 		t.Fatal("rows error", err)
// 	}
// 	fmt.Println(rows[0])
// }

// func TestView(t *testing.T) {
// 	item := unitTestItem{}
// 	results, err := item.View(mappingDB, "test/without_include_docs")
// 	if err != nil {
// 		t.Fatal("view without include docs error", err)
// 	}
// 	rows, err := results.Rows()
// 	if err != nil {
// 		t.Fatal("rows error", err)
// 	}
// 	fmt.Println(rows[0])
//
// 	result, err = item.View(mappingDB, "test/with_include_docs")
// 	if err != nil {
// 		t.Fatal("view with include docs error", err)
// 	}
// 	rows, err := results.Rows()
// 	if err != nil {
// 		t.Fatal("rows error", err)
// 	}
// 	fmt.Println(rows[0])
// }

// func TestWrappedView(t *testing.T) {
// 	item := unitTestItem{}
// 	results, err := mappingDB.View("_all_docs", item.wrapRow, nil)
// 	doc = results.Rows()[0]
// 	mappingDB.Delete(doc["_id"].(string))
// }

// func TestQuery(t *testing.T) {
// 	item := unitTestItem{}
// 	results, err := item.Query(mappingDB, allMapFunc, nil)
// 	results, err := item.Query(mappingDB, allMapFunc, true)
// }
