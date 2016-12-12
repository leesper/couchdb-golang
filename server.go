package couchdb

import(
  "encoding/json"
  "net/url"
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
  _, _, jsonData := s.resource.GetJSON("", nil, url.Values{})
  var version string
  _ = json.Unmarshal(*jsonData["version"], &version)
  return version
}
