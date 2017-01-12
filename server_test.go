package couchdb

import (
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	server   *Server
	testsDB  *Database
	movieDB  *Database
	designDB *Database
	movies   = []map[string]interface{}{
		{
			"_id":     "976059",
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
		{
			"_id":     "976197",
			"title":   "American Psyche",
			"year":    2007,
			"rating":  nil,
			"runtime": "55 min",
			"genre": []string{
				"Documentary",
			},
			"director": "Paul van den Boom",
			"writer": []string{
				"Paul van den Boom (creator)",
				"Franois Le Goarant de Tromelin (creator)",
			},
			"cast": []string{
				"Katherine J. Eakin",
				"Mahnaz M. Shabbir",
				"Rene Doria",
				"Peter Koper",
			},
			"poster": "http://ia.media-imdb.com/images/M/MV5BMTM3NTg5NDE2N15BMl5BanBnXkFtZTcwODI0MjM1MQ@@._V1_SX300.jpg",
			"imdb": map[string]interface{}{
				"rating": 8.2,
				"votes":  77,
				"id":     "tt0976197",
			},
		},
		{
			"_id":     "976221",
			"title":   "Voliminal: Inside the Nine",
			"year":    2006,
			"rating":  nil,
			"runtime": "N/A",
			"genre": []string{
				"Documentary",
			},
			"director": "Shawn Crahan",
			"writer":   nil,
			"cast": []string{
				"Shawn Crahan",
				"Chris Fehn",
				"Paul Gray",
				"Craig Jones",
			},
			"poster": nil,
			"imdb": map[string]interface{}{
				"rating": 8.1,
				"votes":  125,
				"id":     "tt0976221",
			},
		},
		{
			"_id":     "97628",
			"title":   "The Johnstown Flood",
			"year":    1989,
			"rating":  nil,
			"runtime": "26 min",
			"genre": []string{
				"Documentary",
				"Short",
			},
			"director": "Charles Guggenheim",
			"writer": []string{
				"Charles Guggenheim",
			},
			"cast": []string{
				"Len Cariou",
				"Elam Bender",
				"Randy Bender",
				"Clarita Berger",
			},
			"poster": "http://ia.media-imdb.com/images/M/MV5BMTc2NTc3MzQ5MF5BMl5BanBnXkFtZTcwMjU5ODkwNg@@._V1_SX300.jpg",
			"imdb": map[string]interface{}{
				"rating": 7.8,
				"votes":  75,
				"id":     "tt0097628",
			},
		},
		{
			"_id":     "97661",
			"title":   "Gundam 0080: A War in the Pocket",
			"year":    1989,
			"rating":  "NOT RATED",
			"runtime": "N/A",
			"genre": []string{
				"Animation",
				"Action",
				"Drama",
			},
			"director": "N/A",
			"writer":   nil,
			"cast": []string{
				"Daisuke Namikawa",
				"Kji Tsujitani",
				"Megumi Hayashibara",
				"Brianne Brozey",
			},
			"poster": "http://ia.media-imdb.com/images/M/MV5BMTk3NjU2ODQ1Ml5BMl5BanBnXkFtZTcwNzY4MzYyMQ@@._V1_SX300.jpg",
			"imdb": map[string]interface{}{
				"rating": 8,
				"votes":  475,
				"id":     "tt0097661",
			},
		},
		{
			"_id":     "97690",
			"title":   "Kuduz",
			"year":    1989,
			"rating":  nil,
			"runtime": "N/A",
			"genre": []string{
				"Drama",
			},
			"director": "Ademir Kenovic",
			"writer": []string{
				"Ademir Kenovic",
				"Abdulah Sidran",
			},
			"cast": []string{
				"Slobodan Custic",
				"Snezana Bogdanovic",
				"Bozidar Bunjevac",
				"Branko Djuric",
			},
			"poster": nil,
			"imdb": map[string]interface{}{
				"rating": 8.1,
				"votes":  342,
				"id":     "tt0097690",
			},
		},
		{
			"_id":     "977224",
			"title":   "Mere Oblivion",
			"year":    2007,
			"rating":  nil,
			"runtime": "N/A",
			"genre": []string{
				"Short",
				"Comedy",
			},
			"director": "Burleigh Smith",
			"writer": []string{
				"Burleigh Smith",
			},
			"cast": []string{
				"Burleigh Smith",
				"Elizabeth Caiacob",
				"Michael Su",
				"Kate Ritchie",
			},
			"poster": nil,
			"imdb": map[string]interface{}{
				"rating": 7.9,
				"votes":  284,
				"id":     "tt0977224",
			},
		},
		{
			"_id":     "97727",
			"title":   "A legnyanya",
			"year":    1989,
			"rating":  nil,
			"runtime": "80 min",
			"genre": []string{
				"Comedy",
			},
			"director": "Dezs Garas",
			"writer": []string{
				"Dezs Garas (screenplay)",
				"Gyrgy Schwajda",
			},
			"cast": []string{
				"Ferenc Kllai",
				"Kroly Eperjes",
				"Judit Pogny",
				"Dezs Garas",
			},
			"poster": nil,
			"imdb": map[string]interface{}{
				"rating": 7.8,
				"votes":  786,
				"id":     "tt0097727",
			},
		},
		{
			"_id":     "97757",
			"title":   "The Little Mermaid",
			"year":    1989,
			"rating":  "G",
			"runtime": "83 min",
			"genre": []string{
				"Animation",
				"Family",
				"Fantasy",
			},
			"director": "Ron Clements, John Musker",
			"writer": []string{
				"John Musker",
				"Ron Clements",
				"Hans Christian Andersen (fairy tale)",
				"Howard Ashman (additional dialogue)",
				"Gerrit Graham (additional dialogue)",
				"Sam Graham (additional dialogue)",
				"Chris Hubbell (additional dialogue)",
			},
			"cast": []string{
				"Rene Auberjonois",
				"Christopher Daniel Barnes",
				"Jodi Benson",
				"Pat Carroll",
			},
			"poster": "http://ia.media-imdb.com/images/M/MV5BNTAxMzY0MjI1Nl5BMl5BanBnXkFtZTgwMTU2NTYxMTE@._V1_SX300.jpg",
			"imdb": map[string]interface{}{
				"rating": 7.6,
				"votes":  138,
				"id":     "tt0097757",
			},
		},
		{
			"_id":     "977654",
			"title":   "Hijos de la guerra",
			"year":    2007,
			"rating":  nil,
			"runtime": "90 min",
			"genre": []string{
				"Documentary",
			},
			"director": "Alexandre Fuchs, Samantha Belmont, Jeremy Fourteau",
			"writer": []string{
				"Jeremy Fourteau (story)",
				"Jeff Zimbalist",
				"Michael Zimbalist",
			},
			"cast":   nil,
			"poster": "http://ia.media-imdb.com/images/M/MV5BMTIwMzUyMjcwN15BMl5BanBnXkFtZTcwNzIyMzU0MQ@@._V1_SX300.jpg",
			"imdb": map[string]interface{}{
				"rating": 8.1,
				"votes":  80,
				"id":     "tt0977654",
			},
		},
	}
)

func TestMain(m *testing.M) {
	// setup
	var err error
	server, err = NewServer("http://localhost:5984")
	if err != nil {
		os.Exit(1)
	}

	server.Delete("golang-tests")
	testsDB, err = server.Create("golang-tests")
	if err != nil {
		os.Exit(2)
	}

	server.Delete("golang-movies")
	movieDB, err = server.Create("golang-movies")
	if err != nil {
		log.Println("create error", err)
		os.Exit(3)
	}
	_, err = movieDB.Update(movies, nil)
	if err != nil {
		os.Exit(4)
	}

	server.Delete("golang-design")
	designDB, err = server.Create("golang-design")
	if err != nil {
		log.Println("create error", err)
		os.Exit(5)
	}

	// run all the tests
	code := m.Run()

	// tear down
	server.Delete("golang-tests")
	server.Delete("golang-movies")
	server.Delete("golang-design")
	os.Exit(code)
}

func TestNewServer(t *testing.T) {
	testServer, err := NewServer(DefaultBaseURL)
	if err != nil {
		t.Error(`new server error`, err)
	}
	_, err = testServer.Version()
	if err != nil {
		t.Error(`server version error`, err)
	}
}

func TestNewServerNoFullCommit(t *testing.T) {
	testServer, err := NewServerNoFullCommit(DefaultBaseURL)
	if err != nil {
		t.Error(`new server full commit error`, err)
	}
	_, err = testServer.Version()
	if err != nil {
		t.Error(`server version error`, err)
	}
}

func TestServerExists(t *testing.T) {
	testServer, err := NewServer("http://localhost:9999")
	if err != nil {
		t.Error(`new server error`, err)
	}

	_, err = testServer.Version()
	if err == nil {
		t.Error(`server version ok`)
	}
}

func TestServerConfig(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		config, err := server.Config("couchdb@localhost")
		if err != nil {
			t.Error(`server config error`, err)
		}
		if reflect.ValueOf(config).Kind() != reflect.Map {
			t.Error(`config not of type map`)
		}
	}
}

func TestServerString(t *testing.T) {
	testServer, err := NewServer(DefaultBaseURL)
	if err != nil {
		t.Error(`new server error`, err)
	}
	if testServer.String() != "Server http://localhost:5984" {
		t.Error(`server name invalid want "Server http://localhost:5984"`)
	}
}

func TestServerVars(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error(`server version error`, err)
	}
	if reflect.ValueOf(version).Kind() != reflect.String {
		t.Error(`version not of string type`)
	}

	tasks, _ := server.ActiveTasks()
	if reflect.ValueOf(tasks).Kind() != reflect.Slice {
		t.Error(`tasks not of slice type`)
	}
}

func TestServerStats(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		stats, err := server.Stats("couchdb@localhost", "")
		if err != nil {
			t.Error(`server stats error`, err)
		}
		if reflect.ValueOf(stats).Kind() != reflect.Map {
			t.Error(`stats not of map type`)
		}
		stats, err = server.Stats("couchdb@localhost", "couchdb")
		if err != nil {
			t.Error(`server stats httpd/requests error`, err)
		}
		if reflect.ValueOf(stats).Kind() != reflect.Map {
			t.Error(`httpd/requests stats not of map type`)
		}
	}
}

func TestDBs(t *testing.T) {
	aName, bName := "dba", "dbb"
	server.Create(aName)
	defer server.Delete(aName)

	server.Create(bName)
	defer server.Delete(bName)

	dbs, err := server.DBs()
	if err != nil {
		t.Error(`server DBs error`, err)
	}
	var aExist, bExist bool
	for _, v := range dbs {
		if v == aName {
			aExist = true
		} else if v == bName {
			bExist = true
		}
	}

	if !aExist {
		t.Errorf("db %s not existed in dbs", aName)
	}

	if !bExist {
		t.Errorf("db %s not existed in dbs", bName)
	}
}

func TestLen(t *testing.T) {
	aName, bName := "dba", "dbb"
	server.Create(aName)
	defer server.Delete(aName)
	server.Create(bName)
	defer server.Delete(bName)

	len, err := server.Len()
	if err != nil {
		t.Error(`server len error`, err)
	}
	if len < 2 {
		t.Error("server len should be >= 2")
	}
}

func TestGetDBMissing(t *testing.T) {
	_, err := server.Get("golang-missing")
	if err != ErrNotFound {
		t.Errorf("err = %v want ErrNotFound", err)
	}
}

func TestGetDB(t *testing.T) {
	_, err := server.Get("golang-tests")
	if err != nil {
		t.Error(`get db error`, err)
	}
}

func TestCreateDBConflict(t *testing.T) {
	conflictDBName := "golang-conflict"
	_, err := server.Create(conflictDBName)
	if err != nil {
		t.Error(`server create error`, err)
	}
	// defer s.Delete(conflictDBName)
	if !server.Contains(conflictDBName) {
		t.Error(`server not contains`, conflictDBName)
	}
	if _, err = server.Create(conflictDBName); err != ErrPreconditionFailed {
		t.Errorf("err = %v want ErrPreconditionFailed", err)
	}
	server.Delete(conflictDBName)
}

func TestCreateDB(t *testing.T) {
	_, err := server.Create("golang-create")
	if err != nil {
		t.Error(`get db failed`)
	}
	server.Delete("golang-create")
}

func TestCreateDBIllegal(t *testing.T) {
	if _, err := server.Create("_db"); err == nil {
		t.Error(`create illegal _db ok`)
	}
}

func TestDeleteDB(t *testing.T) {
	dbName := "golang-delete"
	server.Create(dbName)
	if !server.Contains(dbName) {
		t.Error(`server not contains`, dbName)
	}
	server.Delete(dbName)
	if server.Contains(dbName) {
		t.Error(`server contains`, dbName)
	}
}

func TestDeleteDBMissing(t *testing.T) {
	dbName := "golang-missing"
	err := server.Delete(dbName)
	if err != ErrNotFound {
		t.Errorf("err = %v want ErrNotFound", err)
	}
}

func TestReplicate(t *testing.T) {
	aName := "dba"
	dba, _ := server.Create(aName)
	defer server.Delete(aName)

	bName := "dbb"
	dbb, _ := server.Create(bName)
	defer server.Delete(bName)

	id, _, err := dba.Save(map[string]interface{}{"test": "a"}, nil)
	if err != nil {
		t.Error(`dba save error`, err)
	}
	result, _ := server.Replicate(aName, bName, nil)
	if v, ok := result["ok"]; !(ok && v.(bool)) {
		t.Error(`result should be ok`)
	}
	doc, err := dbb.Get(id, nil)
	if err != nil {
		t.Errorf("db %s get doc %s error %v", bName, id, err)
	}
	if v, ok := doc["test"]; ok {
		if "a" != v.(string) {
			t.Error(`doc[test] should be a, found`, v.(string))
		}
	}

	doc["test"] = "b"
	dbb.Update([]map[string]interface{}{doc}, nil)
	result, err = server.Replicate(bName, aName, nil)
	if err != nil {
		t.Error(`server replicate error`, err)
	}
	if reflect.ValueOf(result).Kind() != reflect.Map {
		t.Error(`server replicate return non-map result`)
	}

	docA, err := dba.Get(id, nil)
	if err != nil {
		t.Errorf("db %s get doc %s error %v", aName, id, err)
	}
	if v, ok := docA["test"]; ok {
		if "b" != v.(string) {
			t.Error(`docA[test] should be b, found`, v.(string))
		}
	}

	docB, err := dbb.Get(id, nil)
	if err != nil {
		t.Errorf("db %s get doc %s error %v", bName, id, err)
	}
	if v, ok := docB["test"]; ok {
		if "b" != v.(string) {
			t.Error(`docB[test] should be b, found`, v.(string))
		}
	}
}

func TestReplicateContinuous(t *testing.T) {
	aName, bName := "dba", "dbb"
	server.Create(aName)
	defer server.Delete(aName)

	server.Create(bName)
	defer server.Delete(bName)

	result, err := server.Replicate(aName, bName, map[string]interface{}{"continuous": true})
	if err != nil {
		t.Error(`server replicate error`, err)
	}

	if reflect.ValueOf(result).Kind() != reflect.Map {
		t.Error(`server replicate return non-map result`)
	}

	if v, ok := result["ok"]; !(ok && v.(bool)) {
		t.Error(`result should be ok`)
	}
}

func TestMembership(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		allNodes, clusterNodes, err := server.Membership()
		if err != nil {
			t.Error(`server membership error`, err)
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
}

func TestUUIDs(t *testing.T) {
	uuids, err := server.UUIDs(10)
	if err != nil {
		t.Error(`server uuids error`, err)
	}
	if reflect.ValueOf(uuids).Kind() != reflect.Slice {
		t.Error(`server uuids should be of type slice`)
	}
	if len(uuids) != 10 {
		t.Error(`server uuids should be of length 10, not`, len(uuids))
	}
}

func TestBasicAuth(t *testing.T) {
	testServer, _ := NewServer("http://root:password@localhost:5984/")
	_, err := testServer.Create("golang-auth")
	if err != ErrUnauthorized {
		t.Errorf("err = %v want ErrUnauthorized", err)
	}
}

func TestUserManagement(t *testing.T) {
	user := "foo"
	password := "secret"
	roles := []string{"hero"}
	server.AddUser(user, password, roles)

	token, err := server.Login(user, password)
	if err != nil {
		t.Errorf("server add user %s password %s roles %v error %s", user, password, roles, err)
	}

	if err = server.VerifyToken(token); err != nil {
		t.Error("server verify token error", err)
	}

	if err = server.Logout(token); err != nil {
		t.Error("server logout error", err)
	}

	if err = server.RemoveUser("foo"); err != nil {
		t.Error("server remove user error", err)
	}
}
