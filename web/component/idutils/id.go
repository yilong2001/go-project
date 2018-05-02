package idutils

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetId(objectName string) int {
	resp, err := http.Get("http://localhost:8190" + "/?name=" + objectName)
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	if resp.StatusCode != 200 {
		log.Panic(string(body))
	}

	nums := string(body)

	log.Println(nums)

	numsarr := strings.Split(nums, ",")

	res, err := strconv.ParseInt(numsarr[0], 10, 32)
	if err != nil {
		log.Panic(err)
	}

	return int(res)
}
