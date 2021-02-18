package main

import (
	"fmt"
	"log"
	"os"
	"sampleProject/logProject/util"
	"time"
)

var (
	Trace    *log.Logger
	Info     *log.Logger
	Warn     *log.Logger
	Err      *log.Logger
	Critical *log.Logger
)

func testLogWrite() {
	fmt.Println("testLogWrite start")
	i := 0

	Trace = util.Trace
	Info = util.Info
	Warn = util.Warn
	Err = util.Err
	Critical = util.Critical

	for {
		Trace.Printf("testLogWrite : %d\n", i)
		Info.Printf("testLogWrite : %d\n", i)
		Warn.Printf("testLogWrite : %d\n", i)
		Err.Printf("testLogWrite : %d\n", i)
		Critical.Printf("testLogWrite : %d\n", i)
		time.Sleep(time.Millisecond * 5000)
		i++
	}
}

func main() {
	//util.EnableInfoLog|util.EnableTraceLog|util.EnableWarnLog
	basePath := "../logs/" + os.Args[0][2:]
	keepDays := 1
	fmt.Println(basePath)
	util.InitLog(basePath, util.EnableTraceLog|util.EnableWarnLog, keepDays)
	testLogWrite()
}
