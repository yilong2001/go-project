package imgutils

import (
	"testing"
)

func Test_RSA(t *testing.T) {
	path := "/Users/yilongli/Pictures/20160822_185232.jpg"
	npath := "/Users/yilongli/Pictures/20160822_185232_1.jpg"
	//path := "/ebang/web/const/uniapi/public/uploads/12520/portrait/20160822_185232.jpg"
	//npath := "/ebang/web/const/uniapi/public/uploads/12520/portrait/20160822_185232_1.jpg"
	err := ResizeImg(path, npath, 200, 200, true)
	if err != nil {
		t.Fail()
	}
}
