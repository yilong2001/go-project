package objutils

import (
	"errors"
	"fmt"
	//"log"
	"strings"
)

func charToStr(ch rune) string {
	return fmt.Sprintf("%c", ch)
}

func CamelToUnderLine(in string) (string, error) {
	orglen := len(in)
	if orglen != len(strings.Replace(in, " ", "", -1)) {
		return in, errors.New("has blank")
	}

	if orglen != len(strings.Replace(in, "_", "", -1)) {
		return in, errors.New("has _")
	}

	//does not conside number
	out := ""
	for id, ch := range in {
		//upper char
		if charToStr(ch) == strings.ToUpper(charToStr(ch)) {
			if id == 0 {
				out = out + strings.ToLower(charToStr(ch))
			} else {
				if ch < '0' || ch > '9' {
					out = out + "_" + strings.ToLower(charToStr(ch))
				} else {
					out = out + strings.ToLower(charToStr(ch))
				}
			}
		} else {
			out = out + strings.ToLower(charToStr(ch))
		}
	}

	return out, nil
}

func UnderLineToCamel(in string) (string, error) {
	orglen := len(in)
	if orglen != len(strings.Replace(in, " ", "", -1)) {
		return in, errors.New("has blank")
	}

	//does not conside number
	out := ""
	preIsUnderLine := false
	for id, ch := range in {
		if id == 0 {
			out = out + strings.ToUpper(charToStr(ch))
		} else {
			if preIsUnderLine && charToStr(ch) != "_" {
				out = out + strings.ToUpper(charToStr(ch))
				preIsUnderLine = false
			} else if charToStr(ch) != "_" {
				out = out + strings.ToLower(charToStr(ch))
			} else {

			}

			if charToStr(ch) == "_" {
				preIsUnderLine = true
			}
		}
	}

	return out, nil
}
