// Copyright 2018 Mathieu Lonjaret

package main

import (
	"flag"
	"log"
	"os"
	"runtime/trace"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func main() {
	f, err := os.Create("/home/mpl/trace.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := trace.Start(f); err != nil {
		log.Fatal(err)
	}
	defer trace.Stop()
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		log.Fatal("need ldb dir as argument")
	}
	ldbCompaction(args[0])
}

func ldbCompaction(dbdir string) {
	db, err := leveldb.OpenFile(dbdir, &opt.Options{
		DisableCompactionBackoff: true,
	})
	if err != nil {
		log.Fatalf("Could not open ldb dir: %v", err)
	}

	if err := db.CompactRange(util.Range{nil, nil}); err != nil {
		log.Fatalf("compact range error: %v", err)
	}

	for _, v := range []string{"leveldb.stats", "leveldb.sstables", "leveldb.openedtables"} {
		if val, err := db.GetProperty(v); err != nil {
			log.Fatal(err)
		} else {
			println(val)
		}
	}

	if err := db.Close(); err != nil {
		log.Fatalf("close DB err: %v", err)
	}

	return
}
