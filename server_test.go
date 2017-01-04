package couchdb

import (
	"os"
	"testing"
)

var (
	s       *Server
	db      *Database
	movieDB *Database
	movies  = []map[string]interface{}{
		{
			"_id":     "976059",
			"_rev":    "1-0ac57ca6eeb3dfc5f7a9196010385e7d",
			"title":   "Spacecataz",
			"year":    2004,
			"rating":  nil,
			"runtime": "10 min",
			"genre": []string{
				"Animation",
				"Short",
				"Comedy",
				"Sci-Fi",
			},
			"director": "Dave Willis",
			"writer": []string{
				"Matt Maiellaro",
				"Dave Willis",
			},
			"cast": []string{
				"Dave Willis",
				"Matt Maiellaro",
				"Andy Merrill",
				"Mike Schatz",
			},
			"poster": nil,
			"imdb": map[string]interface{}{
				"rating": 8,
				"votes":  130,
				"id":     "tt0976059",
			},
		},
	}
)

func TestMain(m *testing.M) {
	// setup
	var err error
	s, err = NewServer("http://root:likejun@localhost:5984")
	if err != nil {
		os.Exit(1)
	}
	s.Delete("golang-tests")
	db, err = s.Create("golang-tests")
	if err != nil {
		os.Exit(2)
	}

	s.Delete("golang-movies")
	movieDB, err = s.Create("golang-movies")
	if err != nil {
		os.Exit(2)
	}

	// run all the tests
	code := m.Run()

	// tear down
	s.Delete("golang-tests")
	s.Delete("golang-movies")
	os.Exit(code)
}

//
// func TestNewServer(t *testing.T) {
// 	server, err := NewServer("http://root:likejun@localhost:5984")
// 	if err != nil {
// 		t.Error(`new server error`, err)
// 	}
// 	_, err = server.Version()
// 	if err != nil {
// 		t.Error(`server version error`, err)
// 	}
// }
//
// func TestNewServerNoFullCommit(t *testing.T) {
// 	server, err := NewServerNoFullCommit("http://root:likejun@localhost:5984")
// 	if err != nil {
// 		t.Error(`new server full commit error`, err)
// 	}
// 	_, err = server.Version()
// 	if err != nil {
// 		t.Error(`server version error`, err)
// 	}
// }
//
// func TestServerExists(t *testing.T) {
// 	server, err := NewServer("http://localhost:9999")
// 	if err != nil {
// 		t.Error(`new server error`, err)
// 	}
//
// 	_, err = server.Version()
// 	if err == nil {
// 		t.Error(`server version ok`)
// 	}
// }
//
// func TestServerConfig(t *testing.T) {
// 	config, err := s.Config("couchdb@localhost")
// 	if err != nil {
// 		t.Error(`server config error`, err)
// 	}
// 	if reflect.ValueOf(config).Kind() != reflect.Map {
// 		t.Error(`config not of type map`)
// 	}
// }
//
// func TestServerString(t *testing.T) {
// 	server, err := NewServer(DEFAULT_BASE_URL)
// 	if err != nil {
// 		t.Error(`new server error`, err)
// 	}
// 	if server.String() != "Server http://localhost:5984" {
// 		t.Error(`server name invalid want "Server http://localhost:5984"`)
// 	}
// }
//
// func TestServerVars(t *testing.T) {
// 	version, err := s.Version()
// 	if err != nil {
// 		t.Error(`server version error`, err)
// 	}
// 	if reflect.ValueOf(version).Kind() != reflect.String {
// 		t.Error(`version not of string type`)
// 	}
//
// 	tasks, err := s.ActiveTasks()
// 	if reflect.ValueOf(tasks).Kind() != reflect.Slice {
// 		t.Error(`tasks not of slice type`)
// 	}
// }
//
// func TestServerStats(t *testing.T) {
// 	stats, err := s.Stats("couchdb@localhost", "")
// 	if err != nil {
// 		t.Error(`server stats error`, err)
// 	}
// 	if reflect.ValueOf(stats).Kind() != reflect.Map {
// 		t.Error(`stats not of map type`)
// 	}
// 	stats, err = s.Stats("couchdb@localhost", "couchdb")
// 	if err != nil {
// 		t.Error(`server stats httpd/requests error`, err)
// 	}
// 	if reflect.ValueOf(stats).Kind() != reflect.Map {
// 		t.Error(`httpd/requests stats not of map type`)
// 	}
// }
//
// func TestDBs(t *testing.T) {
// 	aName, bName := "dba", "dbb"
// 	s.Create(aName)
// 	defer s.Delete(aName)
//
// 	s.Create(bName)
// 	defer s.Delete(bName)
//
// 	dbs, err := s.DBs()
// 	if err != nil {
// 		t.Error(`server DBs error`, err)
// 	}
// 	var aExist, bExist bool
// 	for _, v := range dbs {
// 		if v == aName {
// 			aExist = true
// 		} else if v == bName {
// 			bExist = true
// 		}
// 	}
//
// 	if !aExist {
// 		t.Errorf("db %s not existed in dbs", aName)
// 	}
//
// 	if !bExist {
// 		t.Errorf("db %s not existed in dbs", bName)
// 	}
// }
//
// func TestLen(t *testing.T) {
// 	aName, bName := "dba", "dbb"
// 	s.Create(aName)
// 	defer s.Delete(aName)
// 	s.Create(bName)
// 	defer s.Delete(bName)
//
// 	len, err := s.Len()
// 	if err != nil {
// 		t.Error(`server len error`, err)
// 	}
// 	if len < 2 {
// 		t.Error("server len should be >= 2")
// 	}
// }
//
// func TestGetDBMissing(t *testing.T) {
// 	_, err := s.Get("golang-missing")
// 	if err != ErrNotFound {
// 		t.Errorf("err = %v want ErrNotFound", err)
// 	}
// }
//
// func TestGetDB(t *testing.T) {
// 	_, err := s.Get("golang-tests")
// 	if err != nil {
// 		t.Error(`get db error`, err)
// 	}
// }
//
// func TestCreateDBConflict(t *testing.T) {
// 	conflictDBName := "golang-conflict"
// 	_, err := s.Create(conflictDBName)
// 	if err != nil {
// 		t.Error(`server create error`, err)
// 	}
// 	// defer s.Delete(conflictDBName)
// 	if !s.Contains(conflictDBName) {
// 		t.Error(`server not contains`, conflictDBName)
// 	}
// 	if _, err = s.Create(conflictDBName); err != ErrPreconditionFailed {
// 		t.Errorf("err = %v want ErrPreconditionFailed", err)
// 	}
// 	s.Delete(conflictDBName)
// }
//
// func TestCreateDB(t *testing.T) {
// 	_, err := s.Create("golang-create")
// 	if err != nil {
// 		t.Error(`get db failed`)
// 	}
// 	s.Delete("golang-create")
// }
//
// func TestCreateDBIllegal(t *testing.T) {
// 	if _, err := s.Create("_db"); err == nil {
// 		t.Error(`create illegal _db ok`)
// 	}
// }
//
// func TestDeleteDB(t *testing.T) {
// 	dbName := "golang-delete"
// 	s.Create(dbName)
// 	if !s.Contains(dbName) {
// 		t.Error(`server not contains`, dbName)
// 	}
// 	s.Delete(dbName)
// 	if s.Contains(dbName) {
// 		t.Error(`server contains`, dbName)
// 	}
// }
//
// func TestDeleteDBMissing(t *testing.T) {
// 	dbName := "golang-missing"
// 	err := s.Delete(dbName)
// 	if err != ErrNotFound {
// 		t.Errorf("err = %v want ErrNotFound", err)
// 	}
// }
//
// func TestReplicate(t *testing.T) {
// 	aName := "dba"
// 	dba, _ := s.Create(aName)
// 	defer s.Delete(aName)
//
// 	bName := "dbb"
// 	dbb, _ := s.Create(bName)
// 	defer s.Delete(bName)
//
// 	id, _, err := dba.Save(map[string]interface{}{"test": "a"}, nil)
// 	if err != nil {
// 		t.Error(`dba save error`, err)
// 	}
// 	result, err := s.Replicate(aName, bName, nil)
// 	if v, ok := result["ok"]; !(ok && v.(bool)) {
// 		t.Error(`result should be ok`)
// 	}
// 	doc, err := dbb.Get(id, nil)
// 	if err != nil {
// 		t.Errorf("db %s get doc %s error %v", bName, id, err)
// 	}
// 	if v, ok := doc["test"]; ok {
// 		if "a" != v.(string) {
// 			t.Error(`doc[test] should be a, found`, v.(string))
// 		}
// 	}
//
// 	doc["test"] = "b"
// 	dbb.Update([]map[string]interface{}{doc}, nil)
// 	result, err = s.Replicate(bName, aName, nil)
// 	if err != nil {
// 		t.Error(`server replicate error`, err)
// 	}
// 	if reflect.ValueOf(result).Kind() != reflect.Map {
// 		t.Error(`server replicate return non-map result`)
// 	}
//
// 	docA, err := dba.Get(id, nil)
// 	if err != nil {
// 		t.Errorf("db %s get doc %s error %v", aName, id, err)
// 	}
// 	if v, ok := docA["test"]; ok {
// 		if "b" != v.(string) {
// 			t.Error(`docA[test] should be b, found`, v.(string))
// 		}
// 	}
//
// 	docB, err := dbb.Get(id, nil)
// 	if err != nil {
// 		t.Errorf("db %s get doc %s error %v", bName, id, err)
// 	}
// 	if v, ok := docB["test"]; ok {
// 		if "b" != v.(string) {
// 			t.Error(`docB[test] should be b, found`, v.(string))
// 		}
// 	}
// }
//
// func TestReplicateContinuous(t *testing.T) {
// 	aName, bName := "dba", "dbb"
// 	s.Create(aName)
// 	defer s.Delete(aName)
//
// 	s.Create(bName)
// 	defer s.Delete(bName)
//
// 	result, err := s.Replicate(aName, bName, map[string]interface{}{"continuous": true})
// 	if err != nil {
// 		t.Error(`server replicate error`, err)
// 	}
//
// 	if reflect.ValueOf(result).Kind() != reflect.Map {
// 		t.Error(`server replicate return non-map result`)
// 	}
//
// 	if v, ok := result["ok"]; !(ok && v.(bool)) {
// 		t.Error(`result should be ok`)
// 	}
// }
//
// func TestMembership(t *testing.T) {
// 	allNodes, clusterNodes, err := s.Membership()
// 	if err != nil {
// 		t.Error(`server membership error`, err)
// 	}
//
// 	kind := reflect.ValueOf(allNodes).Kind()
// 	elemKind := reflect.TypeOf(allNodes).Elem().Kind()
//
// 	if kind != reflect.Slice || elemKind != reflect.String {
// 		t.Error(`clusterNodes should be slice of string`)
// 	}
//
// 	kind = reflect.ValueOf(clusterNodes).Kind()
// 	elemKind = reflect.TypeOf(clusterNodes).Elem().Kind()
//
// 	if kind != reflect.Slice || elemKind != reflect.String {
// 		t.Error(`allNodes should be slice of string`)
// 	}
// }
//
// func TestUUIDs(t *testing.T) {
// 	uuids, err := s.UUIDs(10)
// 	if err != nil {
// 		t.Error(`server uuids error`, err)
// 	}
// 	if reflect.ValueOf(uuids).Kind() != reflect.Slice {
// 		t.Error(`server uuids should be of type slice`)
// 	}
// 	if len(uuids) != 10 {
// 		t.Error(`server uuids should be of length 10, not`, len(uuids))
// 	}
// }
//
// func TestBasicAuth(t *testing.T) {
// 	server, _ := NewServer("http://root:password@localhost:5984/")
// 	_, err := server.Create("golang-auth")
// 	if err != ErrUnauthorized {
// 		t.Errorf("err = %v want ErrUnauthorized", err)
// 	}
// }
//
// func TestUserManagement(t *testing.T) {
// 	user := "foo"
// 	password := "secret"
// 	roles := []string{"hero"}
// 	s.AddUser(user, password, roles)
//
// 	token, err := s.Login(user, password)
// 	if err != nil {
// 		t.Errorf("server add user %s password %s roles %v error %v", user, password, roles)
// 	}
//
// 	if err = s.VerifyToken(token); err != nil {
// 		t.Error("server verify token error", err)
// 	}
//
// 	if err = s.Logout(token); err != nil {
// 		t.Error("server logout error", err)
// 	}
//
// 	if err = s.RemoveUser("foo"); err != nil {
// 		t.Error("server remove user error", err)
// 	}
// }
