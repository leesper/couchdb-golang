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

func TestActiveTasks(t *testing.T) {
  jsonArr := s.ActiveTasks()
  if reflect.ValueOf(jsonArr).Kind() != reflect.Slice {
    t.Error(`jsonArr should be slice`)
  }
}

func TestDBs(t *testing.T) {
  dbs := s.DBs()
  if reflect.ValueOf(dbs).Kind() != reflect.Slice {
    t.Error(`dbs should be slice`)
  }

  if reflect.TypeOf(dbs).Elem().Kind() != reflect.String {
    t.Error(`dbs shold be string slice`)
  }
}
