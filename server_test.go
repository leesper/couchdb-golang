package couchdb

import (
  // "fmt"
  // "reflect"
  // "strings"
  "testing"
)

var s *Server

func init() {
  s, _ = NewServer("http://root:likejun@localhost:5984")
}

func TestNewServer(t *testing.T) {
  _, err := NewServer(DEFAULT_BASE_URL)
  if err != nil {
    t.Error(`new server error`, err)
  }
  // _, err = server.Config()
  // if err != nil {
  //   t.Error(`server config error`, err)
  // }
}

func TestNewServerNoFullCommit(t *testing.T) {
  _, err := NewServerNoFullCommit(DEFAULT_BASE_URL)
  if err != nil {
    t.Error(`new server full commit error`, err)
  }
  // _, err = server.Config()
  // if err != nil {
  //   t.Error(`server config error`, err)
  // }
}

func TestServerExists(t *testing.T) {
  _, err := NewServer("http://localhost:9999")
  if err != nil {
    t.Error(`new server error`, err)
  }
  
}

// func TestServerString(t *testing.T) {
//   server, err := NewServer(DEFAULT_BASE_URL)
//   if err != nil {
//     t.Error(`new server error`, err)
//   }
//   fmt.Println(server)
// }
//
// func TestServerVars(t *testing.T) {
//   version := s.Version()
//   if reflect.ValueOf(version).Kind() != reflect.String {
//     t.Error(`version not of string type`)
//   }
//
//   config := s.Config()
//   if reflect.ValueOf(config).Kind() != reflect.Map {
//     t.Error(`config not of map type`)
//   }
//
//   tasks := s.ActiveTasks()
//   if reflect.ValueOf(tasks).Kind() != reflect.Slice {
//     t.Error(`tasks not of slice type`)
//   }
// }
//
// func TestServerStats(t *testing.T) {
//   stats := s.Stats()
//   if reflect.ValueOf(stats).Kind() != reflect.Map {
//     t.Error(`stats not of map type`)
//   }
//   stats = s.Stats("httpd/requests")
//   if reflect.ValueOf(stats).Kind() != reflect.Map {
//     t.Error(`httpd/requests stats not of map type`)
//   }
//   ok := len(stats) == 1 && len(stats["httpd"]) == 1
//   if !ok {
//     t.Errorf("len(stats) = %d want 1, len(stats[httpd]) = %d, want 1", len(stats), len(stats["httpd"]))
//   }
// }
//
// func TestGetDBMissing(t *testing.T) {
//   _, err := s.Get("golang-tests")
//   if err != ErrNotFound {
//     t.Errorf("err = %v want ErrNotFound", err)
//   }
// }
//
// func TestGetDB(t *testing.T) {
//   _, err := s.Create("golang-tests")
//   if err != nil {
//     t.Error(`create db error`, err)
//   }
//   _, err := s.Get("golang-tests")
//   if err != nil {
//     t.Error(`get db error`, err)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestCreateDBConflict(t *testing.T) {
//   conflictDBName := "golang-tests"
//   db, _ := s.Create(conflictDBName)
//   defer s.Delete(conflictDBName)
//   if _, err := s.Create(conflictDBName); err != ErrPreconditionFailed {
//     t.Errorf("err = %v want ErrPreconditionFailed", err)
//   }
// }
//
// func TestCreateDB(t *testing.T) {
//   _, ok := s.Create("golang-tests")
//   if !ok {
//     t.Error(`get db failed`)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestCreateDatabaseIllegal(t *testing.T) {
//   if _, ok := s.Create("_db"); ok {
//     t.Error(`create _db should not succeed`)
//   }
// }
//
// func TestDeleteDB(t *testing.T) {
//   dbName := "golang-tests"
//   s.Create(dbName)
//   if !s.Contains(dbName) {
//     t.Error(`server not contains`, dbName)
//   }
//   s.Delete(dbName)
//   if s.Contains(dbName) {
//     t.Error(`server contains`, dbName)
//   }
// }
//
// func TestDeleteDBMissing(t *testing.T) {
//   dbName := "golang-tests"
//   err := s.Delete(dbName)
//   if err != ErrNotFound {
//     t.Errorf("err = %v want ErrNotFound", err)
//   }
// }
//
// func TestReplicate(t *testing.T) {
//   aName := "dba"
//   dba, _ = s.Create(aName)
//   defer s.Delete(aName)
//
//   bName := "dbb"
//   dbb, _ = s.Create(bName)
//   defer s.Delete(bName)
//
//   id, _ := dba.Save(map[string]interface{}{"test": "a"})
//   result, err := s.Replicate(aName, bName, nil)
//   if v, ok := result["ok"]; !(ok && v.(bool)) {
//     t.Error(`result should be ok`)
//   }
//   doc, err := dbb.Get(id)
//   if err != nil {
//     t.Errorf("db %s get doc %s error %v", bName, id, err)
//   }
//   if v, ok := doc["test"]; ok {
//     if "a" != v.(string) {
//       t.Error(`doc[test] should be a, found`, v.(string))
//     }
//   }
//
//   doc["test"] = "b"
//   dbb.Update([]map[string]interface{}{doc})
//   s.Replicate(bName, aName, nil)
//
//   docA, err := dba.Get(id)
//   if err != nil {
//     t.Errorf("db %s get doc %s error %v", aName, id, err)
//   }
//   if v, ok := docA["test"]; ok {
//     if "b" != v.(string) {
//       t.Error(`docA[test] should be b, found`, v.(string))
//     }
//   }
//
//   docB, err := dbb.Get(id)
//   if err != nil {
//     t.Errorf("db %s get doc %s error %v", bName, id, err)
//   }
//   if v, ok := docB["test"]; ok {
//     if "b" != v.(string) {
//       t.Error(`docB[test] should be b, found`, v.(string))
//     }
//   }
// }
//
// func TestReplicateContinuous(t *testing.T) {
//   aName, bName := "dba", "dbb"
//   s.Create(aName)
//   defer s.Delete(aName)
//
//   s.Create(bName)
//   defer s.Delete(bName)
//
//   result = s.Replicate(aName, bName, url.Values{"continuous": []string{"true"}})
//   if v, ok := result["ok"]; !(ok && v.(bool)) {
//     t.Error(`result should be ok`)
//   }
// }
//
// func TestDBs(t *testing.T) {
//   aName, bName := "dba", "dbb"
//   s.Create(aName)
//   defer s.Delete(aName)
//
//   s.Create(bName)
//   defer s.Delete(bName)
//
//   dbs := s.DBs()
//   var aExist, bExist bool
//   for _, v := range dbs {
//     if v == aName {
//       aExist = true
//     } else if v == bName {
//       bExist = true
//     }
//   }
//
//   if !aExist {
//     t.Errorf("db %s not existed in dbs", aName)
//   }
//
//   if !bExist {
//     t.Errorf("db %s not existed in dbs", bName)
//   }
// }
//
// func TestLen(t *testing.T) {
//   aName, bName := "dba", "dbb"
//   s.Create(aName)
//   defer s.Delete(aName)
//
//   s.Create(bName)
//   defer s.Delete(bName)
//
//   if s.Len() < 2 {
//     t.Error("server len should be >= 2")
//   }
// }
//
// func TestUUIDs(t *testing.T) {
//   uuids := s.UUIDs(10)
//   if reflect.ValueOf(uuids).Kind() != reflect.Slice {
//     t.Error(`uuids should be of type slice`)
//   }
//   if len(uuids) != 10 {
//     t.Error(`uuids should be of length 10, not`, len(uuids))
//   }
// }
//
// func TestBasicAuth(t *testing.T) {
//   server, _ := NewServer("http://root:password@localhost:5984/")
//   _, err := server.Create("golang-tests")
//   if err != ErrUnauthorized {
//     t.Errorf("err = %v want ErrUnauthorized")
//   }
// }
//
// func TestUserManagement(t *testing.T) {
//   s.AddUser("foo", "secret", []string{"hero"})
//   token = server.Login("foo", "secret")
//   if len(token) == 0 {
//     t.Error(`server login error, token empty`)
//   }
//   if !server.VerifyToken(token) {
//     t.Error("server verify token false")
//   }
//   if !server.Logout(token) {
//     t.Error("server logout false")
//   }
//   server.RemoveUser("foo")
// }
//
// func TestMembership(t *testing.T) {
//   allNodes, clusterNodes := s.Membership()
//   if allNodes == nil || clusterNodes == nil {
//     t.Error(`unauthorized`)
//   }
//
//   kind := reflect.ValueOf(allNodes).Kind()
//   elemKind := reflect.TypeOf(allNodes).Elem().Kind()
//
//   if kind != reflect.Slice || elemKind != reflect.String {
//     t.Error(`clusterNodes should be slice of string`)
//   }
//
//   kind = reflect.ValueOf(clusterNodes).Kind()
//   elemKind = reflect.TypeOf(clusterNodes).Elem().Kind()
//
//   if kind != reflect.Slice || elemKind != reflect.String {
//     t.Error(`allNodes should be slice of string`)
//   }
// }
