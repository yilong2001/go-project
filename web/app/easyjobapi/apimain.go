package main

import (
	"log"
	"os"
	//"net/http"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"web/component/cfgutils"
	_ "web/service/orders"
	_ "web/service/others"
	_ "web/service/platform/admins"
	_ "web/service/resources"
	"web/service/routers"
	_ "web/service/serves"
	_ "web/service/servicers"
	_ "web/service/thirdpay"
	"web/service/tokens"
	_ "web/service/wx"
)

func main() {
	rg_num := len(os.Args)
	if rg_num > 1 {
		cfgutils.SetConfigPath(os.Args[1])
	}

	log.Println("hello")

	m := routers.GetRouterHandler()

	//httpReqParams := &reqparamodel.HttpReqParams{}
	//m.Map(httpReqParams)

	//m.RunOnAddr(cfgutils.GetWebApiConfig().HttpPort)
	//http.Handle("career.atayun.com/", m)
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://*.atayun.com", "http://*.atayunlocal.com", "http://*.atayun.com:3000", "http://*.atayunlocal.com:3000", "http://localhost:8990", "http://192.168.20.*:8990", "http://localhost:3000", "http://localhost:9190", "http://192.168.20.*:3000", "http://192.168.1.*:3000", "http://192.168.191.*:3000", "http://192.168.10.*:3000", "http://192.168.20.*:3000", "http://*.zhieasy.com", "http://*.zhieasy.com:9190", "http://*.zhieasy.com:80"},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "X-ACCESS-TOKEN"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	m.Use(martini.Static("public"))

	m.Use(tokens.TokenHeaderHandler)

	m.RunOnAddr(cfgutils.GetWebApiConfig().HttpPort)
}
