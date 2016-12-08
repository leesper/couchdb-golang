// CouchDB resource
//
// This is the low-level wrapper functions of HTTP methods
//
// Used for communicating with CouchDB Server
package couchdb

import (
  "encoding/json"
  "http"
  "io"
  "io/ioutil"
  "net/url"
  "strings"
)

type StatusCode int

const (
  UnknownError StatusCode = -1
  UrlParseError StatusCode = -2
  OK StatusCode = 200
  Created StatusCode = 201
  Accepted StatusCode = 202
  NotModified StatusCode = 304
  BadRequest StatusCode = 400
  Unauthorized StatusCode = 401
  Forbidden StatusCode = 403
  NotFound StatusCode = 404
  ResouceNotAllowed StatusCode = 405
  NotAcceptable StatusCode = 406
  Conflict StatusCode = 409
  PreconditionFailed StatusCode = 412
  BadContentType StatusCode = 415
  RequestedRangeNotSatisfiable StatusCode = 416
  ExpectationFailed StatusCode = 417
  InternalServerError StatusCode = 500
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
    h = make(http.Header)
  } else {
    h = header
  }

  return &Resource{
    header: h,
    base: u,
  }, nil
}

// Head is a wrapper around http.Head
func (r *Resource)Head(path string, header *http.Header, params url.Values) (StatusCode, *http.Header, []byte) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }
  return request(http.MethodHead, u, header, nil, params)
}

// Get is a wrapper around http.Get
func Get(path string, header *http.Header, params url.Values) (StatusCode, *http.Header, []byte) {
  u, err := r.base.Parse(path)
  if err != nil {
    return UrlParseError, nil, nil
  }
  return request(http.MethodGet, u, header, nil, params)
}

// Post is a wrapper around http.Post
func Post(urlStr string, header http.Header, body []byte, params url.Values) (StatusCode, *http.Header, []byte) {}

// Delete is a wrapper around http.Delete
func Delete(urlStr string, header http.Header, params url.Values) (StatusCode, *http.Header, []byte) {}

// Put is a wrapper around http.Put
func Put(urlStr string, header http.Header, body []byte, params url.Values) (StatusCode, *http.Header, []byte) {}

// GetJSON issues a GET to the specified URL, with data returned as json
func GetJSON(urlStr string, header http.Header, params url.Values) (StatusCode, *http.Header, map[string]interface{}) {}

// PostJSON issues a POST to the specified URL, with data returned as json
func PostJSON(urlStr string, header http.Header, body map[string]interface{}, params url.Values) (StatusCode, *http.Header, map[string]interface{}) {}

// DeleteJSON issues a DELETE to the specified URL, with data returned as json
func DeleteJSON(urlStr string, header http.Header, params url.Values) (StatusCode, *http.Header, map[string]interface{}) {}

// PutJSON issues a PUT to the specified URL, with data returned as json
func PutJSON(urlStr string, header http.Header, body map[string]interface{}, params url.Values) (StatusCode, *http.Header, map[string]*json.RawMessage) {}

// helper function to make real request
func requestJSON(method string, urlStr string, header *http.Header, body io.Reader, params url.Values) (StatusCode, *http.Header, map[string]interface{}) {
  statusCode, header, data := request(method, url, header, body, params)
  if header != nil && data != nil && header.Get("Content-type") == "application/json" {
    var jsonData map[string]interface{}
    err := json.Unmarshal(data, &jsonData)
    if err != nil {
      return UnknownError, nil, nil
    }
    return statusCode, header, jsonData
  }
  return statusCode, header, nil
}

// helper function to make real request
func request(method string, u *url.URL, header *http.Header, body io.Reader, params url.Values) (StatusCode, *http.Header, []byte) {
  method = strings.ToUpper(method)

  u.RawQuery = params.Encode()
  var username, password string
  if u.User != nil {
    username = u.User.Username()
    password, _ := u.User.Password()
  }

  req, err := http.NewRequest(method, u.String(), body)
  if err != nil {
    return UnknownError, nil, nil
  }

  if len(username) > 0 && len(password) > 0 {
    req.SetBasicAuth(username, password)
  }

  // Accept and Content-type are highly recommended for CouchDB
  setDefault(req.Header, "Accept", "application/json")
  setDefault(req.Header, "Content-type", "application/json")
  updateHeader(&req.Header, header)

  rsp, err := http.DefaultClient.Do(req)
  if err != nil {
    return UnknownError, nil, nil
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
  for k, _ := range extra {
    header.Set(extra.Get(k))
  }
}
