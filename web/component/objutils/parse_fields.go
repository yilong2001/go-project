package objutils

import (
	"github.com/gorilla/schema"
	"net/http"
)

func ParseObjectWithForm(obj interface{}, req *http.Request) error {
	req.ParseForm()
	decoder := schema.NewDecoder()
	// r.PostForm is a map of our POST form values
	err := decoder.Decode(obj, req.PostForm)
	return err
}
