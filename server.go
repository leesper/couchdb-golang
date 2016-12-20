package couchdb

import(
  "encoding/json"
  // "fmt"
  // "log"
  "net/http"
  "net/url"
  "strconv"
  "strings"
)

// Server represents a CouchDB server instance.
type Server struct {
  resource *Resource
}

// NewServer creates an object on behalf of CouchDB instance in address urlStr.
func NewServer(urlStr string) (*Server, error) {
  res, err := NewResource(urlStr, nil)
  if err != nil {
    return nil, err
  }

  return &Server{
    resource: res,
  }, nil
}

// NewServerFullCommit creates a CouchDB instance in address urlStr.
// Disable X-Couch-Full-Commit by setting fullCommit to false.
func NewServerFullCommit(urlStr string, fullCommit bool) *Server {
  s := NewServer(urlStr)
  if s == nil {
    return nil
  }

  if !fullCommit {
    s.resource.header.Set("X-Couch-Full-Commit", "false")
  }
  return s
}

// Version returns the version info about CouchDB instance.
func (s *Server)Version() string {
  var jsonMap map[string]interface{}

  _, _, jsonData := s.resource.GetJSON("", nil, nil)
  if jsonData == nil {
    return ""
  }
  _ = json.Unmarshal(*jsonData, &jsonMap)

  return jsonMap["version"].(string)
}

// ActiveTasks lists of running tasks.
func (s *Server)ActiveTasks() []interface{} {
  var jsonArr []interface{}

  _, _, jsonData := s.resource.GetJSON("_active_tasks", nil, nil)
  if jsonData == nil {
    return nil
  }
  _ = json.Unmarshal(*jsonData, &jsonArr)

  return jsonArr
}

// DBs returns a list of all the databases in the CouchDB instance.
func (s *Server)DBs() []string {
  var dbs []string

  _, _, jsonData := s.resource.GetJSON("_all_dbs", nil, nil)
  if jsonData == nil {
    return nil
  }
  _ = json.Unmarshal(*jsonData, &dbs)
  return dbs
}


// Membership displays the nodes that are part of the cluster as clusterNodes.
// The field allNodes displays all nodes this node knows about, including the
// ones that are part of cluster.
func (s *Server)Membership() ([]string, []string) {
  var jsonMap map[string]*json.RawMessage

  _, _, jsonData := s.resource.GetJSON("_membership", nil, nil)
  if jsonData == nil {
    return nil, nil
  }

  _ = json.Unmarshal(*jsonData, &jsonMap)
  if _, ok := jsonMap["error"]; ok {
    return nil, nil
  }

  var allNodes []string
  var clusterNodes []string

  _ = json.Unmarshal(*jsonMap["all_nodes"], &allNodes)
  _ = json.Unmarshal(*jsonMap["cluster_nodes"], &clusterNodes)

  return allNodes, clusterNodes
}

// Replicate requests, configure or stop a replication operation.
func (s *Server)Replicate(source, target string, options map[string]interface{}) map[string]interface{} {
  var jsonMap map[string]interface{}

  body := map[string]interface{} {
    "source": source,
    "target": target,
  }

  if options != nil {
    for k, v := range options {
      body[k] = v
    }
  }

  _, _, jsonData := s.resource.PostJSON("_replicate", nil, body, nil)
  if jsonData == nil {
    return nil
  }
  _ = json.Unmarshal(*jsonData, &jsonMap)

  return jsonMap
}

// Stats returns a JSON object containing the statistics for the running server.
// func (s *Server)Stats(entry string) map[string]interface{} {
//   var jsonMap map[string]interface{}
//   _, _, jsonData := s.resource.GetJSON("_stats", nil, url.Values{})
//   if jsonData != nil {
//     return nil
//   }
//
//   _ = json.Unmarshal(*jsonData, &jsonMap)
//   log.Println(jsonMap, len(jsonMap))
//   return jsonMap
// }

// UUIDs requests one or more Universally Unique Identifiers from the CouchDB instance.
// The response is a JSON object providing a list of UUIDs.
// count - Number of UUIDs to return. Default is 1.
func (s *Server)UUIDs(count int) []string {
  if count <= 0 {
    count = 1
  }

  values := url.Values{}
  values.Set("count", strconv.Itoa(count))

  _, _, jsonData := s.resource.GetJSON("_uuids", nil, values)
  if jsonData == nil {
    return nil
  }

  var jsonMap map[string]*json.RawMessage
  _ = json.Unmarshal(*jsonData, &jsonMap)
  if _, ok := jsonMap["uuids"]; !ok {
    return nil
  }

  var uuids []string
  _ = json.Unmarshal(*jsonMap["uuids"], &uuids)

  return uuids
}

// Create returns a database instance with the given name, returns true if created,
// if database already existed, returns false, *Database will be nil if failed.
func (s *Server)Create(name string) (*Database, bool) {
  status, _, _ := s.resource.PutJSON(name, nil, nil, nil)

  // PreconditionFailed means database with the given name already existed
  if status != Created && status != PreconditionFailed {
    return nil, false
  }

  db := s.GetDatabase(name)

  return db, db != nil && status == Created
}

// Delete deletes a database with the given name. Return false if failed.
func (s *Server)Delete(db string) bool {
  status, _, _ := s.resource.DeleteJSON(db, nil, nil)

  if status == OK {
    return true
  }
  return false
}

// GetDatabase gets a database instance with the given name. Return nil if failed.
func (s *Server)GetDatabase(name string) *Database {
  res := s.resource.NewResourceWithURL(name)
  if res == nil {
    return nil
  }

  db := NewDatabaseWithResource(res)
  if db == nil {
    return nil
  }

  status, _, _ := db.resource.Head("", nil, nil)
  if status != OK {
    return nil
  }
  return db
}

// newResource returns an url string representing a resource under server.
func (s *Server)newResource(resource string) string {
  resourceUrl, err := s.resource.base.Parse(resource)
  if err != nil {
    return ""
  }
  return resourceUrl.String()
}

// AddUser adds regular user in authentication database.
// Returns id and rev of the registered user.
func (s *Server)AddUser(name, password string, roles []string) (string, string) {
  db := s.GetDatabase("_users")
  if db == nil {
    return "", ""
  }

  if roles == nil {
    roles = []string{}
  }

  userDoc := map[string]interface{}{
    "_id": "org.couchdb.user:" + name,
    "name": name,
    "password": password,
    "roles": roles,
    "type": "user",
  }

  return db.Save(userDoc)
}

// Login regular user in CouchDB, returns authentication token.
func (s *Server)Login(name, password string) (string, bool) {
  body := map[string]interface{}{
    "name": name,
    "password": password,
  }
  status, header, _ := s.resource.PostJSON("_session", nil, body, nil)
  if status != OK {
    return "", false
  }

  tokenPart := strings.Split(header.Get("Set-Cookie"), ";")[0]
  token := strings.Split(tokenPart, "=")[1]
  return token, status == OK
}

// Verify regular user token
func (s *Server)VerifyToken(token string) bool {
  header := http.Header{}
  header.Set("Cookie", strings.Join([]string{"AuthSession", token}, "="))
  status, _, _ := s.resource.GetJSON("_session", &header, nil)
  return status == OK
}

// Logout regular user in CouchDB
func (s *Server)Logout(token string) bool {
  header := http.Header{}
  header.Set("Cookie", strings.Join([]string{"AuthSession", token}, "="))
  status, _, _ := s.resource.DeleteJSON("_session", &header, nil)
  return status == OK
}

// RemoveUser removes regular user in authentication database.
func (s *Server)RemoveUser(name string) bool {
  db := s.GetDatabase("_users")
  if db == nil {
    return false
  }
  docId := "org.couchdb.user:" + name
  return db.Delete(docId)
}
