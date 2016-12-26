package couchdb

import (
	"bytes"
	"io/ioutil"
	"math"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewDefaultDB(t *testing.T) {
	dbDefault, err := NewDatabase("golang-default")
	if err != nil {
		t.Errorf("new default database error %v", err)
	}
	if dbDefault.Available() {
		t.Error(`db available`)
	}
}

func TestNewDB(t *testing.T) {
	newDB := "golang-newdb"
	s.Create(newDB)
	defer s.Delete(newDB)
	dbNew, err := NewDatabase("http://root:likejun@localhost:5984/" + newDB)
	if err != nil {
		t.Error(`new database error`, err)
	}
	if !dbNew.Available() {
		t.Error(`db not available`)
	}
}

func TestSaveNew(t *testing.T) {
	doc := map[string]interface{}{"doc": "bar"}
	id, rev, err := db.Save(doc, nil)
	if err != nil {
		t.Error(`db save error`, err)
	}
	if id != doc["_id"].(string) {
		t.Errorf("invalid id: %q", id)
	}
	if rev != doc["_rev"].(string) {
		t.Errorf("invalid rev: %q", rev)
	}
}

func TestSaveNewWithID(t *testing.T) {
	doc := map[string]interface{}{"_id": "foo"}
	id, rev, err := db.Save(doc, nil)
	if err != nil {
		t.Error(`db save error`, err)
	}
	if doc["_id"].(string) != "foo" {
		t.Errorf("doc[_id] = %s, not foo", doc["_id"])
	}
	if id != "foo" {
		t.Errorf("id = %s, not foo", id)
	}
	if rev != doc["_rev"].(string) {
		t.Errorf("invalid rev: %q", rev)
	}
}

func TestSaveExisting(t *testing.T) {
	doc := map[string]interface{}{}
	idOld, revOld, err := db.Save(doc, nil)
	if err != nil {
		t.Error(`db save error`, err)
	}
	doc["foo"] = true
	idNew, revNew, err := db.Save(doc, nil)
	if err != nil {
		t.Error(`db save foo error`, err)
	}
	if idOld != idNew {
		t.Errorf("ids are not equal old %s new %s", idOld, idNew)
	}
	if doc["_rev"].(string) != revNew {
		t.Errorf("invalid rev %s want %s", doc["_rev"].(string), revNew)
	}
	if revOld == revNew {
		t.Errorf("new rev is equal to old rev %s", revOld)
	}
}

func TestSaveNewBatch(t *testing.T) {
	doc := map[string]interface{}{"_id": "foo"}
	_, rev, err := db.Save(doc, url.Values{"batch": []string{"ok"}})
	if err != nil {
		t.Error(`db save batch error`, err)
	}
	if len(rev) > 0 {
		t.Error(`rev not empty`, rev)
	}
	if r, ok := doc["_rev"]; ok {
		t.Error(`doc has _rev field`, r.(string))
	}
}

func TestSaveExistingBatch(t *testing.T) {
	doc := map[string]interface{}{"_id": "bar"}
	idOld, revOld, err := db.Save(doc, nil)
	if err != nil {
		t.Error(`db save error`, err)
	}

	idNew, revNew, err := db.Save(doc, url.Values{"batch": []string{"ok"}})
	if err != nil {
		t.Error(`db save batch error`, err)
	}

	if idOld != idNew {
		t.Errorf("old id %s not equal to new id %s", idOld, idNew)
	}

	if len(revNew) > 0 {
		t.Error(`rev not empty`, revNew)
	}

	if doc["_rev"].(string) != revOld {
		t.Errorf("doc[_rev] %s not equal to old rev %s", doc["_rev"].(string), revOld)
	}
}

func TestDatabaseExists(t *testing.T) {
	if !db.Available() {
		t.Error(`golang-tests not available`)
	}
	dbMissing, _ := NewDatabase("golang-missing")
	if dbMissing.Available() {
		t.Error(`golang-missing available`)
	}
}

func TestDatabaseName(t *testing.T) {
	name, err := db.Name()
	if err != nil {
		t.Error(`db name error`, err)
	}
	if name != "golang-tests" {
		t.Error("db name %s, want golang-tests", name)
	}
}

func TestDatabaseString(t *testing.T) {
	if db.String() != "Database http://root:likejun@localhost:5984/golang-tests" {
		t.Error(`db string invalid`, db)
	}
}

func TestCommit(t *testing.T) {
	if err := db.Commit(); err != nil {
		t.Error(`db commit error`, err)
	}
}

func TestCreateLargeDoc(t *testing.T) {
	var buf bytes.Buffer
	// 10MB
	for i := 0; i < 110*1024; i++ {
		buf.WriteString("0123456789")
	}
	doc := map[string]interface{}{"data": buf.String()}
	if err := db.Set("large", doc); err != nil {
		t.Error(`db set error`, err)
	}
	doc, err := db.Get("large", nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if doc["_id"].(string) != "large" {
		t.Errorf("doc[_id] = %s, want large", doc["_id"].(string))
	}
	err = db.DeleteDoc(doc)
	if err != nil {
		t.Error(`db delete doc error`, err)
	}
}

func TestDocIDQuoting(t *testing.T) {
	doc := map[string]interface{}{"foo": "bar"}
	err := db.Set("foo/bar", doc)
	if err != nil {
		t.Error(`db set error`, err)
	}
	doc, err = db.Get("foo/bar", nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if doc["foo"].(string) != "bar" {
		t.Errorf("doc[foo] = %s want bar", doc["foo"].(string))
	}
	err = db.Delete("foo/bar")
	if err != nil {
		t.Error(`db delete error`, err)
	}
	_, err = db.Get("foo/bar", nil)
	if err == nil {
		t.Error(`db get foo/bar ok`)
	}
}

func TestDisallowNaN(t *testing.T) {
	doc := map[string]interface{}{"number": math.NaN()}
	err := db.Set("foo", doc)
	if err == nil {
		t.Error(`db set NaN ok`)
	}
}

func TestDisallowNilID(t *testing.T) {
	err := db.DeleteDoc(map[string]interface{}{"_id": nil, "_rev": nil})
	if err == nil {
		t.Error(`db delete doc with id nil ok`)
	}
	err = db.DeleteDoc(map[string]interface{}{"_id": "foo", "_rev": nil})
	if err == nil {
		t.Error(`db delete doc with rev nil ok`)
	}
}

func TestDocRevs(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{"bar": 42}
	err := db.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	oldRev := doc["_rev"].(string)
	doc["bar"] = 43
	err = db.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	newRev := doc["_rev"].(string)

	newDoc, err := db.Get(uuid, nil)
	if newRev != newDoc["_rev"].(string) {
		t.Errorf("new doc rev %s not equal to %s", newDoc["_rev"].(string), newRev)
	}
	newDoc, err = db.Get(uuid, url.Values{"rev": []string{newRev}})
	if newRev != newDoc["_rev"].(string) {
		t.Errorf("new doc rev %s not equal to %s", newDoc["_rev"].(string), newRev)
	}
	oldDoc, err := db.Get(uuid, url.Values{"rev": []string{oldRev}})
	if oldRev != oldDoc["_rev"].(string) {
		t.Errorf("old doc rev %s not equal to %s", oldDoc["_rev"].(string), oldRev)
	}

	revs, err := db.Revisions(uuid, nil)
	if err != nil {
		t.Error(`db revisions error`, err)
	}
	if revs[0]["_rev"].(string) != newRev {
		t.Errorf("revs first %s not equal to %s", revs[0]["_rev"].(string), newRev)
	}
	if revs[1]["_rev"].(string) != oldRev {
		t.Errorf("revs second %s not equal to %s", revs[1]["_rev"].(string), oldRev)
	}
	_, err = db.Revisions("crap", nil)
	if err == nil {
		t.Error(`db revisions crap ok`)
	}

	err = db.Compact()
	if err != nil {
		t.Error("db compact error", err)
	}

	info, err := db.Info()
	if err != nil {
		t.Error(`db info error`, err)
	}
	for info["compact_running"].(bool) {
		info, err = db.Info()
		if err != nil {
			t.Error(`db info error`, err)
		}
	}

	doc, err = db.Get(uuid, url.Values{"rev": []string{oldRev}})
	if err == nil {
		t.Errorf("db get compacted doc ok, rev = %s", oldRev)
	}
}

func TestAttachmentCRUD(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{"bar": 42}
	db.Set(uuid, doc)
	oldRev := doc["_rev"].(string)

	db.PutAttachment(doc, []byte("Foo bar"), "foo.txt", "text/plain")
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err := db.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	attachments := reflect.ValueOf(doc["_attachments"])
	foo := reflect.ValueOf(attachments.MapIndex(reflect.ValueOf("foo.txt")).Interface())
	length := int(foo.MapIndex(reflect.ValueOf("length")).Interface().(float64))
	if length != len("Foo bar") {
		t.Errorf("length %d want %d", length, len("Foo bar"))
	}
	contentType := foo.MapIndex(reflect.ValueOf("content_type")).Interface().(string)
	if contentType != "text/plain" {
		t.Errorf("content type %s want text/plain", contentType)
	}

	data, err := db.GetAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`get attachment error`, err)
	}
	if string(data) != "Foo bar" {
		t.Errorf("db get attachment %s want Foo bar", string(data))
	}

	data, err = db.GetAttachmentID(uuid, "foo.txt")
	if err != nil {
		t.Error(`get attachment id error`, err)
	}
	if string(data) != "Foo bar" {
		t.Errorf("db get attachment id %s want Foo bar", string(data))
	}

	oldRev = doc["_rev"].(string)
	err = db.DeleteAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`db delete attachment error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = db.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if _, ok := doc["_attachments"]; ok {
		t.Error(`doc attachments still existed`)
	}
}

func TestAttachmentWithFiles(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{"bar": 42}
	err := db.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	oldRev := doc["_rev"].(string)
	fileObj := []byte("Foo bar baz")

	err = db.PutAttachment(doc, fileObj, "foo.txt", mime.TypeByExtension(".txt"))
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = db.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	attachments := reflect.ValueOf(doc["_attachments"])
	foo := reflect.ValueOf(attachments.MapIndex(reflect.ValueOf("foo.txt")).Interface())
	length := int(foo.MapIndex(reflect.ValueOf("length")).Interface().(float64))
	if length != len("Foo bar baz") {
		t.Errorf("length %d want %d", length, len("Foo bar"))
	}
	contentType := foo.MapIndex(reflect.ValueOf("content_type")).Interface().(string)
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("content type %s want text/plain; charset=utf-8", contentType)
	}

	data, err := db.GetAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`get attachment error`, err)
	}
	if string(data) != "Foo bar baz" {
		t.Errorf("db get attachment %s want Foo bar", string(data))
	}

	data, err = db.GetAttachmentID(uuid, "foo.txt")
	if err != nil {
		t.Error(`get attachment id error`, err)
	}
	if string(data) != "Foo bar baz" {
		t.Errorf("db get attachment id %s want Foo bar", string(data))
	}

	oldRev = doc["_rev"].(string)
	err = db.DeleteAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`db delete attachment error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = db.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if _, ok := doc["_attachments"]; ok {
		t.Error(`doc attachments still existed`)
	}
}

func TestAttachmentCRUDFromFS(t *testing.T) {
	uuid := GenerateUUID()
	content := "Foo bar baz"
	tmpFileName := filepath.Join(os.TempDir(), "foo.txt")
	tmpFile, err := os.Create(tmpFileName)
	if err != nil {
		t.Error(`create file error`, err)
	}
	_, err = tmpFile.Write([]byte(content))
	if err != nil {
		t.Error(`write file error`, err)
	}
	tmpFile.Close()

	tmpFile, err = os.Open(tmpFileName)
	if err != nil {
		t.Error(`open file error`, err)
	}
	defer tmpFile.Close()

	data, err := ioutil.ReadAll(tmpFile)
	if err != nil {
		t.Error(`read tmp file error`, err)
	}

	doc := map[string]interface{}{"bar": 42}
	err = db.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	oldRev := doc["_rev"].(string)
	err = db.PutAttachment(doc, data, "foo.txt", mime.TypeByExtension(filepath.Ext(tmpFileName)))
	if err != nil {
		t.Error(`put attachment error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = db.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	attachment := reflect.ValueOf(doc["_attachments"])
	foo := reflect.ValueOf(attachment.MapIndex(reflect.ValueOf("foo.txt")).Interface())
	length := int(foo.MapIndex(reflect.ValueOf("length")).Interface().(float64))
	if len(content) != length {
		t.Errorf("length %d want %d", length, len(content))
	}
	contentType := foo.MapIndex(reflect.ValueOf("content_type")).Interface().(string)
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("content type %s want text/plain; charset=utf-8", contentType)
	}

	data, err = db.GetAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`get attachment error`, err)
	}
	if string(data) != content {
		t.Error(`get attachment should be `, content)
	}

	data, err = db.GetAttachmentID(uuid, "foo.txt")
	if err != nil {
		t.Error(`get attachment id error`, err)
	}
	if string(data) != content {
		t.Error(`get attachment id should be `, content)
	}

	if err = db.DeleteAttachment(doc, "foo.txt"); err != nil {
		t.Error(`delete attachment file error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = db.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if _, ok := doc["_attachments"]; ok {
		t.Error(`doc attachments still existed`)
	}
}

func TestEmptyAttachment(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{}
	err := db.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	oldRev := doc["_rev"].(string)

	err = db.PutAttachment(doc, []byte(""), "empty.txt", mime.TypeByExtension(".txt"))
	if err != nil {
		t.Error(`put attachment error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = db.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}

	attachment := reflect.ValueOf(doc["_attachments"])
	empty := reflect.ValueOf(attachment.MapIndex(reflect.ValueOf("empty.txt")).Interface())
	length := int(empty.MapIndex(reflect.ValueOf("length")).Interface().(float64))
	if length != 0 {
		t.Errorf("length %d want %d", length, 0)
	}
}

func TestDefaultAttachment(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{}
	err := db.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	_, err = db.GetAttachment(doc, "missing.txt")
	if err == nil {
		t.Error(`db get attachment ok`)
	}
}

func TestAttachmentNoFilename(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{}
	err := db.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	err = db.PutAttachment(doc, []byte(""), "", "")
	if err == nil {
		t.Error(`db put attachment with no file name ok`)
	}
}

func TestJSONAttachment(t *testing.T) {
	doc := map[string]interface{}{}
	err := db.Set(GenerateUUID(), doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	err = db.PutAttachment(doc, []byte("{}"), "test.json", "application/json")
	if err != nil {
		t.Error(`db put attachment json error`, err)
	}
	data, err := db.GetAttachment(doc, "test.json")
	if err != nil {
		t.Error(`db get attachment json error`, err)
	}
	if string(data) != "{}" {
		t.Errorf("data = %s want {}", string(data))
	}
}

// func TestIncludeDocs() {}
// // TODO adding new apis about mango query engine
// func TestQueryMultiGet() {}
func TestBulkUpdateConflict(t *testing.T) {
	docs := []map[string]interface{}{
		{
			"type": "Person",
			"name": "John Doe",
		},
		{
			"type": "Person",
			"name": "Mary Jane",
		},
		{
			"type": "Person",
			"name": "Gotham City",
		},
	}

	db.Update(docs, nil)

	// update the first doc to provoke a conflict in the next bulk update
	doc := map[string]interface{}{}
	for k, v := range docs[0] {
		doc[k] = v
	}
	db.Set(doc["_id"].(string), doc)

	results, err := db.Update(docs, nil)
	if err != nil {
		t.Error(`db update error`, err)
	}
	if results[0].err != ErrConflict {
		t.Errorf("db update conflict err %v want ErrConflict", results[0].err)
	}
}

func TestCopyDocConflict(t *testing.T) {
	db.Set("foo1", map[string]interface{}{"status": "idle"})
	db.Set("bar1", map[string]interface{}{"status": "testing"})
	_, err := db.Copy("foo1", "bar1", "")
	if err != ErrConflict {
		t.Errorf(`db copy returns %v, want ErrConflict`, err)
	}
}

func TestCopyDocOverwrite(t *testing.T) {
	foo2 := map[string]interface{}{"status": "testing"}
	bar2 := map[string]interface{}{"status": "idle"}
	db.Set("foo2", foo2)
	db.Set("bar2", bar2)
	result, err := db.Copy("foo2", "bar2", bar2["_rev"].(string))
	if err != nil {
		t.Error(`db copy error`, err)
	}
	doc, _ := db.Get("bar2", nil)
	if result != doc["_rev"].(string) {
		t.Errorf("db copy returns %s want %s", result, doc["_rev"].(string))
	}
	if doc["status"].(string) != "testing" {
		t.Errorf("db copy status = %s, want testing", doc["status"].(string))
	}
}

// func TestChanges() {}
// func TestChangesConnUsable() {}
// func TestChangesHeartbeat() {}
// func TestPurge() {}
// func TestSecurity() {}
//
//
//
//
//
// func TestDatabaseSave(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   doc := map[string]interface{}{
//     "type": "Person",
//     "name": "John Doe",
//   }
//   id, rev := db.Save(doc)
//   if len(id) == 0 || len(rev) == 0 {
//     t.Error(`should return non-empty id and rev`)
//   }
//
//   doc["name"] = "Jason Statham"
//   id, rev = db.Save(doc)
//   if len(id) == 0 || len(rev) == 0 {
//     t.Error(`should return non-empty id and rev`)
//   }
//
//   doc["type"] = "Movie Star"
//   id, rev = db.Save(doc)
//   if len(id) == 0 || len(rev) == 0 {
//     t.Error(`should return non-empty id and rev`)
//   }
//
//   s.Delete("golang-tests")
// }
//
// func TestDatabaseAvailable(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   if !db.Available() {
//     t.Error(`database should be available`)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestDatabaseContains(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   doc := map[string]interface{}{
//     "type": "Person",
//     "name": "Jason Statham",
//   }
//   id, _ := db.Save(doc)
//   if len(id) <= 0 {
//     t.Error(`should return non-empty id`)
//   }
//   if !db.Contains(id) {
//     t.Error(`should contain id ` + id)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestDatabaseSetGetDelete(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   doc := map[string]interface{}{
//     "type": "Person",
//     "name": "Jason Statham",
//   }
//   if !db.Set("Mechanic", doc) {
//     t.Error(`set should return true`)
//   }
//   fetched := db.Get("Mechanic")
//   if fetched == nil {
//     t.Error(`get should return non-nil`)
//   }
//   if !db.Delete("Mechanic") {
//     t.Error(`delete should return true`)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestDatabaseDocIDsAndLen(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   doc := map[string]interface{}{
//     "type": "Person",
//     "name": "Jason Statham",
//   }
//
//   if !db.Set("Mechanic", doc) {
//     t.Error(`set should return true`)
//   }
//
//   ids := db.DocIDs()
//   if ids == nil {
//     t.Error(`should return slice of string`)
//   }
//
//   if len(ids) != 1 {
//     t.Error(`should return 1`)
//   }
//
//   if db.Len() != 1 {
//     t.Error(`Len() should return 1`)
//   }
//
//   s.Delete("golang-tests")
// }
//
// func TestDatabaseCommit(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   if !db.Commit() {
//     t.Error(`commit should be true`)
//   }
//   s.Delete("golang-tests")
// }
//
//
// func TestUpdateDocuments(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   docs := []map[string]interface{}{
//     {
//       "type": "Person",
//       "name": "Jason Statham",
//     },
//     {
//       "type": "Person",
//       "name": "Sylvester Stallone",
//     },
//     {
//       "type": "Person",
//       "name": "Arnold Schwarzenegger",
//     },
//     {
//       "type": "Person",
//       "name": "Sam Worthington",
//     },
//   }
//
//   idRevs, ok := db.UpdateDocuments(docs, nil)
//
//   if !ok {
//     t.Error(`update documents should return true`)
//   }
//
//   if len(idRevs) != len(docs) {
//     t.Error(`update documents should return id and revs in`, len(docs))
//   }
//
//   s.Delete("golang-tests")
// }
//
// func TestUserManagement(t *testing.T) {
//   id, rev := s.AddUser("foo", "secret", []string{"hero"})
//   if len(id) == 0 || len(rev) == 0 {
//     t.Error(`add user should return non-empty id and rev`)
//   }
//
//   token, ok := s.Login("foo", "secret")
//   if !ok {
//     t.Error(`login should return true`)
//   }
//
//   if !s.VerifyToken(token) {
//     t.Error(`token should be valid`, token)
//   }
//
//   if !s.Logout(token) {
//     t.Error(`logout should return true`)
//   }
//
//   if !s.RemoveUser("foo") {
//     t.Error(`remove user should return true`)
//   }
// }
//
// func TestGetSetRevsLimit(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   if !db.SetRevsLimit(10) {
//     t.Error(`set revs limit should return true`)
//   }
//   limit, ok := db.GetRevsLimit()
//   if !ok {
//     t.Error(`get revs limit should return true`)
//   }
//   if limit != 10 {
//     t.Error(`limit should be 10`)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestChanges(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   options := url.Values{
//     "style": []string{"all_docs"},
//   }
//   changes, ok := db.Changes(options)
//   if !ok {
//     t.Error(`changes should return true`)
//   }
//   if changes == nil {
//     t.Error(`changes should be non-nil`)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestCleanup(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   ok := db.Cleanup()
//   if !ok {
//     t.Error(`cleanup should return true`)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestCompact(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   ok := db.Compact()
//   if !ok {
//     t.Error(`compact should return true`)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestCopy(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   doc := map[string]interface{}{
//     "type": "Person",
//     "name": "Jason Statham",
//   }
//   src, _ := db.Save(doc)
//   dst := GenerateUUID()
//   _, ok := db.Copy(src, dst)
//   if !ok {
//     t.Error(`compact should return true`)
//   }
//   dstDoc := db.Get(dst)
//   if dstDoc == nil {
//     t.Error(`dstDoc should be non-nil`)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestPurge(t *testing.T) {
//   // db, _ := s.Create("golang-tests")
//
//   //TODO
//   // s.Delete("golang-tests")
// }
//
// func TestSecurity(t *testing.T) {
//   db, _ := s.Create("golang-tests")
//   secDoc, ok := db.GetSecurity()
//   if !ok {
//     t.Error(`get security should return true`)
//   }
//   if len(secDoc) > 0 {
//     t.Error(`secDoc should be empty`)
//   }
//   if !db.SetSecurity(map[string]interface{}{
//     "names": []string{"test"},
//     "roles": []string{},
//   }) {
//     t.Error(`set security should return true`)
//   }
//   s.Delete("golang-tests")
// }
