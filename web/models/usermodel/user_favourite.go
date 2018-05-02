package usermodel

import (
	"web/component/objutils"
)

type UserFavouriteInfo struct {
	FavouriteId int
	UserId      int
	DestId      int `schema:"DestId"`
	DestType    int `schema:"DestType"`
	CreateTime  string
}

func (this *UserFavouriteInfo) GetUniqId() int {
	return this.FavouriteId
}

func (this *UserFavouriteInfo) GetUniqIdName() string {
	return "FavouriteId"
}
func (this *UserFavouriteInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *UserFavouriteInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *UserFavouriteInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
