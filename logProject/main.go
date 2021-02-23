package main

import (
	"fmt"
	"log"
	"os"
	"sampleProject/logProject/util"
	"strconv"
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

// 말모이 *.yaml 에 들어갈 설정 데이터
type logConfig struct {
	logPath    string
	keepDays   int
	logLevel   int
	logZipPath string
	logZipTime int // 0 ~ 23 hour 기준
}

func main() {
	//util.EnableInfoLog|util.EnableTraceLog|util.EnableWarnLog
	config := logConfig{
		logPath:    "../logs/" + os.Args[0][2:],
		keepDays:   1,
		logLevel:   5,
		logZipPath: "../zips/" + os.Args[0][2:],
		logZipTime: 11,
	}

	fmt.Printf("%+v\n", config)

	switch os.Args[1] {
	case "-z":
		logZipTime, _ := strconv.Atoi(os.Args[2])
		util.LogCompress(config.logPath, config.logZipPath, logZipTime)
		return
	case "-uz":
		fmt.Println("to-do decompress")
		return
	}

	util.InitLog(config.logPath, config.logLevel, config.keepDays, config.logZipPath, config.logZipTime)
	testLogWrite()
}
