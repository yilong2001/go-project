package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"web/component/cfgutils"
	"web/dal/sqldrv"
)

/*
监听服务器中的key列表的变化情况，当有新的列表被添加进入来的时候，产生新的数据
*/

// 数据表名

// id值的队列
type idChan chan int64

// 正在被监听的key列表
var watchList = make(map[string]idChan)

// 监听协程
func watch() {
	log.Println("start idgen watch ... ")
	makeChan()

	for {
		//makeChan()
		for key, value := range watchList {
			if int64(len(value)) < getPreStep()/3 {
				var max, step int64 = 0, 0
				var err error
				for {
					max, step, err = updateId(key)
					if err != nil {
						log.Panic(err)
						//time.Sleep(5 * time.Millisecond)
					} else {
						break
					}
				}

				if max > 0 {
					for i := int64(max - step + 1); i <= max; i++ {
						value <- i
					}
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

}

// 为数据创建新的监听队列
func makeChan() {
	var list map[string]int64
	var err error
	var db *sql.DB = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	defer db.Close()

	for {
		list, err = idsList(db)
		if err != nil {
			log.Panic(err)
			//time.Sleep(5 * time.Millisecond)
		} else {
			break
		}
	}

	// 将不存在于当前监听列表中的数据添加到监听列表中
	for key, _ := range list {
		if _, exists := watchList[key]; !exists {
			watchList[key] = make(idChan, getPreStep()*2)
		}
	}

	// // 检查存在于当前监听列表中，却已经不再存在数据中的ID
	// for key, _ := range watchList {
	// 	if _, exists := list[key]; !exists {
	// 		delete(watchList, key)
	// 	}
	// }
}

// 获取所有的key的名字列表
func idsList(db *sql.DB) (arr map[string]int64, err error) {
	log.Println("start list all objects in idgen ... ")

	arr = make(map[string]int64, 0)

	sqlstr := "select `object_name`,`next_id` from " + IdGen_DB_Name
	log.Println(sqlstr)
	rows, err := db.Query(sqlstr)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var id int64
		if err = rows.Scan(&name, &id); err == nil {
			arr[name] = id
		} else {
			return
		}
	}
	return
}

// 向name所指定的id中申请新的ID空间
// 向数据库申请成功后返回新申请到的最大数和申请的数量
func updateId(name string) (num int64, preStep int64, err error) {
	var db *sql.DB = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	defer db.Close()

	num = 0

	preStep = getPreStep()
	_, err = db.Exec("update `"+IdGen_DB_Name+"` set next_id=(next_id+?) where object_name=?", preStep, name)
	if err != nil {
		log.Println(err)
		return
	}

	row := db.QueryRow("select `next_id` from "+IdGen_DB_Name+" where object_name=?", name)
	if row == nil {
		log.Println("in QueryRow, select next id is failed!")
		return
	}

	err = row.Scan(&num)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
