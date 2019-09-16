package main

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Foo test data
type Foo struct {
	Key int
	Val string
}

const (
	VAL = `Lorem ipsum dolor`
)

var (
	DBW *sql.DB

	DSN      = "root:12345678@tcp(127.0.0.1:3306)/gotest?charset=utf8"
	SQLQuery = "INSERT INTO `foo` (`item`,`itemval`) VALUES(?,?)"
	StmtMain *sql.Stmt
	wg       sync.WaitGroup
)

func Store(d Foo) {
	defer wg.Done()
	_, err := StmtMain.Exec(d.Key, d.Val)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var (
		errDbw  error
		errStmt error
	)
	concurrencyLevel := runtime.NumCPU() * 8
	DBW, errDbw = sql.Open("mysql", DSN)
	if errDbw != nil {
		log.Fatalln(errDbw)
	}
	DBW.SetMaxIdleConns(concurrencyLevel)
	defer DBW.Close()
	StmtMain, errStmt = DBW.Prepare(SQLQuery)
	if errStmt != nil {
		log.Fatalln(errStmt)
	}
	defer StmtMain.Close()
	//populate data
	dd := Foo{
		Key: 0,
		Val: VAL,
	}
	t0 := time.Now()
	for i := 0; i < 300000; {
		for k := 0; k < concurrencyLevel; k++ {
			i++
			if i > 300000 {
				break
			}
			dd.Key = i
			wg.Add(1)
			go Store(dd)
		}
		wg.Wait()
		if i > 300000 {
			break
		}
	}
	t1 := time.Now()
	fmt.Printf("%v per second.\n", 300000.0/t1.Sub(t0).Seconds())
}
