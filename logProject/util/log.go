package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

const (
	EnableTraceLog = 1
	EnableInfoLog  = 2
	EnableWarnLog  = 4
)

var (
	Trace    *log.Logger
	Info     *log.Logger
	Warn     *log.Logger
	Err      *log.Logger
	Critical *log.Logger
	syncLog  sync.Mutex
)

var lastCheckedLogTruncate time.Time

type logWriter struct {
	basePath string // config directory + binary
	fullPath string // basePath + yyyymm
}

// 파일 체크 여부 확인
func checkLogFile(path string) (*os.File, error) {
	var file *os.File
	fullLogFile := fmt.Sprint(path, "/", fmt.Sprintf("%02d.log", time.Now().Day()))

	if _, err := os.Stat(fullLogFile); err != nil {
		// fmt.Println("create log")
		file, err = os.OpenFile(fullLogFile, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return file, err
		}
	} else {
		// fmt.Println("append log")
		file, err = os.OpenFile(fullLogFile, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return file, err
		}
	}

	return file, nil
}

// 초기 기동 시 로그 디렉토리 생성
func (w *logWriter) MakeLogDirectory() {
	os.Mkdir(w.basePath, 0766) // 이미 존재하는 디렉토리는 skip
	now := time.Now()
	w.fullPath = fmt.Sprint(w.basePath, "/", fmt.Sprintf("%04d%02d", now.Year(), int(now.Month())))
	os.Mkdir(w.fullPath, 0766) // 이미 존재하는 디렉토리는 skip
}

func (w *logWriter) Write(p []byte) (int, error) {
	syncLog.Lock()
	defer syncLog.Unlock()

	file, err := checkLogFile(w.fullPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return file.Write(p)
}

// keepDays 이전 로그 삭제
func (w *logWriter) logTrucate(keepDays int) {
	for {
		time.Sleep(12 * time.Hour)
		now := time.Now()
		duration := now.Sub(lastCheckedLogTruncate)
		// delete log files
		if duration > 24 {
			targetDay := now.AddDate(0, 0, -1*keepDays)
			for day := 1; day <= keepDays; day++ {
				fullPath := fmt.Sprint(w.basePath, ",", fmt.Sprintf("%04d%02d", targetDay.Year(), int(targetDay.Month())), "/", fmt.Sprintf("%02d.log", targetDay.Day()))
				_, err := os.Stat(fullPath)
				if err == nil {
					os.Remove(fullPath)
				}

				targetDay = targetDay.AddDate(0, 0, 1)
			}
		}

		// remove empry directory
		dir, err := ioutil.ReadDir(w.fullPath)
		if err != nil {
			Warn.Printf("Can't check the log folder '%s' due to cleansing the empty previous log folders. %s", w.basePath, err)
		} else {
			removeYn := false
			for _, info := range dir {
				name := info.Name()
				if name != "." && name != ".." {
					removeYn = true
					break
				}
			}

			if removeYn {
				os.Remove(w.fullPath)
			}
		}
	}
}

func Option(condition bool, positive interface{}, negative interface{}) interface{} {
	if condition {
		return positive
	}

	return negative
}

func InitLog(logPath string, enableLogs int, keepDays int) {
	w := &logWriter{
		basePath: logPath,
	}

	// 로그 디렉토리 생성
	w.MakeLogDirectory()

	// 로그 삭제
	if keepDays > 0 {
		go w.logTrucate(keepDays)
	}

	lastCheckedLogTruncate = time.Now().Add(-24 * time.Hour)
	fmt.Println("lastCheckedLogTruncate : ", lastCheckedLogTruncate)

	// 레벨 별 로그 설정
	Trace = log.New(Option(enableLogs&EnableTraceLog == EnableTraceLog, w, ioutil.Discard).(io.Writer), "TRACE\t :", log.Ltime|log.Lshortfile)
	Info = log.New(Option(enableLogs&EnableInfoLog == EnableInfoLog, w, ioutil.Discard).(io.Writer), "INFO\t :", log.Ltime|log.Lshortfile)
	Warn = log.New(Option(enableLogs&EnableWarnLog == EnableWarnLog, w, ioutil.Discard).(io.Writer), "WARN\t :", log.Ltime|log.Lshortfile)
	Err = log.New(w, "ERROR\t :", log.Ltime|log.Lshortfile)
	Critical = log.New(w, "CRITICAL :", log.Ltime|log.Lshortfile)
}
