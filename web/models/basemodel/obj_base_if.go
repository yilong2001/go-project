package basemodel

type ObjectUtilBaseIf interface {
	//IsUserIdValid() bool
	//GetUserId() int
	GetUniqId() int
	GetUniqIdName() string
	//GetFieldsWithNotDefaultValue(url string) (map[string]interface{}, map[string]interface{})
	GetWholeFields() (map[string]interface{}, map[string]interface{})
	GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{})
	GetFieldsWithSpecs(skips []string) (map[string]interface{}, map[string]interface{})
}
