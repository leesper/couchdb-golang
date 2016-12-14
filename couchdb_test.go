package couchdb

import (
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

func TestCreateDatabase(t *testing.T) {
  if _, ok := s.Create("hello_couch"); !ok {
    t.Error(`create db failed`)
  }
}

func TestCreateDatabaseIllegal(t *testing.T) {
  if _, ok := s.Create("_db"); ok {
    t.Error(`create _db should not succeed`)
  }
}

func TestDeleteDatabase(t *testing.T) {
  if ok := s.Delete("hello_couch"); !ok {
    t.Error(`delete db failed`)
  }
}

func TestGetDatabase(t *testing.T) {
  _, ok := s.Create("hello_couch")
  if !ok {
    t.Error(`get db failed`)
  }
  s.Delete("hello_couch")
}

func TestGetNotExistDatabase(t *testing.T) {
  if db := s.GetDatabase("_not_exist"); db != nil {
    t.Error(`db should be nil`)
  }
}

func TestDatabaseName(t *testing.T) {
  s.Create("hello_couch")
  db := NewDatabase("http://root:likejun@localhost:5984/hello_couch")
  if (db == nil) {
    t.Error(`db should be non-nil`)
  }
  if (db.Name() != "hello_couch") {
    t.Error(`should return db name`)
  }
  s.Delete("hello_couch")
}
