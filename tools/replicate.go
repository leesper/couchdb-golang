package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/leesper/couchdb-golang"
)

type tuple struct {
	left, right string
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func findPath(s string) (string, string) {
	if s == "." {
		return couchdb.DefaultBaseURL, ""
	}

	if !strings.HasPrefix(s, "http") {
		return couchdb.DefaultBaseURL, s
	}

	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}

	var base string
	if u.User != nil {
		base = fmt.Sprintf("%s://%s@%s", u.Scheme, u.User.String(), u.Host)
	} else {
		base = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	}
	res, err := couchdb.NewResource(base, nil)
	if err != nil {
		log.Fatal(err)
	}
	_, data, err := res.GetJSON("", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	info := map[string]interface{}{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		log.Fatal(err)
	}

	if _, ok := info["couchdb"]; !ok {
		log.Fatal(fmt.Sprintf("%s doe not appear to be a CouchDB", s))
	}

	return base, strings.TrimLeft(u.Path, "/")
}

func main() {
	isContinous := flag.Bool("continous", false, "trigger continous replication in couchdb")
	isCompact := flag.Bool("compact", false, "compact target database after replication")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] <source> <target>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(os.Args) != 3 {
		log.Fatal("need source and target arguments")
	}

	src, tgt := os.Args[1], os.Args[2]

	sbase, spath := findPath(src)
	source, err := couchdb.NewServer(sbase)
	if err != nil {
		log.Fatal(err)
	}

	tbase, tpath := findPath(tgt)
	target, err := couchdb.NewServer(tbase)
	if err != nil {
		log.Fatal(err)
	}

	// check dtabase name specs
	if strings.Contains(tpath, "*") {
		log.Fatal("invalid target path: must be single db or empty")
	}

	all, err := source.DBs()
	if err != nil {
		log.Fatal(err)
	}

	if spath == "" {
		log.Fatal("source database must be specified")
	}

	sources := []string{}
	for _, db := range all {
		if !strings.HasPrefix(db, "_") { // skip reserved names
			var ok bool
			ok, err = path.Match(spath, db)
			if err != nil {
				log.Fatal(err)
			}
			if ok {
				sources = append(sources, db)
			}
		}
	}

	if len(sources) == 0 {
		log.Fatalf("no source databases match glob %s", spath)
	}

	databases := []tuple{}
	if len(sources) > 1 && tpath != "" {
		log.Fatal("target path must be empty with multiple sources")
	} else if len(sources) == 1 {
		databases = append(databases, tuple{sources[0], tpath})
	} else {
		for _, source := range sources {
			databases = append(databases, tuple{source, source})
		}
	}

	// actual replication
	for _, t := range databases {
		sdb, tdb := t.left, t.right
		start := time.Now()
		fmt.Printf("%s -> %s\n", sdb, tdb)

		if !target.Contains(tdb) {
			_, err = target.Create(tdb)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		sdb = fmt.Sprintf("%s/%s", sbase, sdb)
		var opts map[string]interface{}
		if *isContinous {
			opts = map[string]interface{}{"continuous": true}
		}
		_, err = target.Replicate(sdb, tdb, opts)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%.1fs\n", time.Now().Sub(start).Seconds())
	}

	if *isCompact {
		for _, t := range databases {
			tdb := t.right
			fmt.Println("compact", tdb)
			targetDB, _ := target.Get(tdb)
			targetDB.Compact()
		}
	}
}
