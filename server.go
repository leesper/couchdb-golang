package couchdb

import(
  "encoding/json"
  // "log"
  "net/url"
  "strconv"
)

// Server represents a CouchDB server instance.
type Server struct {
  resource *Resource
}

// NewServer creates an object on behalf of CouchDB instance in address urlStr.
func NewServer(urlStr string) *Server {
  res, err := NewResource(urlStr, nil)
  if err != nil {
    return nil
  }

  return &Server{
    resource: res,
  }
}

// Version returns the version info about CouchDB instance.
func (s *Server)Version() string {
  var jsonMap map[string]interface{}

  _, _, jsonData := s.resource.GetJSON("", nil, url.Values{})
  if jsonData == nil {
    return ""
  }
  _ = json.Unmarshal(*jsonData, &jsonMap)

  return jsonMap["version"].(string)
}

// ActiveTasks lists of running tasks.
func (s *Server)ActiveTasks() []interface{} {
  var jsonArr []interface{}

  _, _, jsonData := s.resource.GetJSON("_active_tasks", nil, url.Values{})
  if jsonData == nil {
    return nil
  }
  _ = json.Unmarshal(*jsonData, &jsonArr)

  return jsonArr
}

// DBs returns a list of all the databases in the CouchDB instance.
func (s *Server)DBs() []string {
  var dbs []string

  _, _, jsonData := s.resource.GetJSON("_all_dbs", nil, url.Values{})
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

  _, _, jsonData := s.resource.GetJSON("_membership", nil, url.Values{})
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

  _, _, jsonData := s.resource.PostJSON("_replicate", nil, body, url.Values{})
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

// Create creates a database with the given name. Return false if failed.
func (s *Server)Create(name string) (*Database, bool) {
  status, _, _ := s.resource.PutJSON(name, nil, nil, url.Values{})
  if status != Created {
    return nil, false
  }

  db := s.GetDatabase(name)
  if db == nil {
    return nil, false
  }

  return db, true
}

// Delete deletes a database with the given name. Return false if failed.
func (s *Server)Delete(db string) bool {
  status, _, _ := s.resource.DeleteJSON(db, nil, url.Values{})

  if status == OK {
    return true
  }
  return false
}

// GetDatabase gets a database instance with the given name. Return nil if failed.
func (s *Server)GetDatabase(name string) *Database {
  dbResStr := s.newResource(name)
  if len(dbResStr) == 0 {
    return nil
  }
  db := NewDatabase(dbResStr)
  if db == nil {
    return nil
  }
  status, _, _ := db.resource.Head("", nil, url.Values{})
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
// RemoveUser removes regular user in authentication database.
