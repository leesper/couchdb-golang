package couchdb

import
// "fmt"

(
	"net/url"
	"testing"
)

func TestNewDefaultDB(t *testing.T) {
	db, err := NewDatabase("golang-tests")
	if err != nil {
		t.Errorf("new default database error %v", err)
	}
	if db.Available() {
		t.Error(`db available`)
	}
}

func TestNewDB(t *testing.T) {
	s.Create("golang-tests")
	db, err := NewDatabase("http://root:likejun@localhost:5984/golang-tests")
	if err != nil {
		t.Error(`new database error`, err)
	}
	if !db.Available() {
		t.Error(`db not available`)
	}
	s.Delete("golang-tests")
}

func TestSaveNew(t *testing.T) {
	db, _ := s.Create("golang-tests")
	defer s.Delete("golang-tests")
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
	db, _ := s.Create("golang-tests")
	defer s.Delete("golang-tests")
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
	db, _ := s.Create("golang-tests")
	defer s.Delete("golang-tests")
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
	db, _ := s.Create("golang-tests")
	defer s.Delete("golang-tests")
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
	db, _ := s.Create("golang-tests")
	defer s.Delete("golang-tests")
	doc := map[string]interface{}{"_id": "foo"}

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

// func TestDatabaseExists() {}
// func TestDatabaseName() {}
// func TestCommit() {}
// func TestCreateLargeDoc() {}
// func TestDocIDQuoting() {}
// func TestDisallowNaN() {}
// func TestDisallowNilID() {}
// func TestDocRevs() {}
// func TestAttachmentCRUD() {}
// func TestAttachmentCRUDWithFiles() {}
// func TestAttachmentFromFS() {}
// func TestEmptyAttachment() {}
// func TestDefaultAttachment() {}
// func TestAttachmentNoFilename() {}
// func TestJSONAttachment() {}
// func TestIncludeDocs() {}
// // TODO adding new apis about mango query engine
// func TestQueryMultiGet() {}
// func TestBulkUpdateConflict() {}
// func TestBulkUpdateAllOrNothing() {}
// func TestBulkUpdateBadDoc() {}
// func TestCopyDocConflict() {}
// func TestCopyDocOverwrite() {}
// func TestCopyDocSrcObj() {}
// func TestCopyDocDestObjNoRev() {}
// func TestCopyDocSrcDictLike() {}
// func TestCopyDocDestDictLike() {}
// func TestCopyDocSrcBadDoc() {}
// func TestCopyDocDestBadDoc() {}
// func TestChanges() {}
// func TestChangesConnUsable() {}
// func TestChangesHeartbeat() {}
// func TestPurge() {}
// func TestSecurity() {}
//
//
//
//
// func TestDatabaseName(t *testing.T) {
//   s.Create("golang-tests")
//   db := NewDatabase("http://root:likejun@localhost:5984/golang-tests")
//   if (db == nil) {
//     t.Error(`db should be non-nil`)
//   }
//   if (db.Name() != "golang-tests") {
//     t.Error(`should return db name`)
//   }
//   s.Delete("golang-tests")
// }
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
// func TestPutGetDeleteAttachment(t *testing.T) {
//   content := "hello couch"
//   db, _ := s.Create("golang-tests")
//   tmpFileName := filepath.Join(os.TempDir(), "test.txt")
//   tmpFile, err := os.Create(tmpFileName)
//   if err != nil {
//     t.Error(`create file error`, err)
//   }
//   _, err = tmpFile.Write([]byte(content))
//   if err != nil {
//     t.Error(`write file error`, err)
//   }
//   tmpFile.Close()
//
//   tmpFile, err = os.Open(tmpFileName)
//   if err != nil {
//     t.Error(`open file error`, err)
//   }
//   defer tmpFile.Close()
//
//   doc := map[string]interface{}{
//     "type": "Person",
//     "name": "Jason Statham",
//   }
//   db.Set(GenerateUUID(), doc)
//   if !db.PutAttachment(doc, tmpFile, mime.TypeByExtension(filepath.Ext(tmpFileName))) {
//     t.Error(`put attachment should return true`)
//   }
//
//   data, ok := db.GetAttachment(doc["_id"].(string), "test.txt")
//   if !ok {
//     t.Error(`get attachment should return true`)
//   }
//
//   if string(data) != content {
//     t.Error(`read data should be `, content)
//   }
//
//   if !db.DeleteAttachment(doc, tmpFileName) {
//     t.Error(`delete attachment file failed`)
//   }
//
//   s.Delete("golang-tests")
// }
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
