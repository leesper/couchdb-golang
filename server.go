package couchdb

type Server struct {
  resource *Resource
}

func NewServer(urlStr string) *Server {
  res, _ := NewResource(urlStr, nil)
  return &Server{
    resource: res,
  }
}
