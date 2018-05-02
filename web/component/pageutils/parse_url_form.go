package pageutils

import (
	//"log"
	"net/http"
	//"net/url"
	"strings"
)

func Parse_URL_Form(req *http.Request) map[string]string {
	out := make(map[string]string)

	req.ParseForm()

	// log.Println(req.URL.String())

	// u, err := url.Parse(req.URL.String())
	// if err != nil {
	// 	panic(err)
	// }

	// m, _ := url.ParseQuery(u.RawQuery)
	// log.Println(m)

	out["page"] = req.Form.Get("page")
	out["pageSize"] = req.Form.Get("pageSize")
	out["sort"] = req.Form.Get("sort")
	out["fields"] = req.Form.Get("fields")
	out["url"] = req.URL.String()
	out["surl"] = strings.Split(req.URL.String(), "?")[0]

	return out
}
