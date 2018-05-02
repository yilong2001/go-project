package uploads

import (
	//"errors"
	"io"
	"log"
	"net/http"
	//"strconv"
	//"time"
	"web/component/fileutils"
)

func UploadFile(req *http.Request, destDirType string, useridStr string) (error, string) {
	log.Println("parsing form")

	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println("ParseMultipartForm = ", err)
		return err, ""
	}

	//files := req.MultipartForm.File["files"]
	//file, err := files[i].Open()
	//for i, _ := range files {
	//}

	log.Println("getting handle to file ... ")

	file, handler, err1 := req.FormFile("upFile")
	if err1 != nil {
		log.Println("get form file failed, ", err1)
		return err1, ""
	}

	defer file.Close()

	srcFileName := handler.Filename

	log.Println("creating destination file ... ")

	//dst, err := os.OpenFile(getUserPortraitDir(params["UserId"])+files[i].Filename, os.O_WRONLY|os.O_CREATE, 0666)
	dst, err2 := fileutils.CreateFile(fileutils.GetUpladFileDir(destDirType, useridStr), srcFileName)
	defer dst.Close()
	if err2 != nil {
		log.Println("create file failed, ", err2, srcFileName)
		return err2, ""
	}

	log.Println("copying the uploaded file to the destination file ... ")

	if _, err := io.Copy(dst, file); err != nil {
		log.Println("copy failed , ", err)
		return err, ""
	}

	outFilePath := fileutils.GetUpladFileRelativeDir(destDirType, useridStr) + srcFileName

	return nil, outFilePath
}
