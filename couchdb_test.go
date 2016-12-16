package couchdb

import (
  // "fmt"
  "mime"
  "os"
  "path/filepath"
  "reflect"
  "strings"
  "testing"
)

var s *Server

func init() {
  s = NewServer("http://root:likejun@localhost:5984")
}

func TestNewServer(t *testing.T) {
  server := NewServer(DEFAULT_BASE_URL)
  if server == nil {
    t.Error(`server should be non-nil`)
  }
}

func TestVersion(t *testing.T) {
  version := s.Version()
  if reflect.ValueOf(version).Kind() != reflect.String {
    t.Error(`version should be string`)
  }
  if !strings.HasPrefix(version, "2") {
    t.Error(`version should be above 2`)
  }
}

func TestActiveTasks(t *testing.T) {
  jsonArr := s.ActiveTasks()
  if reflect.ValueOf(jsonArr).Kind() != reflect.Slice {
    t.Error(`jsonArr should be slice`)
  }
}

func TestDBs(t *testing.T) {
  dbs := s.DBs()
  kind := reflect.ValueOf(dbs).Kind()
  elemKind := reflect.TypeOf(dbs).Elem().Kind()
  if kind != reflect.Slice || elemKind != reflect.String {
    t.Error(`dbs shold be string slice`)
  }
}

func TestMembership(t *testing.T) {
  allNodes, clusterNodes := s.Membership()
  if allNodes == nil || clusterNodes == nil {
    t.Error(`unauthorized`)
  }

  kind := reflect.ValueOf(allNodes).Kind()
  elemKind := reflect.TypeOf(allNodes).Elem().Kind()

  if kind != reflect.Slice || elemKind != reflect.String {
    t.Error(`clusterNodes should be slice of string`)
  }

  kind = reflect.ValueOf(clusterNodes).Kind()
  elemKind = reflect.TypeOf(clusterNodes).Elem().Kind()

  if kind != reflect.Slice || elemKind != reflect.String {
    t.Error(`allNodes should be slice of string`)
  }
}

func TestReplicate(t *testing.T) {
  rsp := s.Replicate("db_a", "db_b", nil)
  if reflect.ValueOf(rsp).Kind() != reflect.Map {
    t.Error(`should return a map`)
  }
}

// func TestStats(t *testing.T) {
//   stats := s.Stats()
//   if reflect.ValueOf(stats).Kind() != reflect.Map {
//     t.Error(`should return a map`)
//   }
// }

func TestUUIDs(t *testing.T) {
  uuids := s.UUIDs(10)
  kind := reflect.ValueOf(uuids).Kind()
  elemKind := reflect.TypeOf(uuids).Elem().Kind()

  if kind != reflect.Slice || elemKind != reflect.String {
    t.Error(`should return slice of string`)
  }
}

func TestCreateDeleteDatabase(t *testing.T) {
  if _, ok := s.Create("golang-tests"); !ok {
    t.Error(`create db failed`)
  }

  if ok := s.Delete("golang-tests"); !ok {
    t.Error(`delete db failed`)
  }
}

func TestCreateDatabaseIllegal(t *testing.T) {
  if _, ok := s.Create("_db"); ok {
    t.Error(`create _db should not succeed`)
  }
}

func TestGetDatabase(t *testing.T) {
  _, ok := s.Create("golang-tests")
  if !ok {
    t.Error(`get db failed`)
  }
  s.Delete("golang-tests")
}

func TestGetNotExistDatabase(t *testing.T) {
  if db := s.GetDatabase("_not_exist"); db != nil {
    t.Error(`db should be nil`)
  }
}

func TestDatabaseName(t *testing.T) {
  s.Create("golang-tests")
  db := NewDatabase("http://root:likejun@localhost:5984/golang-tests")
  if (db == nil) {
    t.Error(`db should be non-nil`)
  }
  if (db.Name() != "golang-tests") {
    t.Error(`should return db name`)
  }
  s.Delete("golang-tests")
}

func TestDatabaseSave(t *testing.T) {
  db, _ := s.Create("golang-tests")
  doc := map[string]interface{}{
    "type": "Person",
    "name": "John Doe",
  }
  id, rev := db.Save(doc)
  if len(id) == 0 || len(rev) == 0 {
    t.Error(`should return non-empty id and rev`)
  }

  doc["name"] = "Jason Statham"
  id, rev = db.Save(doc)
  if len(id) == 0 || len(rev) == 0 {
    t.Error(`should return non-empty id and rev`)
  }

  doc["type"] = "Movie Star"
  id, rev = db.Save(doc)
  if len(id) == 0 || len(rev) == 0 {
    t.Error(`should return non-empty id and rev`)
  }

  s.Delete("golang-tests")
}

func TestDatabaseAvailable(t *testing.T) {
  db, _ := s.Create("golang-tests")
  if !db.Available() {
    t.Error(`database should be available`)
  }
  s.Delete("golang-tests")
}

func TestDatabaseContains(t *testing.T) {
  db, _ := s.Create("golang-tests")
  doc := map[string]interface{}{
    "type": "Person",
    "name": "Jason Statham",
  }
  id, _ := db.Save(doc)
  if len(id) <= 0 {
    t.Error(`should return non-empty id`)
  }
  if !db.Contains(id) {
    t.Error(`should contain id ` + id)
  }
  s.Delete("golang-tests")
}

func TestDatabaseSetGetDelete(t *testing.T) {
  db, _ := s.Create("golang-tests")
  doc := map[string]interface{}{
    "type": "Person",
    "name": "Jason Statham",
  }
  if !db.Set("Mechanic", doc) {
    t.Error(`set should return true`)
  }
  fetched := db.Get("Mechanic")
  if fetched == nil {
    t.Error(`get should return non-nil`)
  }
  if !db.Delete("Mechanic") {
    t.Error(`delete should return true`)
  }
  s.Delete("golang-tests")
}

func TestDatabaseDocIDsAndLen(t *testing.T) {
  db, _ := s.Create("golang-tests")
  doc := map[string]interface{}{
    "type": "Person",
    "name": "Jason Statham",
  }

  if !db.Set("Mechanic", doc) {
    t.Error(`set should return true`)
  }

  ids := db.DocIDs()
  if ids == nil {
    t.Error(`should return slice of string`)
  }

  if len(ids) != 1 {
    t.Error(`should return 1`)
  }

  if db.Len() != 1 {
    t.Error(`Len() should return 1`)
  }

  s.Delete("golang-tests")
}

func TestDatabaseCommit(t *testing.T) {
  db, _ := s.Create("golang-tests")
  if !db.Commit() {
    t.Error(`commit should be true`)
  }
  s.Delete("golang-tests")
}

func TestPutGetDeleteAttachment(t *testing.T) {
  content := "hello couch"
  db, _ := s.Create("golang-tests")
  tmpFileName := filepath.Join(os.TempDir(), "test.txt")
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

  doc := map[string]interface{}{
    "type": "Person",
    "name": "Jason Statham",
  }
  db.Set(GenerateUUID(), doc)
  if !db.PutAttachment(doc, tmpFile, mime.TypeByExtension(filepath.Ext(tmpFileName))) {
    t.Error(`put attachment should return true`)
  }

  data, ok := db.GetAttachment(doc["_id"].(string), "test.txt")
  if !ok {
    t.Error(`get attachment should return true`)
  }

  if string(data) != content {
    t.Error(`read data should be `, content)
  }

  if !db.DeleteAttachment(doc, tmpFileName) {
    t.Error(`delete attachment file failed`)
  }

  s.Delete("golang-tests")
}

func TestUpdateDocuments(t *testing.T) {
  db, _ := s.Create("golang-tests")
  docs := []map[string]interface{}{
    {
      "type": "Person",
      "name": "Jason Statham",
    },
    {
      "type": "Person",
      "name": "Sylvester Stallone",
    },
    {
      "type": "Person",
      "name": "Arnold Schwarzenegger",
    },
    {
      "type": "Person",
      "name": "Sam Worthington",
    },
  }

  idRevs, ok := db.UpdateDocuments(docs, nil)

  if !ok {
    t.Error(`update documents should return true`)
  }

  if len(idRevs) != len(docs) {
    t.Error(`update documents should return id and revs in`, len(docs))
  }

  s.Delete("golang-tests")
}

func TestUserManagement(t *testing.T) {
  id, rev := s.AddUser("foo", "secret", []string{"hero"})
  if len(id) == 0 || len(rev) == 0 {
    t.Error(`add user should return non-empty id and rev`)
  }

  token, ok := s.Login("foo", "secret")
  if !ok {
    t.Error(`login should return true`)
  }

  if !s.VerifyToken(token) {
    t.Error(`token should be valid`, token)
  }

  if !s.Logout(token) {
    t.Error(`logout should return true`)
  }

  if !s.RemoveUser("foo") {
    t.Error(`remove user should return true`)
  }
}

func TestGetSetRevsLimit(t *testing.T) {
  db, _ := s.Create("golang-tests")
  if !db.SetRevsLimit(10) {
    t.Error(`set revs limit should return true`)
  }
  limit, ok := db.GetRevsLimit()
  if !ok {
    t.Error(`get revs limit should return true`)
  }
  if limit != 10 {
    t.Error(`limit should be 10`)
  }
  s.Delete("golang-tests")
}
