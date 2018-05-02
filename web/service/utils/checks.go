package utils

import (
	"fmt"
	//"log"
	"strings"
	"web/component/ruleutils"
)

var globalCheckRulesMap map[string]ruleutils.CheckRule

func init() {
	globalCheckRulesMap = make(map[string]ruleutils.CheckRule)

	globalCheckRulesMap["phone"] = &ruleutils.DefaultPhoneCheckRule{}
	globalCheckRulesMap["login_name"] = &ruleutils.DefaultLoginNameCheckRule{}
	globalCheckRulesMap["email"] = &ruleutils.DefaultEmailCheckRule{}

	globalCheckRulesMap["user_id"] = &ruleutils.DefaultUserIdCheckRule{}
	globalCheckRulesMap["job_id"] = &ruleutils.DefaultUserIdCheckRule{}
	globalCheckRulesMap["service_id"] = &ruleutils.DefaultUserIdCheckRule{}

	globalCheckRulesMap["order_id"] = &ruleutils.DefaultUserIdCheckRule{}
	globalCheckRulesMap["order_sub_id"] = &ruleutils.DefaultUserIdCheckRule{}

	globalCheckRulesMap["firm_id"] = &ruleutils.DefaultUserIdCheckRule{}
	globalCheckRulesMap["course_id"] = &ruleutils.DefaultUserIdCheckRule{}
	globalCheckRulesMap["comment_id"] = &ruleutils.DefaultUserIdCheckRule{}
	globalCheckRulesMap["tag_id"] = &ruleutils.DefaultUserIdCheckRule{}

	globalCheckRulesMap["industry"] = &ruleutils.DefaultIndustryCheckRule{}
	globalCheckRulesMap["sex"] = &ruleutils.DefaultSexCheckRule{}
	globalCheckRulesMap["age"] = &ruleutils.DefaultAgeCheckRule{}
	globalCheckRulesMap["city"] = &ruleutils.DefaultCityCheckRule{}
}

func GetCheckRules(attr string) ruleutils.CheckRule {
	return globalCheckRulesMap[attr]
}

func IsFieldCorrectWithRule(attr string, in string) bool {
	rule := GetCheckRules(attr)
	if rule == nil {
		return true
	}

	return rule.Check(in)
}

func IsFieldsValueOk(fieldArrs map[string]interface{}) (bool, string) {
	res := true
	fns := ""
	for fn, fd := range fieldArrs {
		//log.Println("check : " + fn)
		if !IsFieldCorrectWithRule(strings.ToLower(fn), fmt.Sprintf("%v", fd)) {
			res = false
			fns = fns + fn + ", "
		}
	}

	return res, fns
}
