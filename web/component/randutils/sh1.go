package randutils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
)

func BuildHSA1String(info string, key []byte) string {
	//sha1
	h := sha1.New()
	io.WriteString(h, info)
	log.Printf("%x\n", h.Sum(nil))

	//hmac ,use sha1
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(info))

	out := hex.EncodeToString(mac.Sum(nil))
	log.Println(out)

	return out
}
