package main

import (
    "database/sql"
    "fmt"
    _ "github.com/Go-SQL-Driver/MySQL"
    //"math/rand"
    //"encoding/json"
    "time"
)

type GameManager struct {
    db_handler *sql.DB
}

func (gm *GameManager) Connect() bool {
    var e error
    //说明：第一个参数"mysql"表示打开一个MYSQL链接
    //第二个参数是一个DSN，格式是：
    //[username[:password]@][protocol[(address)]]/dbname[?param1=value1&paramN=valueN]
    gm.db_handler, e = sql.Open("mysql", "root:123456@/test?charset=utf8")
    if e != nil {
        fmt.Printf("conn error.")
        return false
    }
    return true
}

func (gm *GameManager) UpdatePermenentStat(id, uid, value int) {
    _, e := gm.db_handler.Exec("insert into t_gm_permenent values(%d,%d,%d,%d)", id, time.Now().Unix(), uid, value)
    if e != nil {
        fmt.Printf("insert error:", id, uid, value)
        return
    }
}

func (gm *GameManager) UpdateDayStat(id, uid, value int) {
    //func UpdateDayStatistic(db_handler *sql.DB, id, uid, value int) {
    table := fmt.Sprintf("t_gm_day_%4d%02d%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
    str := fmt.Sprintf("insert into %s () values(%d,%d,%d,%d);", table, id, time.Now().Unix(), uid, value)
    _, e := gm.db_handler.Exec(str)
    if e != nil {
        fmt.Printf("insert error:", id, uid, value)
        return
    }
}

func main() {
    fmt.Printf("query:\n")
    var gm *GameManager = new(GameManager)
    //defer gm.db_handler.Close()
    gm.Connect()
    //sqlStr := fmt.Sprintf("select id,logtime from t_gm_day_20140401 where id=101")

    //for _ = range time.Tick(1e10) {
    //e := gm.db_handler.Ping()
    //if e != nil {
    //fmt.Printf("connect error\n")
    //} else {
    //fmt.Printf("connect succ\n")
    //}
    //}
    gm.db_handler.Close()
}
