package randutils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	//"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
	"strings"
	//jwt "github.com/dgrijalva/jwt-go"
)

func RsaEncrypt(publicKey, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public is wrong")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(privateKey, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func RsaSign(privateKey, src []byte, hash crypto.Hash) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = priv.(*rsa.PrivateKey); !ok {
		return nil, errors.New("not private key type")
	}

	// pkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	// if err != nil {
	// 	return nil, err
	// }

	//log.Println(pkey)

	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	//log.Println(hashed)
	//log.Println(src)

	return rsa.SignPKCS1v15(nil, pkey, hash, hashed[:])

	//return jwt.SigningMethodRS256.Sign(string(src), pkey)

	//return rsa.SignPKCS1v15(rand.Reader, pkey, hash, hashed)
}

func RsaVerify(publickey, src []byte, sign []byte, hash crypto.Hash) error {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return errors.New("publickey key,  error!")
	}

	publi, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	publiky, ok := publi.(*rsa.PublicKey)
	if !ok {
		return errors.New("rsa.PublicKey wrong")
	}

	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)

	log.Println(hashed)

	return rsa.VerifyPKCS1v15(publiky, hash, hashed, sign)

	//return jwt.SigningMethodRS256.Verify(string(src), string(sign), publiky)

	//return rsa.VerifyPKCS1v15(publiky, hash, hashed, sign)
}

func EncodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

// Decode JWT specific base64url encoding with padding stripped
func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}
