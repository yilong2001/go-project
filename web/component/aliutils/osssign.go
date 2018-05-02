package aliutils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	//"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"hash"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

type headerSorter struct {
	Keys []string
	Vals []string
}

// Additional function for function SignHeader.
func newHeaderSorter(m map[string]string) *headerSorter {
	hs := &headerSorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]string, 0, len(m)),
	}

	for k, v := range m {
		hs.Keys = append(hs.Keys, k)
		hs.Vals = append(hs.Vals, v)
	}
	return hs
}

// Additional function for function SignHeader.
func (hs *headerSorter) Sort() {
	sort.Sort(hs)
}

// Additional function for function SignHeader.
func (hs *headerSorter) Len() int {
	return len(hs.Vals)
}

// Additional function for function SignHeader.
func (hs *headerSorter) Less(i, j int) bool {
	return bytes.Compare([]byte(hs.Keys[i]), []byte(hs.Keys[j])) < 0
}

// Additional function for function SignHeader.
func (hs *headerSorter) Swap(i, j int) {
	hs.Vals[i], hs.Vals[j] = hs.Vals[j], hs.Vals[i]
	hs.Keys[i], hs.Keys[j] = hs.Keys[j], hs.Keys[i]
}

func OssSign(method, contentMd5, contentType, date string, canonicalizedOSSHeaders, canonicalizedResource string) string {
	signStr := method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + canonicalizedResource
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(AliOss_GetRAMSecret()))
	io.WriteString(h, signStr)

	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signedStr

	// Get the final Authorization' string
	//authorizationStr := "OSS " + AliOss_GetRAMAccessKey() + ":" + signedStr
}

func OssSignHeader(req *http.Request, canonicalizedResource string) {
	// Find out the "x-oss-"'s address in this request'header
	temp := make(map[string]string)

	for k, v := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "x-oss-") {
			temp[strings.ToLower(k)] = v[0]
		}
	}
	hs := newHeaderSorter(temp)

	// Sort the temp by the Ascending Order
	hs.Sort()

	// Get the CanonicalizedOSSHeaders
	canonicalizedOSSHeaders := ""
	for i := range hs.Keys {
		canonicalizedOSSHeaders += hs.Keys[i] + ":" + hs.Vals[i] + "\n"
	}

	// Give other parameters values
	date := req.Header.Get("Date")
	contentType := req.Header.Get("Content-Type")
	contentMd5 := req.Header.Get("Content-MD5")

	// signStr := req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + canonicalizedResource
	// h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(conn.config.AccessKeySecret))
	// io.WriteString(h, signStr)
	// signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	signedStr := OssSign(req.Method, contentMd5, contentType, date, canonicalizedOSSHeaders, canonicalizedResource)

	// Get the final Authorization' string
	authorizationStr := "OSS " + AliOss_GetRAMAccessKey() + ":" + signedStr

	// Give the parameter "Authorization" value
	req.Header.Set("Authorization", authorizationStr)
}

func OssSignCdn(path, date string) string {
	src := path + "-" + date + "-0-0-" + AliOss_GetVodBj0CdnKey()
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func GetSignedUrl(path, exp string) map[string]string {
	//prefix := ""
	//wholepaht := prefix + path
	path2, err := url.QueryUnescape(path)
	if err != nil {
		path2 = path
	}

	resUri, pErr := url.Parse(path2)
	if pErr != nil {
		return map[string]string{}
	}
	//fmt.Println(resUri)
	f2 := resUri.EscapedPath()

	sign := OssSignCdn(f2, exp)
	log.Println(f2, path2, sign)

	return map[string]string{
		"Url":       AliOss_GetHostOfCdnZhiying() + f2,
		"Expires":   exp,
		"Signature": sign,
	}

	//urls := AliOss_GetHostOfZhiying() + f2 + "?OSSAccessKeyId=" + AliOss_GetRAMAccessKey() + "&Expires" + exp + "&Signature=" + sign2

	//return urls
}

func GetSignedUrl_OSS(path, exp string) map[string]string {
	//prefix := ""
	//wholepaht := prefix + path
	path2, err := url.QueryUnescape(path)
	if err != nil {
		path2 = path
	}

	canonicalizedResource := "/" + AliOss_GetResourceBucketOfZhiying() + path2
	sign := OssSign("GET", "", "", exp, "", canonicalizedResource)
	sign2 := url.QueryEscape(sign)

	resUri, pErr := url.Parse(path2)
	if pErr != nil {
		return map[string]string{}
	}
	//fmt.Println(resUri)
	f2 := resUri.EscapedPath()
	log.Println(f2, path2, sign)

	return map[string]string{
		"Url":            AliOss_GetHostOfZhiying() + path2,
		"OSSAccessKeyId": AliOss_GetRAMAccessKey(),
		"Expires":        exp,
		"Signature":      sign2,
	}

	//urls := AliOss_GetHostOfZhiying() + f2 + "?OSSAccessKeyId=" + AliOss_GetRAMAccessKey() + "&Expires" + exp + "&Signature=" + sign2

	//return urls
}
