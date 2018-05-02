package ctrlbase

import (
	//"database/sql"
	//"log"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

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

type FavouriteUpdateBaseController struct {
	CtrlBaseController
	UniqId string
	DestId interface{}

	AddField string
}

func (this *FavouriteUpdateBaseController) CheckParams(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func (this *FavouriteUpdateBaseController) ExInfoInit(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *FavouriteUpdateBaseController) CondCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {
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

func (this *FavouriteUpdateBaseController) CalcAddFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	//fds := map[string]string{"favourite_num": " favourite_num = favourite_num + 1"}
	fds := map[string]string{this.AddField: this.AddField + "  = " + this.AddField + " + 1"}

	return fds
}

func (this *FavouriteUpdateBaseController) CalcDecFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	//fds := map[string]string{"favourite_num": " favourite_num = favourite_num - 1"}
	fds := map[string]string{this.AddField: this.AddField + "  = " + this.AddField + " - 1"}

	return fds
}
