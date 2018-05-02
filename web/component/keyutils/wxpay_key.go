package keyutils

import (
	//"io/ioutil"
	"log"
	"os"
	"strings"
)

var wxpayClientCertPath string
var wxpayClientKeyPath string

func GetWxPayClientCertPath() string {
	return wxpayClientCertPath
}

func GetWxPayClientKeyPath() string {
	return wxpayClientKeyPath
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
		destdir = destdir + "uniapi/src/web/config/key/wx/"
	} else {
		if strings.HasSuffix(baseDir, "bin") {
			destdir = destdir + "/config/key/wx/"
		} else {
			destdir = destdir + "/bin/config/key/wx/"
		}
	}

	wxpayClientCertPath = destdir + "apiclient_cert.pem"

	wxpayClientKeyPath = destdir + "apiclient_key.pem"

	log.Println(wxpayClientCertPath, wxpayClientKeyPath)
}
