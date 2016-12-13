package couchdb

import(
  "encoding/json"
  "net/url"
  // "log"
)

type Server struct {
  resource *Resource
}

func NewServer(urlStr string) *Server {
  res, _ := NewResource(urlStr, nil)
  return &Server{
    resource: res,
  }
}

func (s *Server)Version() string {
  var version string
  var jsonMap map[string]*json.RawMessage

  _, _, jsonData := s.resource.GetJSON("", nil, url.Values{})
  if jsonData == nil {
    return version
  }
  _ = json.Unmarshal(*jsonData, &jsonMap)
  _ = json.Unmarshal(*jsonMap["version"], &version)

  return version
}

func (s *Server)ActiveTasks() []*json.RawMessage {
  var jsonArr []*json.RawMessage

  _, _, jsonData := s.resource.GetJSON("_active_tasks", nil, url.Values{})
  if jsonData == nil {
    return nil
  }
  _ = json.Unmarshal(*jsonData, &jsonArr)

  return jsonArr
}

func (s *Server)DBs() []string {
  var dbs []string

  _, _, jsonData := s.resource.GetJSON("_all_dbs", nil, url.Values{})
  if jsonData == nil {
    return nil
  }
  _ = json.Unmarshal(*jsonData, &dbs)

  return dbs
}

func (s *Server)Membership() ([]string, []string) {
  var allNodes []string
  var clusterNodes []string
  var jsonMap map[string]*json.RawMessage

  _, _, jsonData := s.resource.GetJSON("_membership", nil, url.Values{})
  if jsonData == nil {
    return nil, nil
  }

  _ = json.Unmarshal(*jsonData, &jsonMap)
  if _, ok := jsonMap["error"]; ok {
    return nil, nil
  }

  _ = json.Unmarshal(*jsonMap["all_nodes"], &allNodes)
  _ = json.Unmarshal(*jsonMap["cluster_nodes"], &clusterNodes)

  return allNodes, clusterNodes
}