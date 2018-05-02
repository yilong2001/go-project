package errcode

import (
	"fmt"
	"log"
)

const (
	Const_Error_Alert_High   = 0
	Const_Error_Alert_Medium = 1
)

func AlertLogOutput(level int, prefix string, info string) {
	log.Print(" *** Alert *** "+fmt.Sprint(level)+" *** ,", prefix, info)
}
