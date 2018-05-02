package pageutils

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"web/component/sqlutils"
	"web/models/condmodel"
)

type Paginator struct {
	FirstPage   int64
	PrePage     int64
	CurrentPage int64
	NextPage    int64
	TotalPage   int64
	TotalCounts int64
	PageSize    int
	//UniqIdName string
}

func (this *Paginator) GetSQLLimit() int {
	return this.PageSize
}

func (this *Paginator) GetSQLOffset() int64 {
	return (this.CurrentPage - 1) * int64(this.PageSize)
}

func (this *Paginator) MergePaingCond(db *sql.DB, uniqIdName, tableName string, whereCond map[string]interface{}, ruleCond map[string]string, orderbys string, exCondVal []interface{}) {

	log.Println(uniqIdName, whereCond)
	for fn, _ := range whereCond {
		if fn == uniqIdName {
			return
		}
	}

	if this.isOrderByUnUniqName(uniqIdName, orderbys) {
		return
	}

	rl := this.getCompRuleByOrderby(uniqIdName, orderbys)
	uid := this.fetchLimitId(db, uniqIdName, tableName, whereCond, ruleCond, orderbys, exCondVal)

	whereCond[uniqIdName] = fmt.Sprintf(uniqIdName+" "+rl+" %d", uid)

	ruleCond[uniqIdName] = rl

	return
}

func (this *Paginator) isOrderByUnUniqName(field string, orderbys string) bool {
	ascCount := strings.Count(orderbys, " asc")
	descCount := strings.Count(orderbys, " desc")

	fieldCount1 := strings.Count(orderbys, " "+field+" ")
	fieldCount2 := strings.Count(orderbys, ","+field+" ")

	log.Println(orderbys, field, ascCount, descCount, fieldCount1, fieldCount2)

	if (ascCount + descCount) == (fieldCount1 + fieldCount2) {
		return false
	}

	return true
}

func (this *Paginator) isOrderByAsc(field string, orderbys string) bool {
	if strings.Contains(orderbys, field+" asc") {
		return true
	}

	if strings.Contains(orderbys, field+" desc") {
		return false
	}

	ascid := strings.Index(orderbys, " asc")
	descid := strings.Index(orderbys, " desc")
	log.Println(ascid, descid)
	if descid == ascid && descid == -1 {
		return false
	}

	if ascid < 0 && descid >= 0 {
		return false
	}

	if ascid >= 0 && descid >= 0 && descid < ascid {
		return false
	}

	return true
}

func (this *Paginator) getCompRuleByOrderby(field string, orderbys string) string {
	if this.isOrderByAsc(field, orderbys) {
		return " >= "
	}

	return " <= "
}

func (this *Paginator) MergePaingWithComposerLinker(db *sql.DB, uniqIdName, tableName string, conds *condmodel.CondComposerLinker, orderbys string, exWhere []interface{}) *condmodel.CondComposerLinker {
	//log.Println(conds)

	found := false
	next := conds
	for {
		if next == nil {
			break
		}
		if next.Item != nil {
			for fn, _ := range next.Item.Where {
				if fn == uniqIdName {
					found = true
					break
				}
			}
		}

		next = next.Next
	}

	if found {
		return conds
	}

	if this.isOrderByUnUniqName(uniqIdName, orderbys) {
		return conds
	}

	uid := this.fetchLimitIdWithComposer(db, uniqIdName, tableName, conds, orderbys, exWhere)

	rl := this.getCompRuleByOrderby(uniqIdName, orderbys)

	subsqls := fmt.Sprintf(uniqIdName+" "+rl+" %d", uid)

	where := map[string]interface{}{uniqIdName: subsqls}
	rule := map[string]string{uniqIdName: rl}

	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	composer := condmodel.NewCondComposerLinker("and")
	composer.SetItem(compSub)
	composer.SetNext(conds)

	//conds.AddItem(compSub)

	return composer
}

func (this *Paginator) fetchLimitIdWithComposer(db *sql.DB, uniqIdName, tableName string, conds *condmodel.CondComposerLinker, orderby string, exWhere []interface{}) int {
	log.Println(uniqIdName, conds)

	uniqId := -1
	selFields := make(map[string]interface{})
	selFields[uniqIdName] = &uniqId

	msqls, selargs, whereargs := sqlutils.Sqls_CompSelectWithComposerLinker(tableName, selFields, conds)

	whereargs = append(exWhere, whereargs...)

	//" order by " + uniqIdName + " desc "
	msqls = msqls + orderby + " limit " + fmt.Sprint(this.GetSQLOffset()) + ", 1 "

	log.Println("FetchLimitId: ", msqls)

	err := sqlutils.Sqls_Do_QueryRowAndScan(db, msqls, selargs, whereargs)
	if err != nil {
		log.Println("fetchLimitIdWithComposer", err)
		return -1
	}

	log.Println(uniqIdName, uniqId)

	return uniqId
}

func (this *Paginator) fetchLimitId(db *sql.DB, uniqIdName, tableName string, whereCond map[string]interface{}, ruleCond map[string]string, orderby string, exCondVal []interface{}) int {
	log.Println(uniqIdName, whereCond)

	uniqId := -1
	selFields := make(map[string]interface{})
	selFields[uniqIdName] = &uniqId

	msqls, selargs, whereargs := sqlutils.Sqls_CompSelect(tableName, selFields, whereCond, ruleCond)

	whereargs = append(whereargs, exCondVal...)

	//" order by " + uniqIdName + " desc "
	msqls = msqls + orderby + " limit " + fmt.Sprint(this.GetSQLOffset()) + ", 1 "

	log.Println("FetchLimitId sqls : ", msqls)

	err := sqlutils.Sqls_Do_QueryRowAndScan(db, msqls, selargs, whereargs)
	if err != nil {
		return -1
	}

	return uniqId
}

func (this *Paginator) MergeSQLs(sqls string) string {
	return sqls + " limit " + fmt.Sprint(this.GetSQLLimit())
}

func NewPaginator(total int64, pageSize int, reqPage int64) *Paginator {
	log.Println("total is : ", total, "; pageSize is ", pageSize, "; reqPage is ", reqPage)

	return &Paginator{
		TotalCounts: total,
		PageSize:    pageSize,
		CurrentPage: reqPage,
		TotalPage:   int64(math.Ceil(float64(total) / float64(pageSize))),
	}
}
