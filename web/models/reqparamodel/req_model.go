package reqparamodel

import (
	"github.com/go-martini/martini"
	"web/models/basemodel"
	//jwt "github.com/dgrijalva/jwt-go"
	"net/url"
)

type HttpReqParams struct {
	TokenParams  map[string]string
	RouterParams map[string]string
	URLParams    url.Values
	Url          string
	ShortUrl     string
	PostFields   []string
}

func (this *HttpReqParams) MergeMartiniParams(params martini.Params) {
	for _, idname := range basemodel.Default_All_UniqId_Names {
		this.RouterParams[idname] = params[idname]
	}

	if this.TokenParams["UserId"] != "" {
		this.RouterParams["UserId"] = this.TokenParams["UserId"]
	}
}

func NewHttpReqParams() *HttpReqParams {
	return &HttpReqParams{
		TokenParams:  make(map[string]string),
		RouterParams: make(map[string]string),
		PostFields:   make([]string, 0),
		//URLParams:    make(map[string][]string),
	}
}
