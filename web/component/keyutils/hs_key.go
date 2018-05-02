package keyutils

import (
	//"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var HSKey []byte

func GetHS256Key() []byte {
	return HSKey
}

func init() {
	var err error
	baseDir, _ := os.Getwd()
	log.Println("baseDir", baseDir)

	destdir := baseDir
	// if !strings.HasSuffix(baseDir, "uniapi") {
	// 	destdir = strings.Split(baseDir, "uniapi")[0]
	// 	destdir = destdir + "uniapi"
	// }

	if strings.Contains(baseDir, "uniapi/src/web") {
		destdir = strings.Split(baseDir, "uniapi")[0]
		destdir = destdir + "uniapi/src/web/config/key/hs.key"
	} else {
		if strings.HasSuffix(baseDir, "bin") {
			destdir = destdir + "/config/key/hs.key"
		} else {
			destdir = destdir + "/bin/config/key/hs.key"
		}
	}

	log.Println("destdir", destdir)

	HSKey, err = ioutil.ReadFile(destdir)
	if err != nil {
		panic(err)
	}
}
