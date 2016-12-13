package couchdb

import (
  // "log"
  "reflect"
  "strings"
  "testing"
)

var s *Server

func init() {
  s = NewServer("http://root:likejun@localhost:5984")
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
  kind := reflect.ValueOf(dbs).Kind()
  elemKind := reflect.TypeOf(dbs).Elem().Kind()
  if kind != reflect.Slice || elemKind != reflect.String {
    t.Error(`dbs shold be string slice`)
  }
}

func TestMembership(t *testing.T) {
  allNodes, clusterNodes := s.Membership()
  if allNodes == nil || clusterNodes == nil {
    t.Error(`unauthorized`)
  }

  kind := reflect.ValueOf(allNodes).Kind()
  elemKind := reflect.TypeOf(allNodes).Elem().Kind()

  if kind != reflect.Slice || elemKind != reflect.String {
    t.Error(`clusterNodes should be`)
  }

  kind = reflect.ValueOf(clusterNodes).Kind()
  elemKind = reflect.TypeOf(clusterNodes).Elem().Kind()

  if kind != reflect.Slice || elemKind != reflect.String {
    t.Error(`allNodes should be`)
  }
}

func TestReplicate(t *testing.T) {
  rsp := s.Replicate("db_a", "db_b", nil)
  if reflect.ValueOf(rsp).Kind() != reflect.Map {
    t.Error(`should return a map`)
  }
}
