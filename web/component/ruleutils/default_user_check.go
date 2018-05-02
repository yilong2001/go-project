package ruleutils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type DefaultPhoneCheckRule struct{}
type DefaultLoginNameCheckRule struct{}
type DefaultEmailCheckRule struct{}
type DefaultUserIdCheckRule struct{}
type DefaultIndustryCheckRule struct{}

type DefaultSexCheckRule struct{}
type DefaultAgeCheckRule struct{}
type DefaultCityCheckRule struct{}

func isNoBlankAndNoSpace(in string, minlen int) bool {
	if len(in) < minlen {
		return false
	}

	if strings.Contains(in, " ") {
		return false
	}

	return true
}

func (this *DefaultPhoneCheckRule) Check(in string) bool {
	log.Println("DefaultPhoneCheckRule:" + in)
	if !isNoBlankAndNoSpace(in, 7) {
		return false
	}

	reg := regexp.MustCompile(`\d{7,15}`)

	return reg.MatchString(in)
}

func (this *DefaultLoginNameCheckRule) Check(in string) bool {
	if !isNoBlankAndNoSpace(in, 8) {
		return false
	}

	return true
}

func (this *DefaultEmailCheckRule) Check(in string) bool {
	if !isNoBlankAndNoSpace(in, 8) {
		return false
	}

	return true
}

func (this *DefaultUserIdCheckRule) Check(in string) bool {
	log.Println("DefaultUserIdCheckRule:" + in)
	log.Println("In len is : %", len(in))

	if len(in) < 5 {
		return false
	}

	i, err := strconv.ParseInt(in, 10, 32)
	if err != nil {
		return false
	}

	if i < 10000 {
		return false
	}

	return true
}

func (this *DefaultIndustryCheckRule) Check(in string) bool {
	log.Println("DefaultIndustryCheckRule:" + in)
	_, err := strconv.ParseInt(in, 10, 32)
	if err != nil {
		return false
	}

	//if i < 1001000 {
	//	return false
	//}

	return true
}

func (this *DefaultSexCheckRule) Check(in string) bool {
	log.Println("DefaultSexCheckRule:" + in)
	i, err := strconv.ParseInt(in, 10, 32)
	if err != nil {
		return false
	}

	if i < 0 || i > 2 {
		return false
	}

	return true
}

func (this *DefaultAgeCheckRule) Check(in string) bool {
	i, err := strconv.ParseInt(in, 10, 32)
	if err != nil {
		return false
	}

	if i < 0 || i > 8 {
		return false
	}

	return true
}

func (this *DefaultCityCheckRule) Check(in string) bool {
	i, err := strconv.ParseInt(in, 10, 32)
	if err != nil {
		return false
	}

	if i < 110100 {
		return false
	}

	return true
}
