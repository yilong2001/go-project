package cfgutils

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type InitConfig struct {
	Mode string `xml:"mode"`
}

type WebApiConfig struct {
	AppName       string `xml:"appname"`
	HttpAddr      string `xml:"httpaddr"`
	HttpPort      string `xml:"httpport"`
	DbType        string `xml:"dbtype"`
	DbFile        string `xml:"dbfile"`
	DbUrl         string `xml:"dburl"`
	DbUser        string `xml:"dbuser"`
	DbPw          string `xml:"dbpw"`
	DbName        string `xml:"dbname"`
	SaltFile      string `xml:"saltfile"`
	HttpAddrIdGen string `xml:"httpaddridgen"`
	HttpPortIdGen string `xml:"httpportidgen"`
	UploadDir     string `xml:"uploaddir"`

	AdminHttpAddr string `xml:"adminhttpaddr"`
	AdminHttpPort string `xml:"adminhttpport"`

	WxpayH5Http   string `xml:"wxpayh5http"`
	WxpayH5Domain string `xml:"wxpayh5domain"`

	NewsV2HttpAddr string `xml:"newsv2httpaddr"`
	NewsV2HttpPort string `xml:"newsv2httpport"`
}

var globalInitConfig *InitConfig = &InitConfig{}
var globalWebApiConfig *WebApiConfig = &WebApiConfig{}

var globalNewsConfig *WebApiConfig = &WebApiConfig{}

var globalApiConfigPath string = ""
var globalNewsConfigPath string = ""

func GetInitConfig() *InitConfig {
	return globalInitConfig
}

func GetWebApiConfig() *WebApiConfig {
	if globalApiConfigPath == "" {
		return globalWebApiConfig
	}

	return GetWebApiConfigWithPath(globalApiConfigPath)
}

func GetNewsConfig() *WebApiConfig {
	if globalNewsConfigPath == "" {
		return globalNewsConfig
	}

	return GetWebApiConfigWithPath(globalNewsConfigPath)
}

func SetConfigPath(path string) {
	globalApiConfigPath = path
}

func GetWebApiConfigWithPath(configPath string) *WebApiConfig {
	var webApiConfig *WebApiConfig = &WebApiConfig{}

	log.Println("start load init config ... ")
	baseDir, _ := os.Getwd()
	log.Println("baseDir", baseDir)

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = xml.Unmarshal(content, webApiConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(webApiConfig)

	return webApiConfig
}

func init() {
	log.Println("start load init config ... ")
	baseDir, _ := os.Getwd()

	destdir := baseDir
	// if !strings.HasSuffix(baseDir, "uniapi") {
	// 	destDir = strings.Split(baseDir, "uniapi")[0]
	// 	destDir = destDir + "uniapi"
	// }

	if strings.Contains(baseDir, "uniapi/src/web") {
		destdir = strings.Split(baseDir, "uniapi")[0]
		destdir = destdir + "uniapi"
	}

	iniFileDir := destdir + "/bin/config/load/"
	if strings.HasSuffix(destdir, "bin") {
		iniFileDir = destdir + "/config/load/"
	}

	log.Println("iniFileDir", iniFileDir)

	content, err := ioutil.ReadFile(iniFileDir + "init.xml")
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(content, globalInitConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(globalInitConfig)

	log.Println("start load web api config ... ")

	content, err = ioutil.ReadFile(iniFileDir + "webapi_" + globalInitConfig.Mode + ".xml")
	if err != nil {
		log.Fatal(err)
	}
	err = xml.Unmarshal(content, globalWebApiConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(globalWebApiConfig)

	content1, err1 := ioutil.ReadFile(iniFileDir + "news_" + globalInitConfig.Mode + ".xml")
	if err1 != nil {
		log.Fatal(err1)
	}
	err1 = xml.Unmarshal(content1, globalNewsConfig)
	if err1 != nil {
		log.Fatal(err1)
	}
	log.Println(globalNewsConfig)
}
