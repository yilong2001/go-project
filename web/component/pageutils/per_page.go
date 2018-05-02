package pageutils

import (
	"log"
	//"github.com/go-martini/martini" //
	"net/url"
	"strconv"
)

const (
	Default_Size_Per_Page      int    = 3
	Default_Page_Size_Str_Flag string = "pageSize"
	Default_Page_Num_Str_Flag  string = "page"
)

func GetRequestPageSize(table string, params url.Values) int {
	log.Println("GetRequestPageNo, ", params.Get(Default_Page_Size_Str_Flag))

	ps := Default_Size_Per_Page
	if params.Get(Default_Page_Size_Str_Flag) != "" {
		ct, err := strconv.ParseInt(params.Get(Default_Page_Size_Str_Flag), 10, 32)
		if err != nil {
			log.Println(err)
		} else {
			ps = int(ct)
		}
	}

	log.Println(table+": per page is : ", ps)
	return ps
}

func GetRequestPageNo(table string, params url.Values) int64 {
	log.Println("GetRequestPageNo, ", params.Get(Default_Page_Num_Str_Flag))
	pn := int64(1)
	if params.Get(Default_Page_Num_Str_Flag) != "" {
		ct, err := strconv.ParseInt(params.Get(Default_Page_Num_Str_Flag), 10, 32)
		if err != nil {
			log.Println(err)
		} else {
			pn = (ct)
		}
	}

	log.Println(table+": req page num is : ", pn)
	return pn
}
