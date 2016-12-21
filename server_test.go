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

// func TestNewServer(t *testing.T) {
//   server, err := NewServer("http://root:likejun@localhost:5984")
//   if err != nil {
//     t.Error(`new server error`, err)
//   }
//   _, err = server.Version()
//   if err != nil {
//     t.Error(`server version error`, err)
//   }
// }
//
// func TestNewServerNoFullCommit(t *testing.T) {
//   server, err := NewServerNoFullCommit("http://root:likejun@localhost:5984")
//   if err != nil {
//     t.Error(`new server full commit error`, err)
//   }
//   _, err = server.Version()
//   if err != nil {
//     t.Error(`server version error`, err)
//   }
// }
//
// func TestServerExists(t *testing.T) {
//   server, err := NewServer("http://localhost:9999")
//   if err != nil {
//     t.Error(`new server error`, err)
//   }
//
//   _, err = server.Version()
//   if err == nil {
//     t.Error(`server version ok`)
//   }
// }
//
// func TestServerConfig(t *testing.T) {
//   config, err := s.Config("couchdb@localhost")
//   if err != nil {
//     t.Error(`server config error`, err)
//   }
//   if reflect.ValueOf(config).Kind() != reflect.Map {
//     t.Error(`config not of type map`)
//   }
// }
//
// func TestServerString(t *testing.T) {
//   server, err := NewServer(DEFAULT_BASE_URL)
//   if err != nil {
//     t.Error(`new server error`, err)
//   }
//   if server.String() != "Server http://localhost:5984" {
//     t.Error(`server name invalid want "Server http://localhost:5984"`)
//   }
// }
//
// func TestServerVars(t *testing.T) {
//   version, err := s.Version()
//   if err != nil {
//     t.Error(`server version error`, err)
//   }
//   if reflect.ValueOf(version).Kind() != reflect.String {
//     t.Error(`version not of string type`)
//   }
//
//   tasks, err := s.ActiveTasks()
//   if reflect.ValueOf(tasks).Kind() != reflect.Slice {
//     t.Error(`tasks not of slice type`)
//   }
// }
//
// func TestServerStats(t *testing.T) {
//   stats, err := s.Stats("couchdb@localhost", "")
//   if err != nil {
//     t.Error(`server stats error`, err)
//   }
//   if reflect.ValueOf(stats).Kind() != reflect.Map {
//     t.Error(`stats not of map type`)
//   }
//   stats, err = s.Stats("couchdb@localhost", "couchdb")
//   if err != nil {
//     t.Error(`server stats httpd/requests error`, err)
//   }
//   if reflect.ValueOf(stats).Kind() != reflect.Map {
//     t.Error(`httpd/requests stats not of map type`)
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
//   dbs, err := s.DBs()
//   if err != nil {
//     t.Error(`server DBs error`, err)
//   }
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
//   s.Create(bName)
//   defer s.Delete(bName)
//
//   len, err := s.Len()
//   if err != nil {
//     t.Error(`server len error`, err)
//   }
//   if len < 2 {
//     t.Error("server len should be >= 2")
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
//   _, err = s.Get("golang-tests")
//   if err != nil {
//     t.Error(`get db error`, err)
//   }
//   s.Delete("golang-tests")
// }

func TestCreateDBConflict(t *testing.T) {
  conflictDBName := "golang-tests"
  _, err := s.Create(conflictDBName)
  if err != nil {
    t.Error(`server create error`, err)
  }
  // defer s.Delete(conflictDBName)
  if !s.Contains(conflictDBName) {
    t.Error(`server not contains`, conflictDBName)
  }
  if _, err = s.Create(conflictDBName); err != ErrPreconditionFailed {
    t.Errorf("err = %v want ErrPreconditionFailed", err)
  }
  s.Delete(conflictDBName)
}

// func TestCreateDB(t *testing.T) {
//   _, err := s.Create("golang-tests")
//   if err != nil {
//     t.Error(`get db failed`)
//   }
//   s.Delete("golang-tests")
// }
//
// func TestCreateDBIllegal(t *testing.T) {
//   if _, err := s.Create("_db"); err == nil {
//     t.Error(`create illegal _db ok`)
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
