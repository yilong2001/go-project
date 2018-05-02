package ctrlbase

import (
	//"database/sql"
	//"log"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"fmt"
	"log"
	//"net/http"
	//"strconv"
	//"strings"
	"encoding/json"
	//"errors"
	//"time"
	//"reflect"

	//"web/component/cfgutils"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/condmodel"
	"web/models/reqparamodel"
	//"web/models/servemodel"
	//"web/models/usermodel"
	//"web/service/routers"
	//"web/service/utils"
)

type CommentUpdateBaseController struct {
	CtrlBaseController
	CommentLevel int
	UniqId       string
	DestId       int

	UpFieldName  string
	AvgFieldName string
}

func (this *CommentUpdateBaseController) CheckParams(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func (this *CommentUpdateBaseController) ExInfoInit(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *CommentUpdateBaseController) CondCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {
	compser := condmodel.NewCondComposerLinker(" or ")

	where := map[string]interface{}{this.UniqId: this.DestId}
	rule := map[string]string{this.UniqId: " = "}
	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	compser.SetItem(compSub)
	compser.SetNext(nil)

	ji, _ := json.Marshal(compser)
	log.Println("update "+this.UniqId+" num", string(ji))

	return compser
}

func (this *CommentUpdateBaseController) CalcUpFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {
	fds := map[string]string{}

	fds["avg_star"] = " avg_star = (comment_num * avg_star + " + fmt.Sprint(this.CommentLevel) + ")/(comment_num+1)"

	fds["comment_num"] = " comment_num = comment_num + 1 "

	return fds
}

func (this *CommentUpdateBaseController) CalcServedNumFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	fds := map[string]string{"served_num": " served_num = served_num + 1 "}

	return fds
}

func (this *CommentUpdateBaseController) CalcUpFieldsEx(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {
	fds := map[string]string{}

	if this.AvgFieldName != "" {
		fds[this.AvgFieldName] = " " + this.AvgFieldName + " = (" + this.UpFieldName + " *  " + this.AvgFieldName + " + " + fmt.Sprint(this.CommentLevel) + ")/(" + this.UpFieldName + "+1)"
	}

	fds[this.UpFieldName] = " " + this.UpFieldName + " = " + this.UpFieldName + " + 1 "

	return fds
}
