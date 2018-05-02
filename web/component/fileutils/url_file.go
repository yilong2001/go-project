package fileutils

import (
	"io/ioutil"
	"log"
	"net/http"
)

func GetUrlFile(url string) ([]byte, error) {
	oid, err := http.Get(url)
	if err != nil {
		log.Println("get file : ", err, url)
		return nil, err
	}
	defer oid.Body.Close()
	oids, err2 := ioutil.ReadAll(oid.Body)
	if err2 != nil {
		log.Println("get file : ", err2, url)
		return nil, err2
	}
	return oids, nil
}
