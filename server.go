package couchdb

import(
  "encoding/json"
  "fmt"
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

// NewServer creates a CouchDB server instance in address urlStr.
func NewServer(urlStr string) (*Server, error) {
 return newServer(urlStr, true)
}

// NewServerNoFullCommit creates a CouchDB server instance in address urlStr
// with X-Couch-Full-Commit disabled.
func NewServerNoFullCommit(urlStr string) (*Server, error) {
  return newServer(urlStr, false)
}

func newServer(urlStr string, fullCommit bool) (*Server, error) {
  res, err := NewResource(urlStr, nil)
  if err != nil {
    return nil, err
  }

  s := &Server{
    resource: res,
  }

  if !fullCommit {
    s.resource.header.Set("X-Couch-Full-Commit", "false")
  }
  return s, nil
}

// Config returns the entire CouchDB server configuration as JSON structure.
func (s *Server)Config(node string) (map[string]map[string]string, error) {
  _, data, err := s.resource.GetJSON(fmt.Sprintf("_node/%s/_config", node), nil, nil)
  if err != nil {
    return nil, err
  }
  var config map[string]map[string]string
  err = json.Unmarshal(*data, &config)
  if err != nil {
    return nil, err
  }
  return config, nil
}

// Version returns the version info about CouchDB instance.
func (s *Server)Version() (string, error) {
  var jsonMap map[string]interface{}

  _, data, err := s.resource.GetJSON("", nil, nil)
  if err != nil {
    return "", err
  }
  err = json.Unmarshal(*data, &jsonMap)
  if err != nil {
    return "", err
  }

  return jsonMap["version"].(string), nil
}

func (s *Server)String() string {
  return fmt.Sprintf("Server %s", s.resource.base)
}

// ActiveTasks lists of running tasks.
func (s *Server)ActiveTasks() ([]interface{}, error) {
  var tasks []interface{}
  _, data, err := s.resource.GetJSON("_active_tasks", nil, nil)
  if err != nil {
    return nil, err
  }
  err = json.Unmarshal(*data, &tasks)
  if err != nil {
    return nil, err
  }
  return tasks, nil
}

// DBs returns a list of all the databases in the CouchDB server instance.
func (s *Server)DBs() ([]string, error) {
  var dbs []string
  _, data, err := s.resource.GetJSON("_all_dbs", nil, nil)
  if err != nil {
    return nil, err
  }
  err = json.Unmarshal(*data, &dbs)
  if err != nil {
    return nil, err
  }
  return dbs, nil
}

// Stats returns a JSON object containing the statistics for the running server.
func (s *Server)Stats(node, entry string) (map[string]interface{}, error) {
  var stats map[string]interface{}
  _, data, err := s.resource.GetJSON(fmt.Sprintf("_node/%s/_stats/%s", node, entry), nil, url.Values{})
  if err != nil {
    return nil, err
  }
  err = json.Unmarshal(*data, &stats)
  if err != nil {
    return nil, err
  }
  return stats, nil
}

// Len returns the number of dbs in CouchDB server instance.
func (s *Server)Len() (int, error) {
  dbs, err := s.DBs()
  if err != nil {
    return -1, err
  }
  return len(dbs), nil
}

// Create returns a database instance with the given name, returns true if created,
// if database already existed, returns false, *Database will be nil if failed.
func (s *Server)Create(name string) (*Database, error) {
  _, _, err := s.resource.PutJSON(name, nil, nil, nil)

  // ErrPreconditionFailed means database with the given name already existed
  if err != nil && err != ErrPreconditionFailed {
    return nil, err
  }

  db, getErr := s.Get(name)
  if getErr != nil {
    return nil, getErr
  }
  return db, err
}

// Delete deletes a database with the given name. Return false if failed.
func (s *Server)Delete(name string) error {
  _, _, err := s.resource.DeleteJSON(name, nil, nil)
  return err
}

// Get gets a database instance with the given name. Return nil if failed.
func (s *Server)Get(name string) (*Database, error) {
  res, err := s.resource.NewResourceWithURL(name)
  if err != nil {
    return nil, err
  }

  db, err := NewDatabaseWithResource(res)
  if err != nil {
    return nil, err
  }

  _, _, err = db.resource.Head("", nil, nil)
  if err != nil {
    return nil, err
  }
  return db, nil
}

// Contains returns true if a db with given name exsited.
func (s *Server)Contains(name string) bool {
  _, _, err := s.resource.Head(name, nil, nil)
  return err == nil
}

//////////////////////////////////////////////////////////////////////////////


// Membership displays the nodes that are part of the cluster as clusterNodes.
// The field allNodes displays all nodes this node knows about, including the
// ones that are part of cluster.
func (s *Server)Membership() ([]string, []string) {
  var jsonMap map[string]*json.RawMessage

  _, jsonData, err := s.resource.GetJSON("_membership", nil, nil)
  if err != nil {
    return nil, nil
  }

  json.Unmarshal(*jsonData, &jsonMap)
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

  _, jsonData, err := s.resource.PostJSON("_replicate", nil, body, nil)
  if err != nil {
    return nil
  }
  json.Unmarshal(*jsonData, &jsonMap)

  return jsonMap
}


// UUIDs requests one or more Universally Unique Identifiers from the CouchDB instance.
// The response is a JSON object providing a list of UUIDs.
// count - Number of UUIDs to return. Default is 1.
func (s *Server)UUIDs(count int) []string {
  if count <= 0 {
    count = 1
  }

  values := url.Values{}
  values.Set("count", strconv.Itoa(count))

  _, jsonData, err := s.resource.GetJSON("_uuids", nil, values)
  if err != nil {
    return nil
  }

  var jsonMap map[string]*json.RawMessage
  json.Unmarshal(*jsonData, &jsonMap)
  if _, ok := jsonMap["uuids"]; !ok {
    return nil
  }

  var uuids []string
  _ = json.Unmarshal(*jsonMap["uuids"], &uuids)

  return uuids
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
func (s *Server)AddUser(name, password string, roles []string) (string, string, error) {
  db, err := s.Get("_users")
  if err != nil {
    return "", "", err
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

  id, rev := db.Save(userDoc)
  return id, rev, nil
}

// Login regular user in CouchDB, returns authentication token.
func (s *Server)Login(name, password string) (string, bool) {
  body := map[string]interface{}{
    "name": name,
    "password": password,
  }
  header, _, err := s.resource.PostJSON("_session", nil, body, nil)
  if err != nil {
    return "", false
  }

  tokenPart := strings.Split(header.Get("Set-Cookie"), ";")[0]
  token := strings.Split(tokenPart, "=")[1]
  return token, err == nil
}

// Verify regular user token
func (s *Server)VerifyToken(token string) bool {
  header := http.Header{}
  header.Set("Cookie", strings.Join([]string{"AuthSession", token}, "="))
  _, _, err := s.resource.GetJSON("_session", header, nil)
  return err == nil
}

// Logout regular user in CouchDB
func (s *Server)Logout(token string) bool {
  header := http.Header{}
  header.Set("Cookie", strings.Join([]string{"AuthSession", token}, "="))
  _, _, err := s.resource.DeleteJSON("_session", header, nil)
  return err == nil
}

// RemoveUser removes regular user in authentication database.
func (s *Server)RemoveUser(name string) bool {
  db, err := s.Get("_users")
  if err != nil {
    return false
  }
  docId := "org.couchdb.user:" + name
  return db.Delete(docId)
}
