package keyutils

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type XTokenKey struct {
	Ver        string
	PrivateKey []byte
	PublicKey  []byte
}

func compRSAPrivateKeyPath(dir, ver string) string {
	return dir + ver + "/private_key.pem"
}

func compRSAPublicKeyPath(dir, ver string) string {
	return dir + ver + "/public_key.pem"
}

func getTokens(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func initRSAKeyByVer(dir, ver string) (*XTokenKey, error) {
	prikey, err := getTokens(compRSAPrivateKeyPath(dir, ver))
	if err != nil {
		return nil, err
	}

	pubkey, err := getTokens(compRSAPublicKeyPath(dir, ver))
	if err != nil {
		return nil, err
	}

	return &XTokenKey{
		Ver:        ver,
		PrivateKey: prikey,
		PublicKey:  pubkey,
	}, nil
}

var allRsaKyes = []*XTokenKey{}

func GetAllRSAKeys() []*XTokenKey {
	return allRsaKyes
}

func GetRSAKeyByVer(ver string) *XTokenKey {
	for _, key := range allRsaKyes {
		if key.Ver == ver {
			return key
		}
	}

	return nil
}

func init() {
	baseDir, _ := os.Getwd()
	log.Println("rsa baseDir", baseDir)

	var destdir string
	destdir = baseDir
	// if !strings.HasSuffix(baseDir, "uniapi") {
	// 	destDir = strings.Split(baseDir, "uniapi")[0]
	// 	destDir = destDir + "uniapi"
	// }

	if strings.Contains(baseDir, "uniapi/src/web") {
		destdir = strings.Split(baseDir, "uniapi")[0]
		destdir = destdir + "uniapi/src/web/config/key/rsa/"
	} else {
		if strings.HasSuffix(baseDir, "bin") {
			destdir = destdir + "/config/key/rsa/"
		} else {
			destdir = destdir + "/bin/config/key/rsa/"
		}
	}

	log.Println("rsa destDir", destdir)

	keyv01, err := initRSAKeyByVer(destdir, "v0.1")
	if err != nil {
		panic(err)
	}

	allRsaKyes = append(allRsaKyes, keyv01)
}
