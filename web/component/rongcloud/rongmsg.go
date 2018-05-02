package rongcloud

import (
	"web/component/objutils"
)

//{"content":{"messageName":"TextMessage","content":"where a u","extra":"附加信息"},
//"conversationType":1,"objectName":"RC:TxtMsg","messageDirection":1,
//"senderUserId":"13099101011","sentStatus":30,"sentTime":1467803986977,
//"targetId":"13099101012","messageType":"TextMessage",
//"messageUId":"5ARV-RRC4-44QV-N7KD"}

//{"content":{"messageName":"TextMessage","content":"where a u","extra":"附加信息"},
//"conversationType":1,"objectName":"RC:TxtMsg","messageDirection":2,
//"messageId":"1_1106342","receivedTime":1467803986194,"senderUserId":"13099101011",
//"sentTime":1467803986977,"targetId":"13099101011",
//"messageType":"TextMessage","messageUId":"5ARV-RRC4-44QV-N7KD"}

type RongCloudIMSentMsg struct {
	MessageId        int
	Title            string `shema:"title"`
	ImageUri         string `schema:"imageUri"`
	Url              string `schema:"url"`
	MessageName      string `schema:"messageName"`
	Content          string `schema:"content"`
	Extra            string `schema:"extra"`
	ConversationType int    `schema:"conversationType"`
	ObjectName       string `schema:"objectName"`
	MessageDirection int    `schema:"messageDirection"`
	SentStatus       int    `schema:"sentStatus"`
	SentTime         int64  `schema:"sentTime"`
	SenderUserId     string `schema:"senderUserId"`
	TargetId         string `schema:"targetId"`
	MessageType      string `schema:"messageType"`
	MessageUid       string `schema:"messageUid"`
	IsRead           int    `schema:"isRead"`
}

func (this *RongCloudIMSentMsg) GetUniqId() int {
	return int(this.MessageId)
}

func (this *RongCloudIMSentMsg) GetUniqIdName() string {
	return "MessageId"
}

func (this *RongCloudIMSentMsg) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *RongCloudIMSentMsg) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *RongCloudIMSentMsg) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

type RongCloudIMRcvMsg struct {
	MessageId        int64
	MessageName      string `schema:"messageName"`
	Content          string `schema:"content"`
	Extra            string `schema:"extra"`
	ConversationType int    `schema:"conversationType"`
	ObjectName       string `schema:"objectName"`
	MessageDirection int    `schema:"messageDirection"`
	//MessageId        string `schema:"messageId"`
	ReceivedTime int64 `schema:"receivedTime"`
	//SentStatus       int    `schema:"sentStatus"`
	SentTime     string `schema:"sentTime"`
	SenderUserId string `schema:"senderUserId"`
	TargetId     string `schema:"targetId"`
	MessageType  string `schema:"messageType"`
	MessageUid   string `schema:"messageUId"`
}
