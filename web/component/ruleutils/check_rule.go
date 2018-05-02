package ruleutils

import ()

type CheckRule interface {
	Check(in string) bool
}
