package randutils

import (
	"crypto/md5"
	"encoding/hex"
	//"log"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

func BuildRawMd5String(pw, phone string) string {
	str1 := BuildMd5PWPhoneStringV2(pw, "")
	str2 := BuildMd5PWPhoneStringV2(str1, phone)

	return str2
}

func BuildMd5PWPhoneString(md5pw, phone string) string {
	// hsm5 := md5.New()

	// _, err := hsm5.Write([]byte(md5pw + phone))
	// if err != nil {
	// 	panic(err)
	// }

	// abts2 := hsm5.Sum(nil)
	//log.Println(bts2)

	//return string(abts2)

	t1 := md5.New()
	io.WriteString(t1, md5pw+phone)

	return hex.EncodeToString(t1.Sum(nil))
}

func BuildMd5PWPhoneStringV2(md5pw, phone string) string {
	hsm5 := md5.New()

	_, err := hsm5.Write([]byte(md5pw + phone))
	if err != nil {
		panic(err)
	}

	abts2 := hsm5.Sum(nil)
	//log.Println(abts2)

	return hex.EncodeToString(abts2)
}

func BuildMd5PWPhoneStringV64(md5pw, phone string) string {
	hsm5 := md5.New()

	_, err := hsm5.Write([]byte(md5pw + phone))
	if err != nil {
		panic(err)
	}

	abts2 := hsm5.Sum(nil)
	//log.Println(abts2)

	return base64.StdEncoding.EncodeToString(abts2)
}

func BuildMD5OrderOutId(info string, orderid int) string {
	md5 := BuildMd5PWPhoneStringV64(info, "")
	md5 = strings.ToUpper(md5)

	pluss := fmt.Sprint(orderid)[0:1]
	decss := fmt.Sprint(orderid * 3)[0:1]
	xcss := fmt.Sprint(orderid * 7)[0:1]

	md5 = strings.Replace(md5, "+", pluss, -1)
	md5 = strings.Replace(md5, "-", decss, -1)
	md5 = strings.Replace(md5, "/", xcss, -1)

	out := md5[0:4] + "-" + md5[4:8] + "-" + md5[8:12] + "-" + md5[12:16]
	return out
}
