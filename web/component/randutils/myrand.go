package randutils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// 随机字符串
func KRand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func KRandNum(size int) []byte {
	return KRand(size, KC_RAND_KIND_NUM)
}

func KRandLower(size int) []byte {
	return KRand(size, KC_RAND_KIND_LOWER)
}

func KRandUpper(size int) []byte {
	return KRand(size, KC_RAND_KIND_UPPER)
}

func KRandAll(size int) []byte {
	return KRand(size, KC_RAND_KIND_ALL)
}

func KRandMd5() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	sessionId := MyMd5(MyMd5(strconv.FormatInt(nano, 10)) + MyMd5(strconv.FormatInt(rndNum, 10)))
	log.Println(sessionId)

	return sessionId
}

func MyMd5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

const (
	base64Table = "LM+N/OP26ABijklmCK78vDE345FJwxyzQRSTXYZaUVWbcdefghnoGHIpqrstu019"
)

var coder = base64.NewEncoding(base64Table)

func KBase64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func KBase64EncodeString(src []byte) string {
	return coder.EncodeToString(src)
}

func KBase64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

func HashUID(uid int) string {
	log.Println(uid)

	uidBts := []byte(fmt.Sprint(uid))

	//log.Println(uidBts, len(uidBts))

	rndBts := KRandAll(128 - len(uidBts))
	//log.Println(rndBts, len(rndBts))

	allBts := make([]byte, 0)
	for _, bt := range uidBts {
		allBts = append(allBts, bt)
	}

	for _, bt := range rndBts {
		allBts = append(allBts, bt)
	}

	sl1 := allBts[0:8]
	sl2 := allBts[8:128]

	//log.Println(allBts, len(allBts))

	hm5 := md5.New()

	hm5.Write(sl2)

	bts := hm5.Sum(sl1)

	b64 := KBase64EncodeString(bts)
	he := hex.EncodeToString(bts)

	log.Println(b64)

	log.Println(he)

	return b64
}
