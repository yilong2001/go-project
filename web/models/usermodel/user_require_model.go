package usermodel

import (
	"web/component/objutils"
)

type UserRequireInfo struct {
	RequireId     int
	RequireUserId int
	RequireType   string `schema:"RequireType"`
	Address       string `schema:"Address"`
	Direction     string `schema:"Direction"`
	Profession    string `schema:"Profession"`
	Title         string `schema:"Title"`
	ExpectArrange string `schema:"ExpectArrange"`
	Reason        string `schema:"Reason"`
	Object        string `schema:"Object"`
	Goal          string `schema:"Goal"`
	Detail        string `schema:"Detail"`
	Duration      string `schema:"Duration"`
	PeopleNum     string `schema:"PeopleNum"`
	Firm          string `schema:"Firm"`
	Contact       string `schema:"Contact"`
	Phone         string `schema:"Phone"`
	Email         string `schema:"Email"`
	CreateTime    string
}

func (this *UserRequireInfo) GetUniqId() int {
	return this.RequireId
}

func (this *UserRequireInfo) GetUniqIdName() string {
	return "RequireId"
}

func (this *UserRequireInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *UserRequireInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *UserRequireInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
