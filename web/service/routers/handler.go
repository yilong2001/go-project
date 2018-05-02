package routers

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var globalClassMartini *martini.ClassicMartini = &martini.ClassicMartini{}

func init() {
	globalClassMartini = martini.Classic()
	globalClassMartini.Use(render.Renderer())
}

func GetRouterHandler() *martini.ClassicMartini {
	return globalClassMartini
}
