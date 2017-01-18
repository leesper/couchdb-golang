package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewDefaultDB(t *testing.T) {
	dbDefault, err := NewDatabase("golang-default")
	if err != nil {
		t.Errorf("new default database error %v", err)
	}
	if err = dbDefault.Available(); err == nil {
		t.Error(`db available`)
	}
}

func TestNewDB(t *testing.T) {
	newDB := "golang-newdb"
	server.Create(newDB)
	defer server.Delete(newDB)
	dbNew, err := NewDatabase(fmt.Sprintf("%s/%s", DefaultBaseURL, newDB))
	if err != nil {
		t.Error(`new database error`, err)
	}
	if err = dbNew.Available(); err != nil {
		t.Error(`db not available, error`, err)
	}
}

func TestSaveNew(t *testing.T) {
	doc := map[string]interface{}{"doc": "bar"}
	id, rev, err := testsDB.Save(doc, nil)
	if err != nil {
		t.Error(`db save error`, err)
	}
	if id != doc["_id"].(string) {
		t.Errorf("invalid id: %q", id)
	}
	if rev != doc["_rev"].(string) {
		t.Errorf("invalid rev: %q", rev)
	}
}

func TestSaveNewWithID(t *testing.T) {
	doc := map[string]interface{}{"_id": "foo"}
	id, rev, err := testsDB.Save(doc, nil)
	if err != nil {
		t.Error(`db save error`, err)
	}
	if doc["_id"].(string) != "foo" {
		t.Errorf("doc[_id] = %s, not foo", doc["_id"])
	}
	if id != "foo" {
		t.Errorf("id = %s, not foo", id)
	}
	if rev != doc["_rev"].(string) {
		t.Errorf("invalid rev: %q", rev)
	}
}

func TestSaveExisting(t *testing.T) {
	doc := map[string]interface{}{}
	idOld, revOld, err := testsDB.Save(doc, nil)
	if err != nil {
		t.Error(`db save error`, err)
	}
	doc["foo"] = true
	idNew, revNew, err := testsDB.Save(doc, nil)
	if err != nil {
		t.Error(`db save foo error`, err)
	}
	if idOld != idNew {
		t.Errorf("ids are not equal old %s new %s", idOld, idNew)
	}
	if doc["_rev"].(string) != revNew {
		t.Errorf("invalid rev %s want %s", doc["_rev"].(string), revNew)
	}
	if revOld == revNew {
		t.Errorf("new rev is equal to old rev %s", revOld)
	}
}

func TestSaveNewBatch(t *testing.T) {
	doc := map[string]interface{}{"_id": "foo"}
	_, rev, err := testsDB.Save(doc, url.Values{"batch": []string{"ok"}})
	if err != nil {
		t.Error(`db save batch error`, err)
	}
	if len(rev) > 0 {
		t.Error(`rev not empty`, rev)
	}
	if r, ok := doc["_rev"]; ok {
		t.Error(`doc has _rev field`, r.(string))
	}
}

func TestSaveExistingBatch(t *testing.T) {
	doc := map[string]interface{}{"_id": "bar"}
	idOld, revOld, err := testsDB.Save(doc, nil)
	if err != nil {
		t.Error(`db save error`, err)
	}

	idNew, revNew, err := testsDB.Save(doc, url.Values{"batch": []string{"ok"}})
	if err != nil {
		t.Error(`db save batch error`, err)
	}

	if idOld != idNew {
		t.Errorf("old id %s not equal to new id %s", idOld, idNew)
	}

	if len(revNew) > 0 {
		t.Error(`rev not empty`, revNew)
	}

	if doc["_rev"].(string) != revOld {
		t.Errorf("doc[_rev] %s not equal to old rev %s", doc["_rev"].(string), revOld)
	}
}

func TestDatabaseExists(t *testing.T) {
	if err := testsDB.Available(); err != nil {
		t.Error(`golang-tests not available, error`, err)
	}
	dbMissing, _ := NewDatabase("golang-missing")
	if err := dbMissing.Available(); err == nil {
		t.Error(`golang-missing available`)
	}
}

func TestDatabaseName(t *testing.T) {
	name, err := testsDB.Name()
	if err != nil {
		t.Error(`db name error`, err)
	}
	if name != "golang-tests" {
		t.Errorf("db name %s, want golang-tests", name)
	}
}

func TestDatabaseString(t *testing.T) {
	if testsDB.String() != "Database http://localhost:5984/golang-tests" {
		t.Error(`db string invalid`, testsDB)
	}
}

func TestCommit(t *testing.T) {
	if err := testsDB.Commit(); err != nil {
		t.Error(`db commit error`, err)
	}
}

func TestCreateLargeDoc(t *testing.T) {
	var buf bytes.Buffer
	// 10MB
	for i := 0; i < 110*1024; i++ {
		buf.WriteString("0123456789")
	}
	doc := map[string]interface{}{"data": buf.String()}
	if err := testsDB.Set("large", doc); err != nil {
		t.Error(`db set error`, err)
	}
	doc, err := testsDB.Get("large", nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if doc["_id"].(string) != "large" {
		t.Errorf("doc[_id] = %s, want large", doc["_id"].(string))
	}
	err = testsDB.DeleteDoc(doc)
	if err != nil {
		t.Error(`db delete doc error`, err)
	}
}

func TestDocIDQuoting(t *testing.T) {
	doc := map[string]interface{}{"foo": "bar"}
	err := testsDB.Set("foo/bar", doc)
	if err != nil {
		t.Error(`db set error`, err)
	}
	doc, err = testsDB.Get("foo/bar", nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if doc["foo"].(string) != "bar" {
		t.Errorf("doc[foo] = %s want bar", doc["foo"].(string))
	}
	err = testsDB.Delete("foo/bar")
	if err != nil {
		t.Error(`db delete error`, err)
	}
	_, err = testsDB.Get("foo/bar", nil)
	if err == nil {
		t.Error(`db get foo/bar ok`)
	}
}

func TestDisallowNaN(t *testing.T) {
	doc := map[string]interface{}{"number": math.NaN()}
	err := testsDB.Set("foo", doc)
	if err == nil {
		t.Error(`db set NaN ok`)
	}
}

func TestDisallowNilID(t *testing.T) {
	err := testsDB.DeleteDoc(map[string]interface{}{"_id": nil, "_rev": nil})
	if err == nil {
		t.Error(`db delete doc with id nil ok`)
	}
	err = testsDB.DeleteDoc(map[string]interface{}{"_id": "foo", "_rev": nil})
	if err == nil {
		t.Error(`db delete doc with rev nil ok`)
	}
}

func TestDocRevs(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{"bar": 42}
	err := testsDB.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	oldRev := doc["_rev"].(string)
	doc["bar"] = 43
	err = testsDB.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	newRev := doc["_rev"].(string)

	newDoc, err := testsDB.Get(uuid, nil)
	if err != nil {
		t.Error("db get error", err)
	}
	if newRev != newDoc["_rev"].(string) {
		t.Errorf("new doc rev %s want %s", newDoc["_rev"].(string), newRev)
	}
	newDoc, err = testsDB.Get(uuid, url.Values{"rev": []string{newRev}})
	if err != nil {
		t.Error("db get error", err)
	}
	if newRev != newDoc["_rev"].(string) {
		t.Errorf("new doc rev %s want %s", newDoc["_rev"].(string), newRev)
	}
	oldDoc, err := testsDB.Get(uuid, url.Values{"rev": []string{oldRev}})
	if err != nil {
		t.Error("db get error", err)
	}
	if oldRev != oldDoc["_rev"].(string) {
		t.Errorf("old doc rev %s want %s", oldDoc["_rev"].(string), oldRev)
	}

	revs, err := testsDB.Revisions(uuid, nil)
	if err != nil {
		t.Error(`db revisions error`, err)
	}
	if revs[0]["_rev"].(string) != newRev {
		t.Errorf("revs first %s want %s", revs[0]["_rev"].(string), newRev)
	}
	if revs[1]["_rev"].(string) != oldRev {
		t.Errorf("revs second %s not equal to %s", revs[1]["_rev"].(string), oldRev)
	}
	_, err = testsDB.Revisions("crap", nil)
	if err == nil {
		t.Error(`db revisions crap ok`)
	}

	err = testsDB.Compact()
	if err != nil {
		t.Error("db compact error", err)
	}

	info, err := testsDB.Info("")
	if err != nil {
		t.Error(`db info error`, err)
	}
	for info["compact_running"].(bool) {
		info, err = testsDB.Info("")
		if err != nil {
			t.Error(`db info error`, err)
		}
	}

	_, err = testsDB.Get(uuid, url.Values{"rev": []string{oldRev}})
	if err == nil {
		t.Errorf("db get compacted doc ok, rev = %s", oldRev)
	}
}

func TestAttachmentCRUD(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{"bar": 42}
	testsDB.Set(uuid, doc)
	oldRev := doc["_rev"].(string)

	testsDB.PutAttachment(doc, []byte("Foo bar"), "foo.txt", "text/plain")
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err := testsDB.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	attachments := reflect.ValueOf(doc["_attachments"])
	foo := reflect.ValueOf(attachments.MapIndex(reflect.ValueOf("foo.txt")).Interface())
	length := int(foo.MapIndex(reflect.ValueOf("length")).Interface().(float64))
	if length != len("Foo bar") {
		t.Errorf("length %d want %d", length, len("Foo bar"))
	}
	contentType := foo.MapIndex(reflect.ValueOf("content_type")).Interface().(string)
	if contentType != "text/plain" {
		t.Errorf("content type %s want text/plain", contentType)
	}

	data, err := testsDB.GetAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`get attachment error`, err)
	}
	if string(data) != "Foo bar" {
		t.Errorf("db get attachment %s want Foo bar", string(data))
	}

	data, err = testsDB.GetAttachmentID(uuid, "foo.txt")
	if err != nil {
		t.Error(`get attachment id error`, err)
	}
	if string(data) != "Foo bar" {
		t.Errorf("db get attachment id %s want Foo bar", string(data))
	}

	oldRev = doc["_rev"].(string)
	err = testsDB.DeleteAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`db delete attachment error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = testsDB.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if _, ok := doc["_attachments"]; ok {
		t.Error(`doc attachments still existed`)
	}
}

func TestAttachmentWithFiles(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{"bar": 42}
	err := testsDB.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	oldRev := doc["_rev"].(string)
	fileObj := []byte("Foo bar baz")

	err = testsDB.PutAttachment(doc, fileObj, "foo.txt", mime.TypeByExtension(".txt"))
	if err != nil {
		t.Error("db put attachment error", err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = testsDB.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	attachments := reflect.ValueOf(doc["_attachments"])
	foo := reflect.ValueOf(attachments.MapIndex(reflect.ValueOf("foo.txt")).Interface())
	length := int(foo.MapIndex(reflect.ValueOf("length")).Interface().(float64))
	if length != len("Foo bar baz") {
		t.Errorf("length %d want %d", length, len("Foo bar"))
	}
	contentType := foo.MapIndex(reflect.ValueOf("content_type")).Interface().(string)
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("content type %s want text/plain; charset=utf-8", contentType)
	}

	data, err := testsDB.GetAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`get attachment error`, err)
	}
	if string(data) != "Foo bar baz" {
		t.Errorf("db get attachment %s want Foo bar", string(data))
	}

	data, err = testsDB.GetAttachmentID(uuid, "foo.txt")
	if err != nil {
		t.Error(`get attachment id error`, err)
	}
	if string(data) != "Foo bar baz" {
		t.Errorf("db get attachment id %s want Foo bar", string(data))
	}

	oldRev = doc["_rev"].(string)
	err = testsDB.DeleteAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`db delete attachment error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = testsDB.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if _, ok := doc["_attachments"]; ok {
		t.Error(`doc attachments still existed`)
	}
}

func TestAttachmentCRUDFromFS(t *testing.T) {
	uuid := GenerateUUID()
	content := "Foo bar baz"
	tmpFileName := filepath.Join(os.TempDir(), "foo.txt")
	tmpFile, err := os.Create(tmpFileName)
	if err != nil {
		t.Error(`create file error`, err)
	}
	_, err = tmpFile.Write([]byte(content))
	if err != nil {
		t.Error(`write file error`, err)
	}
	tmpFile.Close()

	tmpFile, err = os.Open(tmpFileName)
	if err != nil {
		t.Error(`open file error`, err)
	}
	defer tmpFile.Close()

	data, err := ioutil.ReadAll(tmpFile)
	if err != nil {
		t.Error(`read tmp file error`, err)
	}

	doc := map[string]interface{}{"bar": 42}
	err = testsDB.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	oldRev := doc["_rev"].(string)
	err = testsDB.PutAttachment(doc, data, "foo.txt", mime.TypeByExtension(filepath.Ext(tmpFileName)))
	if err != nil {
		t.Error(`put attachment error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = testsDB.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	attachment := reflect.ValueOf(doc["_attachments"])
	foo := reflect.ValueOf(attachment.MapIndex(reflect.ValueOf("foo.txt")).Interface())
	length := int(foo.MapIndex(reflect.ValueOf("length")).Interface().(float64))
	if len(content) != length {
		t.Errorf("length %d want %d", length, len(content))
	}
	contentType := foo.MapIndex(reflect.ValueOf("content_type")).Interface().(string)
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("content type %s want text/plain; charset=utf-8", contentType)
	}

	data, err = testsDB.GetAttachment(doc, "foo.txt")
	if err != nil {
		t.Error(`get attachment error`, err)
	}
	if string(data) != content {
		t.Error(`get attachment should be `, content)
	}

	data, err = testsDB.GetAttachmentID(uuid, "foo.txt")
	if err != nil {
		t.Error(`get attachment id error`, err)
	}
	if string(data) != content {
		t.Error(`get attachment id should be `, content)
	}

	if err = testsDB.DeleteAttachment(doc, "foo.txt"); err != nil {
		t.Error(`delete attachment file error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = testsDB.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	if _, ok := doc["_attachments"]; ok {
		t.Error(`doc attachments still existed`)
	}
}

func TestEmptyAttachment(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{}
	err := testsDB.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	oldRev := doc["_rev"].(string)

	err = testsDB.PutAttachment(doc, []byte(""), "empty.txt", mime.TypeByExtension(".txt"))
	if err != nil {
		t.Error(`put attachment error`, err)
	}
	if oldRev == doc["_rev"].(string) {
		t.Error(`doc[_rev] == oldRev`)
	}

	doc, err = testsDB.Get(uuid, nil)
	if err != nil {
		t.Error(`db get error`, err)
	}

	attachment := reflect.ValueOf(doc["_attachments"])
	empty := reflect.ValueOf(attachment.MapIndex(reflect.ValueOf("empty.txt")).Interface())
	length := int(empty.MapIndex(reflect.ValueOf("length")).Interface().(float64))
	if length != 0 {
		t.Errorf("length %d want %d", length, 0)
	}
}

func TestDefaultAttachment(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{}
	err := testsDB.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	_, err = testsDB.GetAttachment(doc, "missing.txt")
	if err == nil {
		t.Error(`db get attachment ok`)
	}
}

func TestAttachmentNoFilename(t *testing.T) {
	uuid := GenerateUUID()
	doc := map[string]interface{}{}
	err := testsDB.Set(uuid, doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	err = testsDB.PutAttachment(doc, []byte(""), "", "")
	if err == nil {
		t.Error(`db put attachment with no file name ok`)
	}
}

func TestJSONAttachment(t *testing.T) {
	doc := map[string]interface{}{}
	err := testsDB.Set(GenerateUUID(), doc)
	if err != nil {
		t.Error(`db set doc error`, err)
	}
	err = testsDB.PutAttachment(doc, []byte("{}"), "test.json", "application/json")
	if err != nil {
		t.Error(`db put attachment json error`, err)
	}
	data, err := testsDB.GetAttachment(doc, "test.json")
	if err != nil {
		t.Error(`db get attachment json error`, err)
	}
	if string(data) != "{}" {
		t.Errorf("data = %s want {}", string(data))
	}
}

func TestBulkUpdateConflict(t *testing.T) {
	docs := []map[string]interface{}{
		{
			"type": "Person",
			"name": "John Doe",
		},
		{
			"type": "Person",
			"name": "Mary Jane",
		},
		{
			"type": "Person",
			"name": "Gotham City",
		},
	}

	testsDB.Update(docs, nil)

	// update the first doc to provoke a conflict in the next bulk update
	doc := map[string]interface{}{}
	for k, v := range docs[0] {
		doc[k] = v
	}
	testsDB.Set(doc["_id"].(string), doc)

	results, err := testsDB.Update(docs, nil)
	if err != nil {
		t.Error(`db update error`, err)
	}
	if results[0].Err != ErrConflict {
		t.Errorf("db update conflict err %v want ErrConflict", results[0].Err)
	}
}

func TestCopyDocConflict(t *testing.T) {
	testsDB.Set("foo1", map[string]interface{}{"status": "idle"})
	testsDB.Set("bar1", map[string]interface{}{"status": "testing"})
	_, err := testsDB.Copy("foo1", "bar1", "")
	if err != ErrConflict {
		t.Errorf(`db copy returns %v, want ErrConflict`, err)
	}
}

func TestCopyDocOverwrite(t *testing.T) {
	foo2 := map[string]interface{}{"status": "testing"}
	bar2 := map[string]interface{}{"status": "idle"}
	testsDB.Set("foo2", foo2)
	testsDB.Set("bar2", bar2)
	result, err := testsDB.Copy("foo2", "bar2", bar2["_rev"].(string))
	if err != nil {
		t.Error(`db copy error`, err)
	}
	doc, _ := testsDB.Get("bar2", nil)
	if result != doc["_rev"].(string) {
		t.Errorf("db copy returns %s want %s", result, doc["_rev"].(string))
	}
	if doc["status"].(string) != "testing" {
		t.Errorf("db copy status = %s, want testing", doc["status"].(string))
	}
}

func TestChanges(t *testing.T) {
	options := url.Values{
		"style": []string{"all_docs"},
	}
	_, err := testsDB.Changes(options)
	if err != nil {
		t.Error(`db change error`, err)
	}
}

func TestPurge(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// TODO: purge not implemented in CouchDB 2.0.0
	if !strings.HasPrefix(version, "2") {
		doc := map[string]interface{}{"a": "b"}
		err := testsDB.Set("purge", doc)
		if err != nil {
			t.Error(`db set error`, err)
		}
		result, err := testsDB.Purge([]map[string]interface{}{doc})
		if err != nil {
			t.Error(`db purge error`, err)
		}

		purgeSeq := int(result["purge_seq"].(float64))
		if purgeSeq != 1 {
			t.Errorf("db purge seq=%d want 1", purgeSeq)
		}
	}
}

func TestSecurity(t *testing.T) {
	secDoc, err := testsDB.GetSecurity()
	if err != nil {
		t.Error(`get security should return true`)
	}
	if len(secDoc) > 0 {
		t.Error(`secDoc should be empty`)
	}
	if testsDB.SetSecurity(map[string]interface{}{
		"names": []string{"test"},
		"roles": []string{},
	}) != nil {
		t.Error(`set security should return true`)
	}
}

func TestDBContains(t *testing.T) {
	doc := map[string]interface{}{
		"type": "Person",
		"name": "Jason Statham",
	}
	id, _, err := testsDB.Save(doc, nil)
	if err != nil {
		t.Error(`db save error`, err)
	}
	if err = testsDB.Contains(id); err != nil {
		t.Error(`db contains error`, err)
	}
}

func TestDBSetGetDelete(t *testing.T) {
	doc := map[string]interface{}{
		"type": "Person",
		"name": "Jason Statham",
	}
	err := testsDB.Set("Mechanic", doc)
	if err != nil {
		t.Error(`db set error`, err)
	}
	_, err = testsDB.Get("Mechanic", nil)
	if err != nil {
		t.Error(`db get error`, err)
	}
	err = testsDB.Delete("Mechanic")
	if err != nil {
		t.Error(`db delete error`, err)
	}
}

func TestDBDocIDsAndLen(t *testing.T) {
	doc := map[string]interface{}{
		"type": "Person",
		"name": "Jason Statham",
	}

	err := testsDB.Set("Mechanic", doc)
	if err != nil {
		t.Error(`db set error`, err)
	}

	ids, err := testsDB.DocIDs()
	if err != nil {
		t.Error(`db doc ids error`, err)
	}

	length, err := testsDB.Len()
	if err != nil {
		t.Error(`db len error`, err)
	}
	if length != len(ids) {
		t.Errorf("Len() returns %d want %d", length, len(ids))
	}
}

func TestGetSetRevsLimit(t *testing.T) {
	err := testsDB.SetRevsLimit(10)
	if err != nil {
		t.Error(`db set revs limit error`, err)
	}
	limit, err := testsDB.GetRevsLimit()
	if err != nil {
		t.Error(`db get revs limit error`, err)
	}
	if limit != 10 {
		t.Error(`limit should be 10`)
	}
}

func TestCleanup(t *testing.T) {
	err := testsDB.Cleanup()
	if err != nil {
		t.Error(`db clean up error`, err)
	}
}

func TestParseSelectorSyntax(t *testing.T) {
	_, err := parseSelectorSyntax(`title == "Spacecataz" && year == 2004 && director == "Dave Willis"`)
	if err != nil {
		t.Error("parse selector syntax error", err)
	}

	_, err = parseSelectorSyntax(`year >= 1990 && (director == "George Lucas" || director == "Steven Spielberg")`)
	if err != nil {
		t.Error("parse selector syntax error", err)
	}

	_, err = parseSelectorSyntax(`year >= 1900 && year <= 2000 && nor(year == 1990, year == 1989, year == 1997)`)
	if err != nil {
		t.Error("parse selector syntax error", err)
	}

	_, err = parseSelectorSyntax(`_id > nil && all(genre, []string{"Comedy", "Short"})`)
	if err != nil {
		t.Error("parse selector syntax error", err)
	}

	_, err = parseSelectorSyntax(`_id > nil && any(genre, genre == "Short" || genre == "Horror" || score >= 8)`)
	if err != nil {
		t.Error("parse selector syntax error", err)
	}

	_, err = parseSelectorSyntax(`exists(director, true, "wrongParam")`)
	if err == nil {
		t.Error("parse exists function ok, should be 2 parameters")
	}

	_, err = parseSelectorSyntax(`typeof(genre, "array", "wrongParam")`)
	if err == nil {
		t.Error("parse typeof function ok, should be 2 parameters")
	}

	_, err = parseSelectorSyntax(`in(director, []string{"Mike Portnoy", "Vitali Kanevsky"}, "wrongParam")`)
	if err == nil {
		t.Error("parse in function ok, should be 2 parameters")
	}

	_, err = parseSelectorSyntax(`nin(year, []int{1990, 1992, 1998}, "wrongParam")`)
	if err == nil {
		t.Error("parse nin function ok, should be 2 parameters")
	}

	_, err = parseSelectorSyntax(`size(genre, 2, "wrongParam")`)
	if err == nil {
		t.Error("parse size function ok, should be 2 parameters")
	}

	_, err = parseSelectorSyntax(`mod(year, 2, 1, "wrongParam")`)
	if err == nil {
		t.Error("parse mod function ok, should be 3 parameters")
	}

	_, err = parseSelectorSyntax(`regex(title, "^A", "wrongParam")`)
	if err == nil {
		t.Error("parse regex function ok, should be 2 parameters")
	}
}

func TestParseSortSyntax(t *testing.T) {
	_, err := parseSortSyntax([]string{"fieldNameA", "fieldNameB"})
	if err != nil {
		t.Error("parse sort syntax error", err)
	}

	_, err = parseSortSyntax([]string{"fieldNameA.subFieldA", "fieldNameB.subFieldB"})
	if err != nil {
		t.Error("parse sort syntax error", err)
	}

	_, err = parseSortSyntax([]string{"desc(fieldName1)", "asc(fieldName2)"})
	if err != nil {
		t.Error("parse sort syntax error", err)
	}

	_, err = parseSortSyntax([]string{"desc(fieldName1.subField1)", "asc(fieldName2.subField2)"})
	if err != nil {
		t.Error("parse sort syntax error", err)
	}
}

func TestQueryYearAndID(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `_id > nil && in(year, []int{2007, 2004})`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
		    "selector": {
					"$and": [
		        	{
		            "_id": { "$gt": null }
		        	},
		        	{
		            "year": {
		              "$in": [2007, 2004]
		          	}
		        	}
		    	]
				}
		}`

		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query year and id not equal")
		}
	}
}

func TestQueryYearOrDirector(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `year == 1989 && (director == "Ademir Kenovic" || director == "Dezs Garas")`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
			"selector": {
				"year": 1989,
				"$or": [
					{ "director": "Ademir Kenovic" },
					{ "director": "Dezs Garas" }
				]
			}
		}`

		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query year or director not equal")
		}
	}
}

func TestQueryYearGteLteNot(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `year >= 1989 && year <= 2006 && year != 2004`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
			"selector": {
				"year": {
		      "$gte": 1989
		    },
		    "year": {
		      "$lte": 2006
		    },
		    "$not": {
		      "year": 2004
		    }
			}
		}`

		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query year gte lte not not equal")
		}
	}
}

func TestQueryIMDBRatingNor(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `imdb.rating >= 6 && imdb.rating <= 9 && nor(imdb.rating == 8.1, imdb.rating == 8.2)`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
	  	"selector": {
	    	"imdb.rating": {
	        "$gte": 6
	    	},
	    	"imdb.rating": {
	        "$lte": 9
	    	},
	    	"$nor": [
	        { "imdb.rating": 8.1 },
	        { "imdb.rating": 8.2 }
	    	]
			}
		}`

		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query imdb rating nor not equal")
		}
	}
}

func TestQueryGenreAll(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `_id > nil && all(genre, []string{"Comedy", "Short"})`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
	  	"selector": {
	    	"_id": {
	        "$gt": null
	    	},
	    	"genre": {
	        "$all": ["Comedy","Short"]
	    	}
			}
		}`

		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query genre all not equal")
		}
	}
}

func TestQueryGenreElemMatch(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `_id > nil && any(genre, genre == "Horror" || genre == "Comedy")`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
	  	"selector": {
				"$and": [
					{
						"_id": {
						"$gt": null
						}
					},
					{
						"genre": {
							"$elemMatch": {
								"$or": [
									{
										"$eq": "Horror"
									},
									{
										"$eq": "Comedy"
									}
								]
							}
						}
					}
				]
			}
		}`

		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query genre elem match not equal")
		}
	}
}

func TestQueryRegex(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `regex(director, "^D")`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
	  	"selector": {
	    	"director": {"$regex": "^D"}
			}
		}`
		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query regex not equal")
		}
	}
}

func TestQueryYearIDNin(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `_id > nil && nin(year, []int{1989, 1990})`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
	  	"selector": {
	    	"$and": [
	        {
	          "_id": { "$gt": null }
	        },
	        {
	          "year": {
	          	"$nin": [1989, 1990]
	          }
	        }
	    	]
			}
		}`
		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query year id nin not equal")
		}
	}
}

func TestQueryTypeAndExists(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `typeof(poster, "string") && exists(runtime, true)`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
	  	"selector": {
	    	"$and": [
	        {
		    		"poster": {
							"$type": "string"
		    		}
					},
					{
		    		"runtime": {
							"$exists": true
		    		}
					}
	    	]
			}
		}`
		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query type and exists not equal")
		}
	}
}

func TestQuerySizeAndMod(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, `size(writer, 2) && mod(year, 2, 0)`, nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
	  	"selector": {
	    	"$and": [
	        {
		    		"writer": {
							"$size": 2
		    		}
					},
					{
		    		"year": {
						"$mod": [2, 0]
		    		}
					}
	    	]
			}
		}`
		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query size and mod not equal")
		}
	}
}

func TestRatingOrYear(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		docsQuery, err := movieDB.Query(nil, "rating != nil || year < 2000", nil, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
	  	"selector":{
	    	"$or": [
	        {
		    		"rating": {
							"$ne": null
		    		}
					},
					{
		    		"year": {
	  					"$lt": 2000
		    		}
					}
	    	]
			}
		}`
		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsRaw, docsQuery) {
			t.Error("db query rating or year not equal")
		}
	}
}

func TestQuerySortLimitSkip(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		fields := []string{"_id", "_rev", "year", "director"}
		selector := `year > 1989`
		sorts := []string{"desc(_id)"}
		docsQuery, err := movieDB.Query(fields, selector, sorts, 5, 2, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
	    "selector": {
	      "year": {"$gt": 1989}
	    },
	    "fields": ["_id", "_rev", "year", "director"],
	    "limit": 5,
	    "skip": 2,
	    "sort": [{"_id": "desc"}]
		}`
		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query sort limit skip not equal")
		}
	}
}

func TestIndexCRUD(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		designName, indexName, err := movieDB.PutIndex([]string{"asc(year)"}, "", "year-index")
		if err != nil {
			t.Error("db put index error", err)
		}

		indexResult, err := movieDB.GetIndex()
		if err != nil {
			t.Error("db get index error", err)
		}

		var totalRows float64
		err = json.Unmarshal(*indexResult["total_rows"], &totalRows)
		if err != nil {
			t.Error("json unmarshal total rows error", err)
		}
		if int(totalRows) < 2 {
			t.Error("index total rows should be >= 2")
		}

		var idxes = []*json.RawMessage{}
		err = json.Unmarshal(*indexResult["indexes"], &idxes)
		if err != nil {
			t.Error("json unmarshal indexes error", err)
		}
		idxMap := map[string]*json.RawMessage{}
		found := false
		idxName := ""
		idxType := ""
		idxDef := map[string]*json.RawMessage{}
		defFields := []*json.RawMessage{}
		fieldMap := map[string]string{}
		for _, idx := range idxes {
			json.Unmarshal(*idx, &idxMap)
			json.Unmarshal(*idxMap["name"], &idxName)
			if idxName != "year-index" {
				continue
			}
			found = true
			json.Unmarshal(*idxMap["type"], &idxType)
			if idxType != "json" {
				t.Error("index type not json")
			}
			json.Unmarshal(*idxMap["def"], &idxDef)
			json.Unmarshal(*idxDef["fields"], &defFields)
			if len(defFields) != 1 {
				t.Error("index def fields != 1")
			}
			json.Unmarshal(*defFields[0], &fieldMap)
			if fieldMap["year"] != "asc" {
				t.Errorf("index year order %s want asc", fieldMap["year"])
			}
		}
		if !found {
			t.Error("index year not found")
		}

		err = movieDB.DeleteIndex(designName, indexName)
		if err != nil {
			t.Error("db delete index error", err)
		}
	}
}

func TestQueryDoubleSort(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		designName, indexName, err := movieDB.PutIndex([]string{"imdb.rating", "imdb.votes"}, "", "imdb-index")
		if err != nil {
			t.Error("db put index error", err)
		}

		fields := []string{"_id", "_rev", "year", "title"}
		selector := `year > 1989 && imdb.rating > 6 && imdb.votes > 100`
		sorts := []string{"imdb.rating", "imdb.votes"}
		docsQuery, err := movieDB.Query(fields, selector, sorts, nil, nil, nil)
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
		  "selector": {
		  	"year": {"$gt": 1989},
				"imdb.rating": {"$gt": 6},
				"imdb.votes": {"$gt": 100}
		  },
		  "fields": ["_id", "_rev", "year", "title"],
		  "sort": [{"imdb.rating": "asc"}, {"imdb.votes": "asc"}]
		}`
		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}

		if !reflect.DeepEqual(docsQuery, docsRaw) {
			t.Error("db query double sort not equal")
		}

		movieDB.DeleteIndex(designName, indexName)
	}
}

func TestQueryUseIndex(t *testing.T) {
	version, err := server.Version()
	if err != nil {
		t.Error("server version error", err)
	}
	// CouchDB 2.0 feature
	if strings.HasPrefix(version, "2") {
		designName, indexName, err := movieDB.PutIndex([]string{"year"}, "", "year-index")
		if err != nil {
			t.Error("db put index error", err)
		}
		if indexName != "year-index" {
			t.Errorf("db put index return %s want year-index", indexName)
		}

		fields := []string{"_id", "_rev", "year", "title"}
		selector := `year > 1989`
		sorts := []string{"asc(year)"}
		docsQuery, err := movieDB.Query(fields, selector, sorts, 5, 0, []string{designName, indexName})
		if err != nil {
			t.Error("db query error", err)
		}

		var rawJSON = `
		{
		  "selector": {
		    "year": {"$gt": 1989}
		  },
		  "fields": ["_id", "_rev", "year", "title"],
		  "sort": [{"year": "asc"}],
		  "limit": 5,
		  "skip": 0,
			"use_index": [%q, %q]
		}`
		rawJSON = fmt.Sprintf(rawJSON, designName, indexName)
		docsRaw, err := movieDB.QueryJSON(rawJSON)
		if err != nil {
			t.Error("db query json error", err)
		}
		if !reflect.DeepEqual(docsRaw, docsQuery) {
			t.Error("db query usd index not equal")
		}
	}
}
