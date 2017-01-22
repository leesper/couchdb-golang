package couchdb

type Test struct {
	D MapField
}

// func TestMutableFields(t *testing.T) {
// 	a := Test{}
// 	b := Test{}
// 	a.D["x"] = true
// 	if !a.D["x"].(bool) {
// 		t.Error("a.D[x] false")
// 	}
// 	if b.D["x"].(bool) {
// 		t.Error("b.D[x] true")
// 	}
// }
//
// type Post struct {
// 	Title TextField
// }
//
// func TestAutomaticID(t *testing.T) {
// 	post := Post{title: "Foo bar"}
// 	if post.ID != "" {
// 		t.Error("post ID not empty", post.ID)
// 	}
// 	err := post.Store(mappingDB)
// 	if err != nil {
// 		t.Fatal("document store error", err)
// 	}
// 	if post.ID == "" {
// 		t.Error("post ID empty")
// 	}
// 	doc, err := mappingDB.Get(post.ID, nil)
// 	if err != nil {
// 		t.Fatal("db get error", err)
// 	}
// 	if doc["title"].(string) != "Foo bar" {
// 		t.Errorf("doc title %s want Foo bar", doc["title"].(string))
// 	}
// }
//
// func TestExplicitIDViaInit(t *testing.T) {
// 	post := Post{ID: "foo_bar", Title: "Foo bar"}
// 	if post.ID != "foo_bar" {
// 		t.Fatalf("post ID %s want foo_bar", post.ID)
// 	}
// 	err := post.Store(mappingDB)
// 	if err != nil {
// 		t.Fatal("document store error", err)
// 	}
// 	doc, err := mappingDB.Get(post.ID, nil)
// 	if err != nil {
// 		t.Fatal("db get error", err)
// 	}
// 	if doc["title"].(string) != "Foo bar" {
// 		t.Errorf("doc title %s want Foo bar", doc["title"].(string))
// 	}
// }
//
// func TestExplicitIDViaSetter(t *testing.T) {
// 	post := Post{Title: "Foo bar"}
// 	post.ID = "foo_bar"
// 	if post.ID != "foo_bar" {
// 		t.Errorf("post ID %s want foo_bar", post.ID)
// 	}
// 	err := post.Store(mappingDB)
// 	if err != nil {
// 		t.Fatal("document store error", err)
// 	}
// 	doc, err := mappingDB.Get(post.ID, nil)
// 	if err != nil {
// 		t.Fatal("db get error", err)
// 	}
// 	if doc["title"].(string) != "Foo bar" {
// 		t.Errorf("doc title %s want Foo bar", doc["title"].(string))
// 	}
// }
//
// func TestChangeIDFailure(t *testing.T) {
// 	post := Post{Title: "Foo bar"}
// 	err := post.Store(mappingDB)
// 	if err != nil {
// 		t.Fatal("document store error", err)
// 	}
// 	post, err = post.Load(mappingDB, post.ID)
// 	if err != nil {
// 		t.Fatal("document load error", err)
// 	}
// 	err = post.SetID("foo_bar")
// 	if err != ErrSetID {
// 		t.Error("document set id error %v want %v", err, ErrSetID)
// 	}
// }
//
// func TestBatchUpdate(t *testing.T) {
// 	post1 := Post{Title: "Foo bar"}
// 	post2 := Post{Title: "Foo baz"}
// 	results, err := mappingDB.Update([]map[string]interface{}{post1.ToMap(), post2.ToMap()}, nil)
// 	if len(results) != 2 {
// 		t.Fatalf("len(results) = %d want 2", len(results))
// 	}
// 	for idx, res := range results {
// 		if res.Err != nil {
// 			t.Errorf("result %d error %v", idx, res.Err)
// 		}
// 	}
// }
//
// func TestStoreExisting(t *testing.T) {
// 	post := Post{Title: "Foo bar"}
// 	err := post.Store(mappingDB)
// 	if err != nil {
// 		t.Fatal("document store error", err)
// 	}
// 	err = post.Store(mappingDB)
// 	if err != nil {
// 		t.Fatal("document store error", err)
// 	}
// 	results, err := mappingDB.View("_all_docs", nil, nil)
// 	if err != nil {
// 		t.Fatal("db view _all_docs error", err)
// 	}
// 	rows, err := results.Rows()
// 	if err != nil {
// 		t.Fatal("rows error", err)
// 	}
// 	if len(rows) != 1 {
// 		t.Errorf("len(rows) %d want 1", len(rows))
// 	}
// }
//
// func TestOldDateTime(t *testing.T) {}
//
// func TestDateTimeWithMicroseconds(t *testing.T) {}
//
// func TestDateTimeToJSON(t *testing.T) {
// 	dt := DateTimeField{}
// 	d := time.Now()
// 	if dt.ToJSON(d) != d {
// 		t.Error("date time not equal")
// 	}
// }
//
// func TestGetHasDefault(t *testing.T) {
// 	doc := Document{}
// 	doc.get("foo")
// 	doc.getDefault("foo", nil)
// }
//
// type PostWithComment struct {
// 	Title    TextField
// 	Comments ListField
// }
//
// func TestListFieldToJSON(t *testing.T) {
// 	post := PostWithComment{Title: "Foo bar"}
// 	comment := MapField{"author": "myself", "content": "Bla bla"}
// 	post.Comments = append(post.Comments, comment)
// 	if !reflect.DeepEqual(post.Comments, []MapField{comment}) {
// 		t.Errorf("post comment %v not equal to %v", post.Comments, []MapField{comment})
// 	}
// }
//
// type Thing struct {
// 	Numbers ListField
// }
//
// func TestListFieldProxyAppend(t *testing.T) {
// 	thing := Thing{Numbers: []DecimalField{DecimalField(1.0), DecimalField(2.0)}}
// 	thing.Numbers = append(thing.Numbers, DecimalField(3.0))
// 	if len(thing.Numbers != 3) {
// 		t.Fatalf("thing numbers length %d want 3", len(thing.Numbers))
// 	}
// 	if !reflect.DeepEqual(DecimalField(3.0), thing.Numbers[2]) {
// 		t.Errorf("thing numbers[2] = %v want %v", thing.Numbers[2], DecimalField(3.0))
// 	}
// }
//
// func TestListFieldProxyContains(t *testing.T) {
// 	thing := Thing{Numbers: []DecimalField{DecimalField(1.0), DecimalField(2.0)}}
// 	found := false
// 	count := 0
// 	index := -1
// 	for idx, num := range thing.Numbers {
// 		if num == DecimalField(1.0) {
// 			found = true
// 			count++
// 			index = idx
// 		}
// 	}
// 	if !found {
// 		t.Error("Decimal 1.0 not in thing")
// 	}
// 	if count != 1 {
// 		t.Errorf("Decimal 1.0 count %d want 1", count)
// 	}
// 	if index != 0 {
// 		t.Errorf("Decimal 1.0 index %d want 0", index)
// 	}
// }
//
// func TestListFieldProxyInsert(t *testing.T) {
// 	thing := Thing{Numbers: []DecimalField{DecimalField(1.0), DecimalField(2.0)}}
// 	thing.Numbers = append([]DecimalField{DecimalField(0.0)}, thing.Numbers...)
// 	if len(thing.Numbers) != 3 {
// 		t.Errorf("thing numbers length %d want 3", len(thing.Numbers))
// 	}
// 	if thing.Numbers[0] != DecimalField(0.0) {
// 		t.Errorf("thing numbers[0] %v want %v", thing.Numbers[0], DecimalField(0.0))
// 	}
// }
//
// func TestListFieldProxyIter(t *testing.T) {
// 	err := mappingDB.Set("test", map[string]interface{}{"numbers": []float64{1.0, 2.0}})
// 	if err != nil {
// 		t.Fatal("db set error", err)
// 	}
// 	thing, err := Thing.Load(mappingDB, "test")
// 	if err != nil {
// 		t.Fatal("document load error", err)
// 	}
// 	if thing.Numbers[0] != DecimalField(1.0) {
// 		t.Errorf("thing numbers[0] %v want %v", thing.Numbers[0], DecimalField(1.0))
// 	}
// }
//
// func TestListFieldIterDict(t *testing.T) {
// 	comments := []map[string]interface{}{
// 		{"author": "Joe", "content": "Hey"},
// 	}
// 	err := mappingDB.Set("test", comments)
// 	if err != nil {
// 		t.Fatal("db set error", err)
// 	}
// 	post, err := Post.Load(mappingDB, "test")
// 	if err != nil {
// 		t.Fatal("document load error", err)
// 	}
// 	if !reflect.DeepEqual(post.Comments, comments) {
// 		t.Error("post comments %v want %v", post.Comments, comments)
// 	}
// }
//
// func TestListFieldProxyPop(t *testing.T) {
// 	thing := Thing{}
// 	thing.Numbers = []DecimalField{DecimalField(0.0), DecimalField(1.0), DecimalField(2.0)}
// 	last := thing.Numbers[len(thing.Numbers)-1]
// 	if last != DecimalField(2.0) {
// 		t.Errorf("last %v want %v", last, DecimalField(2.0))
// 	}
// 	thing.Numbers = thing.Numbers[:len(thing.Numbers)-1]
// 	if len(thing.Numbers) != 2 {
// 		t.Errorf("thing numbers length %d want 2", len(thing.Numbers))
// 	}
// 	first := thing.Numbers[0]
// 	if first != DecimalField(0.0) {
// 		t.Errorf("first %v want %v", first, DecimalField(0.0))
// 	}
// }
//
// func TestListFieldProxySlices(t *testing.T) {
// 	thing := Thing{}
// 	thing.Numbers = []DecimalField{DecimalField(0.0), DecimalField(1.0), DecimalField(2.0), DecimalField(3.0), DecimalField(4.0)}
// 	ll := thing.Numbers[1:3]
// 	if len(ll) != 2 {
// 		t.Errorf("ll length %d want 2", len(ll))
// 	}
// 	if ll[0] != DecimalField(1.0) {
// 		t.Errorf("ll[0] = %v want %v", ll[0], DecimalField(1.0))
// 	}
// 	thing.Numbers[2:4] = []DecimalField{DecimalField(6.0), DecimalField(7.0)}
// 	if thing.Numbers[2] != DecimalField(6.0) {
// 		t.Errorf("thing numbers[2] = %v want %v", thing.Numbers[2], DecimalField(6.0))
// 	}
// 	if thing.Numbers[4] != DecimalField(4.0) {
// 		t.Errorf("thing numbers[4] = %v want %v", thing.Numbers[4], DecimalField(4.0))
// 	}
// 	if len(thing.Numbers) != 5 {
// 		t.Errorf("ll length %d want 5", len(thing.Numbers))
// 	}
// 	thing.Numbers = thing.Numbers[3:]
// 	if len(thing.Numbers) != 3 {
// 		t.Errorf("ll length %d want 3", len(thing.Numbers))
// 	}
// }
//
// func TestListFieldMutableFields(t *testing.T) {
// 	thing := Thing.Wrap(map[string]interface{}{"_id": "foo", "_rev": 1})
// 	thing.Numbers = append(thing.Numbers, DecimalField(1.0))
// 	thing2 := Thing{"_id": "thing2"}
// 	if len(thing2.Numbers) != 0 {
// 		t.Errorf("ll length %d want 0", len(thing.Numbers))
// 	}
// }
//
// func TestViewFieldProperty(t *testing.T) {
// 	item := unitTestItem{}
// 	err := item.Store(mappingDB)
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
//
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
//
// func TestWrappedView(t *testing.T) {
// 	item := unitTestItem{}
// 	results, err := mappingDB.View("_all_docs", item.wrapRow, nil)
// 	doc = results.Rows()[0]
// 	mappingDB.Delete(doc["_id"].(string))
// }
//
// func TestQuery(t *testing.T) {
// 	item := unitTestItem{}
// 	results, err := item.Query(mappingDB, allMapFunc, nil)
// 	results, err := item.Query(mappingDB, allMapFunc, true)
// }
