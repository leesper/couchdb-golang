package couchdb

import(
  "encoding/json"
  "net/url"
  // "log"
)

type Server struct {
  resource *Resource
}

// NewServer creates an object on behalf of CouchDB instance in address urlStr.
func NewServer(urlStr string) *Server {
  res, _ := NewResource(urlStr, nil)
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
