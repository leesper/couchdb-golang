package couchdb

import (
  // "log"
  "reflect"
  "strings"
  "testing"
)

var s *Server

func init() {
  s = NewServer("http://localhost:5984")
}

func TestVersion(t *testing.T) {
  version := s.Version()
  if reflect.ValueOf(version).Kind() != reflect.String {
    t.Error(`version should be string`)
  }
  if !strings.HasPrefix(version, "2") {
    t.Error(`version should be above 2`)
  }
}
