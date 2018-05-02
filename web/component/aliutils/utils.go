package aliutils

import (
	//"crypto/sha1"
	//"encoding/base64"
	//"encoding/hex"
	//"encoding/json"
	"fmt"
	"log"
	"sort"
	//"strconv"
	//"strings"
	//"web/component/keyutils"
	//"web/component/randutils"
)

func genAlipayUrlString(mapBody map[string]interface{}) string {
	var signStrings string

	index := 0
	for k, v := range mapBody {
		log.Println("k=", k, "v =", v)
		value := fmt.Sprintf("%v", v)
		if value != "" {
			signStrings = signStrings + k + "=" + value
		}
		//最后一项后面不要&
		if index < len(mapBody)-1 {
			signStrings = signStrings + "&"
		}
		index++
	}

	return signStrings
}

func genAlipaySignString(mapBody map[string]interface{}) (sign string, err error) {
	sorted_keys := make([]string, 0)
	for k, _ := range mapBody {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	var signStrings string

	index := 0
	for _, k := range sorted_keys {
		//log.Println("k=", k, "v =", mapBody[k])
		value := fmt.Sprintf("%v", mapBody[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value
		}
		//最后一项后面不要&
		if index < len(sorted_keys)-1 {
			signStrings = signStrings + "&"
		}
		index++
	}

	return signStrings, nil
}
