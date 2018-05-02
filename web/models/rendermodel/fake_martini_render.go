package rendermodel

import (
	"github.com/martini-contrib/render"
	"html/template"
	"net/http"
)

type FakeMrtiniRender struct {
	req    *http.Request
	val    interface{}
	status int
}

func (this *FakeMrtiniRender) GetVal() interface{} {
	return this.val
}

func (this *FakeMrtiniRender) GetStatus() int {
	return this.status
}

func (this *FakeMrtiniRender) JSON(status int, v interface{}) {
	this.val = v
	this.status = status
}

func (this *FakeMrtiniRender) HTML(status int, name string, v interface{}, htmlOpt ...render.HTMLOptions) {

}

func (this *FakeMrtiniRender) XML(status int, v interface{}) {
	this.val = v
	this.status = status
}

func (this *FakeMrtiniRender) Status(status int) {
	this.status = status
}

func (this *FakeMrtiniRender) Data(status int, v []byte) {}

func (this *FakeMrtiniRender) Text(status int, v string) {}

func (this *FakeMrtiniRender) Error(status int) {
	this.status = status
}

func (this *FakeMrtiniRender) Redirect(location string, status ...int) {}

func (this *FakeMrtiniRender) Template() *template.Template {
	return nil
}

func (this *FakeMrtiniRender) Header() http.Header {
	if this.req == nil {
		return nil
	}

	return this.req.Header
}
