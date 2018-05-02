package fileutils

import (
	"log"
	"os"
	//"web/component/cfgutils"
)

func GetPortraitRelativeDir(userid string) string {
	return "/uploads/" + userid + "/portrait/"
}

func GetLogoRelativeDir(userid string) string {
	return "/uploads/" + userid + "/logo/"
}

func GetPortraitDir(userid string) string {
	baseDir, _ := os.Getwd()
	log.Println("baseDir", baseDir)

	return baseDir + "/public" + GetPortraitRelativeDir(userid)
}

func GetUpladFileRelativeDir(destType, userid string) string {
	return "/uploads/" + destType + "/" + userid + "/"
}

func GetUpladFileDir(destType string, userid string) string {
	baseDir, _ := os.Getwd()
	log.Println("baseDir", baseDir)

	return baseDir + "/public" + GetUpladFileRelativeDir(destType, userid)
}

func GetLogoDir(userid string) string {
	baseDir, _ := os.Getwd()
	log.Println("baseDir", baseDir)

	return baseDir + "/public" + GetLogoRelativeDir(userid)
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func CreateFile(dir string, name string) (*os.File, error) {
	src := dir + name
	if IsExist(src) {
		return os.OpenFile(src, os.O_WRONLY|os.O_CREATE, 0666)
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		if os.IsPermission(err) {
			log.Println("你不够权限创建文件")
		}
		return nil, err
	}

	return os.OpenFile(src, os.O_WRONLY|os.O_CREATE, 0666)
}
