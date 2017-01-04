package couchdb

import (
	"encoding/json"
	"reflect"
	"testing"
)

// func TestNewDefaultDB(t *testing.T) {
// 	dbDefault, err := NewDatabase("golang-default")
// 	if err != nil {
// 		t.Errorf("new default database error %v", err)
// 	}
// 	if dbDefault.Available() {
// 		t.Error(`db available`)
// 	}
// }
//
// func TestNewDB(t *testing.T) {
// 	newDB := "golang-newdb"
// 	s.Create(newDB)
// 	defer s.Delete(newDB)
// 	dbNew, err := NewDatabase("http://root:likejun@localhost:5984/" + newDB)
// 	if err != nil {
// 		t.Error(`new database error`, err)
// 	}
// 	if !dbNew.Available() {
// 		t.Error(`db not available`)
// 	}
// }
//
// func TestSaveNew(t *testing.T) {
// 	doc := map[string]interface{}{"doc": "bar"}
// 	id, rev, err := db.Save(doc, nil)
// 	if err != nil {
// 		t.Error(`db save error`, err)
// 	}
// 	if id != doc["_id"].(string) {
// 		t.Errorf("invalid id: %q", id)
// 	}
// 	if rev != doc["_rev"].(string) {
// 		t.Errorf("invalid rev: %q", rev)
// 	}
// }
//
// func TestSaveNewWithID(t *testing.T) {
// 	doc := map[string]interface{}{"_id": "foo"}
// 	id, rev, err := db.Save(doc, nil)
// 	if err != nil {
// 		t.Error(`db save error`, err)
// 	}
// 	if doc["_id"].(string) != "foo" {
// 		t.Errorf("doc[_id] = %s, not foo", doc["_id"])
// 	}
// 	if id != "foo" {
// 		t.Errorf("id = %s, not foo", id)
// 	}
// 	if rev != doc["_rev"].(string) {
// 		t.Errorf("invalid rev: %q", rev)
// 	}
// }
//
// func TestSaveExisting(t *testing.T) {
// 	doc := map[string]interface{}{}
// 	idOld, revOld, err := db.Save(doc, nil)
// 	if err != nil {
// 		t.Error(`db save error`, err)
// 	}
// 	doc["foo"] = true
// 	idNew, revNew, err := db.Save(doc, nil)
// 	if err != nil {
// 		t.Error(`db save foo error`, err)
// 	}
// 	if idOld != idNew {
// 		t.Errorf("ids are not equal old %s new %s", idOld, idNew)
// 	}
// 	if doc["_rev"].(string) != revNew {
// 		t.Errorf("invalid rev %s want %s", doc["_rev"].(string), revNew)
// 	}
// 	if revOld == revNew {
// 		t.Errorf("new rev is equal to old rev %s", revOld)
// 	}
// }
//
// func TestSaveNewBatch(t *testing.T) {
// 	doc := map[string]interface{}{"_id": "foo"}
// 	_, rev, err := db.Save(doc, url.Values{"batch": []string{"ok"}})
// 	if err != nil {
// 		t.Error(`db save batch error`, err)
// 	}
// 	if len(rev) > 0 {
// 		t.Error(`rev not empty`, rev)
// 	}
// 	if r, ok := doc["_rev"]; ok {
// 		t.Error(`doc has _rev field`, r.(string))
// 	}
// }
//
// func TestSaveExistingBatch(t *testing.T) {
// 	doc := map[string]interface{}{"_id": "bar"}
// 	idOld, revOld, err := db.Save(doc, nil)
// 	if err != nil {
// 		t.Error(`db save error`, err)
// 	}
//
// 	idNew, revNew, err := db.Save(doc, url.Values{"batch": []string{"ok"}})
// 	if err != nil {
// 		t.Error(`db save batch error`, err)
// 	}
//
// 	if idOld != idNew {
// 		t.Errorf("old id %s not equal to new id %s", idOld, idNew)
// 	}
//
// 	if len(revNew) > 0 {
// 		t.Error(`rev not empty`, revNew)
// 	}
//
// 	if doc["_rev"].(string) != revOld {
// 		t.Errorf("doc[_rev] %s not equal to old rev %s", doc["_rev"].(string), revOld)
// 	}
// }
//
// func TestDatabaseExists(t *testing.T) {
// 	if !db.Available() {
// 		t.Error(`golang-tests not available`)
// 	}
// 	dbMissing, _ := NewDatabase("golang-missing")
// 	if dbMissing.Available() {
// 		t.Error(`golang-missing available`)
// 	}
// }
//
// func TestDatabaseName(t *testing.T) {
// 	name, err := db.Name()
// 	if err != nil {
// 		t.Error(`db name error`, err)
// 	}
// 	if name != "golang-tests" {
// 		t.Error("db name %s, want golang-tests", name)
// 	}
// }
//
// func TestDatabaseString(t *testing.T) {
// 	if db.String() != "Database http://root:likejun@localhost:5984/golang-tests" {
// 		t.Error(`db string invalid`, db)
// 	}
// }
//
// func TestCommit(t *testing.T) {
// 	if err := db.Commit(); err != nil {
// 		t.Error(`db commit error`, err)
// 	}
// }
//
// func TestCreateLargeDoc(t *testing.T) {
// 	var buf bytes.Buffer
// 	// 10MB
// 	for i := 0; i < 110*1024; i++ {
// 		buf.WriteString("0123456789")
// 	}
// 	doc := map[string]interface{}{"data": buf.String()}
// 	if err := db.Set("large", doc); err != nil {
// 		t.Error(`db set error`, err)
// 	}
// 	doc, err := db.Get("large", nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
// 	if doc["_id"].(string) != "large" {
// 		t.Errorf("doc[_id] = %s, want large", doc["_id"].(string))
// 	}
// 	err = db.DeleteDoc(doc)
// 	if err != nil {
// 		t.Error(`db delete doc error`, err)
// 	}
// }
//
// func TestDocIDQuoting(t *testing.T) {
// 	doc := map[string]interface{}{"foo": "bar"}
// 	err := db.Set("foo/bar", doc)
// 	if err != nil {
// 		t.Error(`db set error`, err)
// 	}
// 	doc, err = db.Get("foo/bar", nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
// 	if doc["foo"].(string) != "bar" {
// 		t.Errorf("doc[foo] = %s want bar", doc["foo"].(string))
// 	}
// 	err = db.Delete("foo/bar")
// 	if err != nil {
// 		t.Error(`db delete error`, err)
// 	}
// 	_, err = db.Get("foo/bar", nil)
// 	if err == nil {
// 		t.Error(`db get foo/bar ok`)
// 	}
// }
//
// func TestDisallowNaN(t *testing.T) {
// 	doc := map[string]interface{}{"number": math.NaN()}
// 	err := db.Set("foo", doc)
// 	if err == nil {
// 		t.Error(`db set NaN ok`)
// 	}
// }
//
// func TestDisallowNilID(t *testing.T) {
// 	err := db.DeleteDoc(map[string]interface{}{"_id": nil, "_rev": nil})
// 	if err == nil {
// 		t.Error(`db delete doc with id nil ok`)
// 	}
// 	err = db.DeleteDoc(map[string]interface{}{"_id": "foo", "_rev": nil})
// 	if err == nil {
// 		t.Error(`db delete doc with rev nil ok`)
// 	}
// }
//
// func TestDocRevs(t *testing.T) {
// 	uuid := GenerateUUID()
// 	doc := map[string]interface{}{"bar": 42}
// 	err := db.Set(uuid, doc)
// 	if err != nil {
// 		t.Error(`db set doc error`, err)
// 	}
// 	oldRev := doc["_rev"].(string)
// 	doc["bar"] = 43
// 	err = db.Set(uuid, doc)
// 	if err != nil {
// 		t.Error(`db set doc error`, err)
// 	}
// 	newRev := doc["_rev"].(string)
//
// 	newDoc, err := db.Get(uuid, nil)
// 	if newRev != newDoc["_rev"].(string) {
// 		t.Errorf("new doc rev %s not equal to %s", newDoc["_rev"].(string), newRev)
// 	}
// 	newDoc, err = db.Get(uuid, url.Values{"rev": []string{newRev}})
// 	if newRev != newDoc["_rev"].(string) {
// 		t.Errorf("new doc rev %s not equal to %s", newDoc["_rev"].(string), newRev)
// 	}
// 	oldDoc, err := db.Get(uuid, url.Values{"rev": []string{oldRev}})
// 	if oldRev != oldDoc["_rev"].(string) {
// 		t.Errorf("old doc rev %s not equal to %s", oldDoc["_rev"].(string), oldRev)
// 	}
//
// 	revs, err := db.Revisions(uuid, nil)
// 	if err != nil {
// 		t.Error(`db revisions error`, err)
// 	}
// 	if revs[0]["_rev"].(string) != newRev {
// 		t.Errorf("revs first %s not equal to %s", revs[0]["_rev"].(string), newRev)
// 	}
// 	if revs[1]["_rev"].(string) != oldRev {
// 		t.Errorf("revs second %s not equal to %s", revs[1]["_rev"].(string), oldRev)
// 	}
// 	_, err = db.Revisions("crap", nil)
// 	if err == nil {
// 		t.Error(`db revisions crap ok`)
// 	}
//
// 	err = db.Compact()
// 	if err != nil {
// 		t.Error("db compact error", err)
// 	}
//
// 	info, err := db.Info()
// 	if err != nil {
// 		t.Error(`db info error`, err)
// 	}
// 	for info["compact_running"].(bool) {
// 		info, err = db.Info()
// 		if err != nil {
// 			t.Error(`db info error`, err)
// 		}
// 	}
//
// 	doc, err = db.Get(uuid, url.Values{"rev": []string{oldRev}})
// 	if err == nil {
// 		t.Errorf("db get compacted doc ok, rev = %s", oldRev)
// 	}
// }
//
// func TestAttachmentCRUD(t *testing.T) {
// 	uuid := GenerateUUID()
// 	doc := map[string]interface{}{"bar": 42}
// 	db.Set(uuid, doc)
// 	oldRev := doc["_rev"].(string)
//
// 	db.PutAttachment(doc, []byte("Foo bar"), "foo.txt", "text/plain")
// 	if oldRev == doc["_rev"].(string) {
// 		t.Error(`doc[_rev] == oldRev`)
// 	}
//
// 	doc, err := db.Get(uuid, nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
// 	attachments := reflect.ValueOf(doc["_attachments"])
// 	foo := reflect.ValueOf(attachments.MapIndex(reflect.ValueOf("foo.txt")).Interface())
// 	length := int(foo.MapIndex(reflect.ValueOf("length")).Interface().(float64))
// 	if length != len("Foo bar") {
// 		t.Errorf("length %d want %d", length, len("Foo bar"))
// 	}
// 	contentType := foo.MapIndex(reflect.ValueOf("content_type")).Interface().(string)
// 	if contentType != "text/plain" {
// 		t.Errorf("content type %s want text/plain", contentType)
// 	}
//
// 	data, err := db.GetAttachment(doc, "foo.txt")
// 	if err != nil {
// 		t.Error(`get attachment error`, err)
// 	}
// 	if string(data) != "Foo bar" {
// 		t.Errorf("db get attachment %s want Foo bar", string(data))
// 	}
//
// 	data, err = db.GetAttachmentID(uuid, "foo.txt")
// 	if err != nil {
// 		t.Error(`get attachment id error`, err)
// 	}
// 	if string(data) != "Foo bar" {
// 		t.Errorf("db get attachment id %s want Foo bar", string(data))
// 	}
//
// 	oldRev = doc["_rev"].(string)
// 	err = db.DeleteAttachment(doc, "foo.txt")
// 	if err != nil {
// 		t.Error(`db delete attachment error`, err)
// 	}
// 	if oldRev == doc["_rev"].(string) {
// 		t.Error(`doc[_rev] == oldRev`)
// 	}
//
// 	doc, err = db.Get(uuid, nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
// 	if _, ok := doc["_attachments"]; ok {
// 		t.Error(`doc attachments still existed`)
// 	}
// }
//
// func TestAttachmentWithFiles(t *testing.T) {
// 	uuid := GenerateUUID()
// 	doc := map[string]interface{}{"bar": 42}
// 	err := db.Set(uuid, doc)
// 	if err != nil {
// 		t.Error(`db set doc error`, err)
// 	}
// 	oldRev := doc["_rev"].(string)
// 	fileObj := []byte("Foo bar baz")
//
// 	err = db.PutAttachment(doc, fileObj, "foo.txt", mime.TypeByExtension(".txt"))
// 	if oldRev == doc["_rev"].(string) {
// 		t.Error(`doc[_rev] == oldRev`)
// 	}
//
// 	doc, err = db.Get(uuid, nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
// 	attachments := reflect.ValueOf(doc["_attachments"])
// 	foo := reflect.ValueOf(attachments.MapIndex(reflect.ValueOf("foo.txt")).Interface())
// 	length := int(foo.MapIndex(reflect.ValueOf("length")).Interface().(float64))
// 	if length != len("Foo bar baz") {
// 		t.Errorf("length %d want %d", length, len("Foo bar"))
// 	}
// 	contentType := foo.MapIndex(reflect.ValueOf("content_type")).Interface().(string)
// 	if contentType != "text/plain; charset=utf-8" {
// 		t.Errorf("content type %s want text/plain; charset=utf-8", contentType)
// 	}
//
// 	data, err := db.GetAttachment(doc, "foo.txt")
// 	if err != nil {
// 		t.Error(`get attachment error`, err)
// 	}
// 	if string(data) != "Foo bar baz" {
// 		t.Errorf("db get attachment %s want Foo bar", string(data))
// 	}
//
// 	data, err = db.GetAttachmentID(uuid, "foo.txt")
// 	if err != nil {
// 		t.Error(`get attachment id error`, err)
// 	}
// 	if string(data) != "Foo bar baz" {
// 		t.Errorf("db get attachment id %s want Foo bar", string(data))
// 	}
//
// 	oldRev = doc["_rev"].(string)
// 	err = db.DeleteAttachment(doc, "foo.txt")
// 	if err != nil {
// 		t.Error(`db delete attachment error`, err)
// 	}
// 	if oldRev == doc["_rev"].(string) {
// 		t.Error(`doc[_rev] == oldRev`)
// 	}
//
// 	doc, err = db.Get(uuid, nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
// 	if _, ok := doc["_attachments"]; ok {
// 		t.Error(`doc attachments still existed`)
// 	}
// }
//
// func TestAttachmentCRUDFromFS(t *testing.T) {
// 	uuid := GenerateUUID()
// 	content := "Foo bar baz"
// 	tmpFileName := filepath.Join(os.TempDir(), "foo.txt")
// 	tmpFile, err := os.Create(tmpFileName)
// 	if err != nil {
// 		t.Error(`create file error`, err)
// 	}
// 	_, err = tmpFile.Write([]byte(content))
// 	if err != nil {
// 		t.Error(`write file error`, err)
// 	}
// 	tmpFile.Close()
//
// 	tmpFile, err = os.Open(tmpFileName)
// 	if err != nil {
// 		t.Error(`open file error`, err)
// 	}
// 	defer tmpFile.Close()
//
// 	data, err := ioutil.ReadAll(tmpFile)
// 	if err != nil {
// 		t.Error(`read tmp file error`, err)
// 	}
//
// 	doc := map[string]interface{}{"bar": 42}
// 	err = db.Set(uuid, doc)
// 	if err != nil {
// 		t.Error(`db set doc error`, err)
// 	}
// 	oldRev := doc["_rev"].(string)
// 	err = db.PutAttachment(doc, data, "foo.txt", mime.TypeByExtension(filepath.Ext(tmpFileName)))
// 	if err != nil {
// 		t.Error(`put attachment error`, err)
// 	}
// 	if oldRev == doc["_rev"].(string) {
// 		t.Error(`doc[_rev] == oldRev`)
// 	}
//
// 	doc, err = db.Get(uuid, nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
// 	attachment := reflect.ValueOf(doc["_attachments"])
// 	foo := reflect.ValueOf(attachment.MapIndex(reflect.ValueOf("foo.txt")).Interface())
// 	length := int(foo.MapIndex(reflect.ValueOf("length")).Interface().(float64))
// 	if len(content) != length {
// 		t.Errorf("length %d want %d", length, len(content))
// 	}
// 	contentType := foo.MapIndex(reflect.ValueOf("content_type")).Interface().(string)
// 	if contentType != "text/plain; charset=utf-8" {
// 		t.Errorf("content type %s want text/plain; charset=utf-8", contentType)
// 	}
//
// 	data, err = db.GetAttachment(doc, "foo.txt")
// 	if err != nil {
// 		t.Error(`get attachment error`, err)
// 	}
// 	if string(data) != content {
// 		t.Error(`get attachment should be `, content)
// 	}
//
// 	data, err = db.GetAttachmentID(uuid, "foo.txt")
// 	if err != nil {
// 		t.Error(`get attachment id error`, err)
// 	}
// 	if string(data) != content {
// 		t.Error(`get attachment id should be `, content)
// 	}
//
// 	if err = db.DeleteAttachment(doc, "foo.txt"); err != nil {
// 		t.Error(`delete attachment file error`, err)
// 	}
// 	if oldRev == doc["_rev"].(string) {
// 		t.Error(`doc[_rev] == oldRev`)
// 	}
//
// 	doc, err = db.Get(uuid, nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
// 	if _, ok := doc["_attachments"]; ok {
// 		t.Error(`doc attachments still existed`)
// 	}
// }
//
// func TestEmptyAttachment(t *testing.T) {
// 	uuid := GenerateUUID()
// 	doc := map[string]interface{}{}
// 	err := db.Set(uuid, doc)
// 	if err != nil {
// 		t.Error(`db set doc error`, err)
// 	}
// 	oldRev := doc["_rev"].(string)
//
// 	err = db.PutAttachment(doc, []byte(""), "empty.txt", mime.TypeByExtension(".txt"))
// 	if err != nil {
// 		t.Error(`put attachment error`, err)
// 	}
// 	if oldRev == doc["_rev"].(string) {
// 		t.Error(`doc[_rev] == oldRev`)
// 	}
//
// 	doc, err = db.Get(uuid, nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
//
// 	attachment := reflect.ValueOf(doc["_attachments"])
// 	empty := reflect.ValueOf(attachment.MapIndex(reflect.ValueOf("empty.txt")).Interface())
// 	length := int(empty.MapIndex(reflect.ValueOf("length")).Interface().(float64))
// 	if length != 0 {
// 		t.Errorf("length %d want %d", length, 0)
// 	}
// }
//
// func TestDefaultAttachment(t *testing.T) {
// 	uuid := GenerateUUID()
// 	doc := map[string]interface{}{}
// 	err := db.Set(uuid, doc)
// 	if err != nil {
// 		t.Error(`db set doc error`, err)
// 	}
// 	_, err = db.GetAttachment(doc, "missing.txt")
// 	if err == nil {
// 		t.Error(`db get attachment ok`)
// 	}
// }
//
// func TestAttachmentNoFilename(t *testing.T) {
// 	uuid := GenerateUUID()
// 	doc := map[string]interface{}{}
// 	err := db.Set(uuid, doc)
// 	if err != nil {
// 		t.Error(`db set doc error`, err)
// 	}
// 	err = db.PutAttachment(doc, []byte(""), "", "")
// 	if err == nil {
// 		t.Error(`db put attachment with no file name ok`)
// 	}
// }
//
// func TestJSONAttachment(t *testing.T) {
// 	doc := map[string]interface{}{}
// 	err := db.Set(GenerateUUID(), doc)
// 	if err != nil {
// 		t.Error(`db set doc error`, err)
// 	}
// 	err = db.PutAttachment(doc, []byte("{}"), "test.json", "application/json")
// 	if err != nil {
// 		t.Error(`db put attachment json error`, err)
// 	}
// 	data, err := db.GetAttachment(doc, "test.json")
// 	if err != nil {
// 		t.Error(`db get attachment json error`, err)
// 	}
// 	if string(data) != "{}" {
// 		t.Errorf("data = %s want {}", string(data))
// 	}
// }
//
// func TestBulkUpdateConflict(t *testing.T) {
// 	docs := []map[string]interface{}{
// 		{
// 			"type": "Person",
// 			"name": "John Doe",
// 		},
// 		{
// 			"type": "Person",
// 			"name": "Mary Jane",
// 		},
// 		{
// 			"type": "Person",
// 			"name": "Gotham City",
// 		},
// 	}
//
// 	db.Update(docs, nil)
//
// 	// update the first doc to provoke a conflict in the next bulk update
// 	doc := map[string]interface{}{}
// 	for k, v := range docs[0] {
// 		doc[k] = v
// 	}
// 	db.Set(doc["_id"].(string), doc)
//
// 	results, err := db.Update(docs, nil)
// 	if err != nil {
// 		t.Error(`db update error`, err)
// 	}
// 	if results[0].err != ErrConflict {
// 		t.Errorf("db update conflict err %v want ErrConflict", results[0].err)
// 	}
// }
//
// func TestCopyDocConflict(t *testing.T) {
// 	db.Set("foo1", map[string]interface{}{"status": "idle"})
// 	db.Set("bar1", map[string]interface{}{"status": "testing"})
// 	_, err := db.Copy("foo1", "bar1", "")
// 	if err != ErrConflict {
// 		t.Errorf(`db copy returns %v, want ErrConflict`, err)
// 	}
// }
//
// func TestCopyDocOverwrite(t *testing.T) {
// 	foo2 := map[string]interface{}{"status": "testing"}
// 	bar2 := map[string]interface{}{"status": "idle"}
// 	db.Set("foo2", foo2)
// 	db.Set("bar2", bar2)
// 	result, err := db.Copy("foo2", "bar2", bar2["_rev"].(string))
// 	if err != nil {
// 		t.Error(`db copy error`, err)
// 	}
// 	doc, _ := db.Get("bar2", nil)
// 	if result != doc["_rev"].(string) {
// 		t.Errorf("db copy returns %s want %s", result, doc["_rev"].(string))
// 	}
// 	if doc["status"].(string) != "testing" {
// 		t.Errorf("db copy status = %s, want testing", doc["status"].(string))
// 	}
// }
//
// func TestChanges(t *testing.T) {
// 	options := url.Values{
// 		"style": []string{"all_docs"},
// 	}
// 	_, err := db.Changes(options)
// 	if err != nil {
// 		t.Error(`db change error`, err)
// 	}
// }
//
// // FIXME: Purge not implemented in CouchDB 2.0.0 yet.
// func TestPurge(t *testing.T) {
// 	doc := map[string]interface{}{"a": "b"}
// 	err := db.Set("purge", doc)
// 	if err != nil {
// 		t.Error(`db set error`, err)
// 	}
// 	_, err = db.Purge([]map[string]interface{}{doc})
// 	if err == nil {
// 		t.Error(`db purge ok`, err)
// 	}
// 	/*
// 		purgeSeq := int(result["purge_seq"].(float64))
// 		if purgeSeq != 1 {
// 			t.Errorf("db purge seq=%d want 1", purgeSeq)
// 		}
// 	*/
// }
//
// func TestSecurity(t *testing.T) {
// 	secDoc, err := db.GetSecurity()
// 	if err != nil {
// 		t.Error(`get security should return true`)
// 	}
// 	if len(secDoc) > 0 {
// 		t.Error(`secDoc should be empty`)
// 	}
// 	if db.SetSecurity(map[string]interface{}{
// 		"names": []string{"test"},
// 		"roles": []string{},
// 	}) != nil {
// 		t.Error(`set security should return true`)
// 	}
// }
//
// func TestDBContains(t *testing.T) {
// 	doc := map[string]interface{}{
// 		"type": "Person",
// 		"name": "Jason Statham",
// 	}
// 	id, _, err := db.Save(doc, nil)
// 	if err != nil {
// 		t.Error(`db save error`, err)
// 	}
// 	if err = db.Contains(id); err != nil {
// 		t.Error(`db contains error`, err)
// 	}
// }
//
// func TestDBSetGetDelete(t *testing.T) {
// 	doc := map[string]interface{}{
// 		"type": "Person",
// 		"name": "Jason Statham",
// 	}
// 	err := db.Set("Mechanic", doc)
// 	if err != nil {
// 		t.Error(`db set error`, err)
// 	}
// 	_, err = db.Get("Mechanic", nil)
// 	if err != nil {
// 		t.Error(`db get error`, err)
// 	}
// 	err = db.Delete("Mechanic")
// 	if err != nil {
// 		t.Error(`db delete error`, err)
// 	}
// }
//
// func TestDBDocIDsAndLen(t *testing.T) {
// 	doc := map[string]interface{}{
// 		"type": "Person",
// 		"name": "Jason Statham",
// 	}
//
// 	err := db.Set("Mechanic", doc)
// 	if err != nil {
// 		t.Error(`db set error`, err)
// 	}
//
// 	ids, err := db.DocIDs()
// 	if err != nil {
// 		t.Error(`db doc ids error`, err)
// 	}
//
// 	length, err := db.Len()
// 	if err != nil {
// 		t.Error(`db len error`, err)
// 	}
// 	if length != len(ids) {
// 		t.Errorf("Len() returns %d want %d", length, len(ids))
// 	}
// }
//
// func TestGetSetRevsLimit(t *testing.T) {
// 	err := db.SetRevsLimit(10)
// 	if err != nil {
// 		t.Error(`db set revs limit error`, err)
// 	}
// 	limit, err := db.GetRevsLimit()
// 	if err != nil {
// 		t.Error(`db get revs limit error`, err)
// 	}
// 	if limit != 10 {
// 		t.Error(`limit should be 10`)
// 	}
// }
//
// func TestCleanup(t *testing.T) {
// 	err := db.Cleanup()
// 	if err != nil {
// 		t.Error(`db clean up error`, err)
// 	}
// }

// TODO adding new apis for mango query engine
// func TestQuery(t *testing.T) {
// result, err := parseSelectorSyntax(`title == "Spacecataz" && year == 2004 && director == "Dave Willis"`)
// result, err := parseSelectorSyntax(`title == "Spacecataz" && year == 2004`)
// result, err := parseSelectorSyntax(`year == 2004`)
// result, err := parseSelectorSyntax(`year >= 1990 && (director == "George Lucas" || director == "Steven Spielberg")`)
// result, err := parseSelectorSyntax(`director == "George Lucas" || director == "Steven Spielberg"`)
// result, err := parseSelectorSyntax(`year >= 1900 && year <= 2000 && nor(year == 1990, year == 1989, year == 1997)`)
// result, err := parseSelectorSyntax("year >= 1990 && year <= 1910")
// result, err := parseSelectorSyntax(`_id > nil && all(genre, []string{"Comedy", "Short"})`)
// result, err := parseSelectorSyntax(`_id > nil && any(genre, genre == "Horror" || genre == "Comedy" || genre == "Short")`)
// result, err := parseSelectorSyntax(`_id > nil && any(genre, genre == "Horror" || genre == "Short" || score >= 8)`)
// result, err := parseSelectorSyntax(`exists(director, true)`)
// result, err := parseSelectorSyntax(`typeof(genre, "array")`)
// result, err := parseSelectorSyntax(`in(director, []string{"Mike Portnoy", "Vitali Kanevsky"})`)
// result, err := parseSelectorSyntax(`nin(year, []int{1990, 1992, 1998})`)
// result, err := parseSelectorSyntax(`size(genre, 2)`)
// result, err := parseSelectorSyntax(`mod(year, 2, 1)`)
// result, err := parseSelectorSyntax(`regex(title, "^A")`)
// result, err := parseSortSyntax([]string{"fieldNameA", "fieldNameB"})
// result, err := parseSortSyntax([]string{"fieldNameA.subFieldA", "fieldNameB.subFieldB"})
// result, err := parseSortSyntax([]string{"desc(fieldName1)", "asc(fieldName2)"})
// result, err := parseSortSyntax([]string{"desc(fieldName1.subField1)", "asc(fieldName2.subField2)"})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	s, _ := beautifulJSONString(result)
// 	fmt.Println(s)
// }

func TestQueryYearAndID(t *testing.T) {
	docsQuery, err := movieDB.Query(nil, `_id > nil && in(year, []int{2007, 2004})`, nil, nil, nil, nil)
	if err != nil {
		t.Error("db query error", err)
	}

	var rawJSON = []byte(`
	{
	    "selector": {
				"$and": [
	        	{
	            "_id": { "$gt": null }
	        	},
	        	{
	            "year": {
	              "$in": [2007, 2004]
	          	}
	        	}
	    	]
			}
	}`)
	queryMap := map[string]interface{}{}
	err = json.Unmarshal(rawJSON, &queryMap)
	if err != nil {
		t.Error("json unmarshal error", err)
	}

	docsRaw, err := movieDB.QueryJSON(queryMap)
	if err != nil {
		t.Error("db query json error", err)
	}

	if !reflect.DeepEqual(docsQuery, docsRaw) {
		t.Error("db query year and id not equal")
	}
}

func TestQueryYearOrDirector(t *testing.T) {
	docsQuery, err := movieDB.Query(nil, `year == 1989 && (director == "Ademir Kenovic" || director == "Dezs Garas")`, nil, nil, nil, nil)
	if err != nil {
		t.Error("db query error", err)
	}

	var rawJSON = []byte(`
	{
		"selector": {
			"year": 1989,
			"$or": [
				{ "director": "Ademir Kenovic" },
				{ "director": "Dezs Garas" }
			]
		}
	}`)
	queryMap := map[string]interface{}{}
	err = json.Unmarshal(rawJSON, &queryMap)
	if err != nil {
		t.Error("json unmarshal error", err)
	}

	docsRaw, err := movieDB.QueryJSON(queryMap)
	if err != nil {
		t.Error("db query json error", err)
	}

	if !reflect.DeepEqual(docsQuery, docsRaw) {
		t.Error("db query year or director not equal")
	}
}

func TestQueryYearGteLteNot(t *testing.T) {
	docsQuery, err := movieDB.Query(nil, `year >= 1989 && year <= 2006 && year != 2004`, nil, nil, nil, nil)
	if err != nil {
		t.Error("db query error", err)
	}

	var rawJSON = []byte(`
	{
		"selector": {
			"year": {
	      "$gte": 1989
	    },
	    "year": {
	      "$lte": 2006
	    },
	    "$not": {
	      "year": 2004
	    }
		}
	}`)
	queryMap := map[string]interface{}{}
	err = json.Unmarshal(rawJSON, &queryMap)
	if err != nil {
		t.Error("json unmarshal error", err)
	}

	docsRaw, err := movieDB.QueryJSON(queryMap)
	if err != nil {
		t.Error("db query json error", err)
	}

	if !reflect.DeepEqual(docsQuery, docsRaw) {
		t.Error("db query year gte lte not not equal")
	}
}

// {
//     "imdb.rating": {
//         "$gte": 6
//     },
//     "imdb.rating": {
//         "$lte": 9
//     },
//     "$nor": [
//         { "imdb.rating": 8.1 },
//         { "imdb.rating": 8.2 },
//         { "imdb.rating": 7.8 }
//     ]
// }
func TestQueryIMDBRatingNor(t *testing.T) {}

// {
//     "_id": {
//         "$gt": null
//     },
//     "genre": {
//         "$all": ["Comedy","Short"]
//     }
// }
func TestQueryGenreAll(t *testing.T) {}

// {
//     "_id": { "$gt": null },
//     "genre": {
//         "$elemMatch": {
//             "$eq": "Horror"
//         }
//     }
// }
func TestQueryGenreElemMatch(t *testing.T) {}

// {
//     "selector": {
//         "afieldname": {"$regex": "^A"}
//     }
// }
func TestQueryRegex(t *testing.T) {}

// {
//     "$and": [
//         {
//             "_id": { "$gt": null }
//         },
//         {
//             "year": {
//                 "$nin": [1989, 1990]
//             }
//         }
//     ]
// }
func TestQueryYearIDNin(t *testing.T) {}

// {
//     "$and": [
//         {
// 	    "poster": {
// 	        "$type": "string"
// 	    }
// 	},
// 	{
// 	    "runtime": {
//     	        "$exists": true
// 	    }
// 	}
//     ]
// }
func TestQueryTypeAndExists(t *testing.T) {}

// {
//     "$and": [
//         {
// 	    "writer": {
// 		"$size": 2
// 	    }
// 	},
// 	{
// 	    "year": {
// 		"$mod": [2, 0]
// 	    }
// 	}
//     ]
// }
func TestQuerySizeAndMod(t *testing.T) {}

// {
//     "$or": [
//         {
// 	    "rating": {
// 		"$ne": null
// 	    }
// 	},
// 	{
// 	    "year": {
//   		"$lt": 2000
// 	    }
// 	}
//     ]
// }
func TestRatingOrYear(t *testing.T) {}

// {
//     "selector": {
//         "year": {"$gt": 1989}
//     },
//     "fields": ["_id", "_rev", "year", "title"],
//     "sort": [{"year": "asc"}],
//     "limit": 5,
//     "skip": 2,
// }
func TestQuerySortLimitSkip(t *testing.T) {}

// {
//     "selector": {
//         "year": {"$gt": 1989}
//     },
//     "fields": ["_id", "_rev", "year", "title"],
//     "sort": [{"imdb.rating": "desc"}, {"imdb.votes": "desc"}],
// }
func TestQueryDoubleSort(t *testing.T) {}

func TestIndexCRUD(t *testing.T) {}

// {
//     "selector": {
//         "year": {"$gt": 1989}
//     },
//     "fields": ["_id", "_rev", "year", "title"],
//     "sort": [{"year": "asc"}],
//     "limit": 5,
//     "skip": 2,
//	   "use_index": xxx
// }
func TestQueryUseIndex(t *testing.T) {}
