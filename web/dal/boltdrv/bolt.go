package boltdrv

import (
	"log"
	"net/http"

	"web/component/cfgutils"

	"github.com/boltdb/bolt"
	"github.com/gorilla/sessions"
	"github.com/yosssi/boltstore/reaper"
	"github.com/yosssi/boltstore/store"
)

var db *bolt.DB = nil

func GetDb(webapicfg *cfgutils.WebApiConfig) *bolt.DB {
	if db != nil {
		return db
	}

	var err error
	// Open a Bolt database.
	db, err = bolt.Open(webapicfg.DbFile, 0666, nil)
	if err != nil {
		log.Panic(err)
	}

	return db
}

func CloseDb() {
	db.Close()
	db = nil
}

func BoltDbSession