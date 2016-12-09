// CouchDB resource
//
// This is the low-level wrapper functions of HTTP methods
//
// Used for communicating with CouchDB Server
package couchdb

import (
  "bytes"
  "encoding/json"
  "io"
  "io/ioutil"
  "net/http"
  "net/url"
  "strings"
)

const (
  UnknownError = -1
  UrlParseError = -2
  JSONMarshalError = -3
  JSONUnmarshalError = -4
  OK = 200
  Created = 201
  Accepted = 202
  NotModified = 304
  BadRequest = 400
  Unauthorized = 401
  Forbidden = 403
  NotFound = 404
  ResouceNotAllowed = 405
  NotAcceptable = 406
  Conflict = 409
  PreconditionFailed = 412
  BadContentType = 415
  RequestedRangeNotSatisfiable = 416
  ExpectationFailed = 417
  InternalServerError = 500
)

// Resource handles all requests to CouchDB
type Resource struct {
  header *http.Header
  base *url.URL
}

func NewResource(urlStr string, header *http.Header) (*Resource, error) {
  u, err := url.Parse(urlStr)
  if err != nil {
    return nil, err
  }
  var h *http.Header
  if header == nil {
    h = new(http.Header)
  } else {
    h = header
  }

  return &Resource{
    header: h,
    base: u,
  }, nil
}

// Head is a wrapper around http.Head
func (r *Resource)Head(path string, header *http.Header, params url.Values) (int, http.Header, []byte) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }
  return request(http.MethodHead, u, header, nil, params)
}

// Get is a wrapper around http.Get
func (r *Resource)Get(path string, header *http.Header, params url.Values) (int, http.Header, []byte) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }
  return request(http.MethodGet, u, header, nil, params)
}

// Post is a wrapper around http.Post
func (r *Resource)Post(path string, header *http.Header, body []byte, params url.Values) (int, http.Header, []byte) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }
  return request(http.MethodPost, u, header, bytes.NewReader(body), params)
}

// Delete is a wrapper around http.Delete
func (r *Resource)Delete(path string, header *http.Header, params url.Values) (int, http.Header, []byte) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }
  return request(http.MethodDelete, u, header, nil, params)
}

// Put is a wrapper around http.Put
func (r *Resource)Put(path string, header *http.Header, body []byte, params url.Values) (int, http.Header, []byte) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }
  return request(http.MethodPut, u, header, bytes.NewReader(body), params)
}

// GetJSON issues a GET to the specified URL, with data returned as json
func (r *Resource)GetJSON(path string, header *http.Header, params url.Values) (int, http.Header, map[string]*json.RawMessage) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }
  return requestJSON(http.MethodGet, u, header, nil, params)
}

// PostJSON issues a POST to the specified URL, with data returned as json
func (r *Resource)PostJSON(path string, header *http.Header, body map[string]interface{}, params url.Values) (int, http.Header, map[string]*json.RawMessage) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }

  jsonBody, err := json.Marshal(body)
  if err != nil {
    return JSONMarshalError, nil, nil
  }

  return requestJSON(http.MethodPost, u, header, bytes.NewReader(jsonBody), params)
}

// DeleteJSON issues a DELETE to the specified URL, with data returned as json
func (r *Resource)DeleteJSON(path string, header *http.Header, params url.Values) (int, http.Header, map[string]*json.RawMessage) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }

  return requestJSON(http.MethodDelete, u, header, nil, params)
}

// PutJSON issues a PUT to the specified URL, with data returned as json
func (r *Resource)PutJSON(path string, header *http.Header, body map[string]interface{}, params url.Values) (int, http.Header, map[string]*json.RawMessage) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }

  jsonBody, err := json.Marshal(body)
  if err != nil {
    return JSONMarshalError, nil, nil
  }

  return requestJSON(http.MethodPut, u, header, bytes.NewReader(jsonBody), params)
}

// helper function to make real request
func requestJSON(method string, u *url.URL, header *http.Header, body io.Reader, params url.Values) (int, http.Header, map[string]*json.RawMessage) {
  s, h, d := request(method, u, header, body, params)
  if d != nil && h.Get("Content-type") == "application/json" {
    var jsonData map[string]*json.RawMessage
    err := json.Unmarshal(d, &jsonData)
    if err != nil {
      return JSONUnmarshalError, h, nil
    }
    return s, h, jsonData
  }
  return s, h, nil
}

// helper function to make real request
func request(method string, u *url.URL, header *http.Header, body io.Reader, params url.Values) (int, http.Header, []byte) {
  method = strings.ToUpper(method)

  u.RawQuery = params.Encode()
  var username, password string
  if u.User != nil {
    username = u.User.Username()
    password, _ = u.User.Password()
  }

  req, err := http.NewRequest(method, u.String(), body)
  if err != nil {
    return UnknownError, req.Header, nil
  }

  if len(username) > 0 && len(password) > 0 {
    req.SetBasicAuth(username, password)
  }

  // Accept and Content-type are highly recommended for CouchDB
  setDefault(&req.Header, "Accept", "application/json")
  setDefault(&req.Header, "Content-type", "application/json")
  updateHeader(&req.Header, header)

  rsp, err := http.DefaultClient.Do(req)
  if err != nil {
    return UnknownError, rsp.Header, nil
  }
  defer rsp.Body.Close()
  data, err := ioutil.ReadAll(rsp.Body)

  return rsp.StatusCode, rsp.Header, data
}

// setDefault sets the default value if key not existe in header
func setDefault(header *http.Header, key, value string) {
  if header.Get(key) == "" {
    header.Set(key, value)
  }
}

// updateHeader updates existing header with new values
func updateHeader(header *http.Header, extra *http.Header) {
  for k, _ := range *extra {
    header.Set(k, extra.Get(k))
  }
}
