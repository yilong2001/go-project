package main

import (
	"time"
	"web/component/idutils"
)

func main() {
	go func() {
		count := 10000
		cur := 1
		for {
			if cur > count {
				break
			}

			(idutils.GetId("user_id"))
		}
	}()

	go func() {
		count := 10000
		cur := 1
		for {
			if cur > count {
				break
			}

			(idutils.GetId("user_id"))
		}
	}()

	time.Sleep(60 * time.Second)
}
