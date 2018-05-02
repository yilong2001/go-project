package keyutils

import (
	//"io/ioutil"
	"log"
	"os"
	"strings"
)

var alipayprivatekeypath string
var alipaypublickeypath string
var zhifubaopublickeypath string

func GetAlipayPrivateKeyPath() string {
	return alipayprivatekeypath
}

func GetAlipayPublicKeyPath() string {
	return alipaypublickeypath
}

func GetZhifubaoPublicKeyPath() string {
	return zhifubaopublickeypath
}

func GetAlipayPrivateKeyStr() []byte {
	key, err := getTokens(GetAlipayPrivateKeyPath())
	if err != nil {
		panic(err)
	}
	return (key)
}

func GetAlipayPublicKeyStr() []byte {
	key, err := getTokens(GetAlipayPublicKeyPath())
	if err != nil {
		panic(err)
	}
	return (key)
}

func GetZhifubaoPublicKeyStr() []byte {
	key, err := getTokens(GetZhifubaoPublicKeyPath())
	if err != nil {
		panic(err)
	}
	return (key)
}

func init() {
	baseDir, _ := os.Getwd()
	log.Println("rsa baseDir", baseDir)

	var destdir string
	destdir = baseDir
	// if !strings.HasSuffix(baseDir, "uniapi") {
	//  destDir = strings.Split(baseDir, "uniapi")[0]
	//  destDir = destDir + "uniapi"
	// }

	if strings.Contains(baseDir, "uniapi/src/web") {
		destdir = strings.Split(baseDir, "uniapi")[0]
		destdir = destdir + "uniapi/src/web/config/key/ali/"
	} else {
		if strings.HasSuffix(baseDir, "bin") {
			destdir = destdir + "/config/key/ali/"
		} else {
			destdir = destdir + "/bin/config/key/ali/"
		}
	}

	alipayprivatekeypath = destdir + "private_key.pem"

	alipaypublickeypath = destdir + "public_key.pem"

	zhifubaopublickeypath = destdir + "ali_public_key.pem"

	log.Println(GetAlipayPrivateKeyStr())

	log.Println(GetAlipayPublicKeyStr())
}
