package idutils

import (
	"testing"
	"time"
)

func Test_IdGenSync(t *testing.T) {
	go func() {
		count := 10000
		cur := 1
		for {
			if cur > count {
				break
			}

			(GetId("user_id"))
		}
	}()

	go func() {
		count := 10000
		cur := 1
		for {
			if cur > count {
				break
			}

			(GetId("user_id"))
		}
	}()

	time.Sleep(30 * time.Second)
}
