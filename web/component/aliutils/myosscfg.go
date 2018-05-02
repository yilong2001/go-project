package aliutils

//python ./sts.py AssumeRole RoleArn=acs:ram::1574166049024626:role/ezhibao-app-user RoleSessionName=usr001 Policy='{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:ListObjects","oss:GetObject"],"Resource":["acs:oss:*:*:ezhibao-zhiying-out/build.sh"]}]}' DurationSeconds=3600 --id=LTAI2SD5eA0tz68H --secret=itraV0cdOQfJBDnlvzkll2DgULjPNP

const (
	Oss_RAMAccessKey    = ""
	Oss_RAMSecret       = ""
	Oss_RoleArn         = ""
	Oss_DurationSeconds = 3600

	Oss_Version   = "1"
	Oss_Statement = `[{"Effect":"Allow","Action":["oss:ListObjects","oss:GetObject"]`
)

func AliOss_GetHostOfCdnZhiying() string {
	return ""
}

func AliOss_GetHostOfZhiying() string {
	return AliOss_GetResourceBucketOfZhiying() + ".oss-cn-beijing.aliyuncs.com"
}

func AliOss_GetHostOf3rd(uid string) string {
	return AliOss_GetResourceBucketOf3rd(uid) + ".oss-cn-beijing.aliyuncs.com"
}

func AliOss_GetResourceBucketOfZhiying() string {
	return "ezhibao-zhiying" //"ezhibao-zhiying-out2"
}

func AliOss_GetResourceBucketOf3rd(uid string) string {
	return ""
}

func AliOss_GetRAMAccessKey() string {
	return Oss_RAMAccessKey
}

func AliOss_GetRAMSecret() string {
	return Oss_RAMSecret
}

func AliOss_GetRoleArn() string {
	return Oss_RoleArn
}

func AliOss_GetDurationSeonds() int {
	return Oss_DurationSeconds
}

func AliOss_GetVersion() string {
	return Oss_Version
}

func AliOss_GetStatement() string {
	return Oss_Statement
}

func AliOss_GetPolicyOfZhiying(path string) string {
	out := ""
	out += `{"Version":"` + AliOss_GetVersion() + `",`
	out += `"Statement":[{"Effect":"Allow","Action":["oss:ListObjects","oss:GetObject"],`
	out += `"Resource":["acs:oss:*:*:`
	out += AliOss_GetResourceBucketOfZhiying()
	out += `/` + path + `"]}]}`

	return out
}

func AliOss_GetPolicyOf3rd(uid string, path string) string {
	out := ""
	out += `{"Version":"` + AliOss_GetVersion() + `",`
	out += `"Statement":[{"Effect":"Allow","Action":["oss:ListObjects","oss:GetObject"],`
	out += `"Resource":["acs:oss:*:*:`
	out += AliOss_GetResourceBucketOf3rd(uid)
	out += `/` + path + `"]}]}`

	return out
}

func AliOss_GetVodBj1CdnKey() string {
	return ""
}

func AliOss_GetVodBj0CdnKey() string {
	return ""
}
