package randutils

import (
	"log"
	"testing"
)

func Test_Rand(t *testing.T) {
	log.Println(string(KRandAll(256)))
}
