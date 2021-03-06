package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// FilePermissionAll is the File Permission for the read/write file
	FilePermissionAll = 0666
	// GeneralFolderPermission is the Folder Permission for the general folder
	GeneralFolderPermission = 0766
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
	basePath string // config directory
	fullPath string // basePath + yyyymm
	zipPath  string // config directory
}

// 파일 체크 여부 확인
func checkLogFile(path string) (*os.File, error) {
	var file *os.File
	fullLogFile := fmt.Sprint(path, "/", fmt.Sprintf("%02d.log", time.Now().Day()))

	if _, err := os.Stat(fullLogFile); err != nil {
		// fmt.Println("create log")
		file, err = os.OpenFile(fullLogFile, os.O_CREATE|os.O_RDWR, FilePermissionAll)
		if err != nil {
			return file, err
		}
	} else {
		// fmt.Println("append log")
		file, err = os.OpenFile(fullLogFile, os.O_APPEND|os.O_RDWR, FilePermissionAll)
		if err != nil {
			return file, err
		}
	}

	return file, nil
}

// 초기 기동 시 로그 디렉토리 생성
func (w *logWriter) MakeLogDirectory() {
	os.MkdirAll(w.basePath, GeneralFolderPermission) // 로그 기본 디렉토리
	now := time.Now()
	w.fullPath = fmt.Sprint(w.basePath, "/", fmt.Sprintf("%04d%02d", now.Year(), int(now.Month())))
	os.Mkdir(w.fullPath, GeneralFolderPermission) // basePath + 년월일 디렉토리 생성

	os.MkdirAll(w.zipPath, GeneralFolderPermission) // 로그 압축 디렉토리 생성
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

func isEmptyDirectory(filePath string) (bool, error) {
	infos, err := ioutil.ReadDir(filePath)
	if err != nil {
		return false, err
	}

	for _, info := range infos {
		if info.Name() != "." || info.Name() != ".." {
			return false, nil
		}
	}

	return true, nil
}

// keepDays 이전 로그 삭제
func (w *logWriter) logTrucate(keepDays int) {
	fmt.Println("logTrucate start")
	for {
		time.Sleep(12 * time.Hour)
		// now := time.Now()
		// duration := now.Sub(lastCheckedLogTruncate)
		// fmt.Println("lastCheckedLogTruncate : ", duration.Hours())
		// // delete log files
		// if duration.Hours() > 24 {
		// 	targetDay := now.AddDate(0, 0, -1*keepDays)
		// 	for day := 1; day < keepDays; day++ {
		// 		fullPath := fmt.Sprint(w.basePath, ",", fmt.Sprintf("%04d%02d", targetDay.Year(), int(targetDay.Month())), "/", fmt.Sprintf("%02d.log", targetDay.Day()))
		// 		_, err := os.Stat(fullPath)
		// 		fmt.Print(fullPath)
		// 		if err == nil {
		// 			fmt.Println(" is removed")
		// 			os.Remove(fullPath)
		// 		}

		// 		targetDay = targetDay.AddDate(0, 0, -1)
		// 	}
		// }

		// remove empry directory
		fmt.Println("check empty directory")
		dirs, err := ioutil.ReadDir(w.basePath)
		if err != nil {
			fmt.Printf("Can't check the log folder '%s' due to cleansing the empty previous log folders. %s", w.basePath, err)
		} else {
			for _, dir := range dirs {
				removePath := filepath.Join(w.basePath, dir.Name())
				empty, err := isEmptyDirectory(removePath)
				if empty && err == nil {
					os.Remove(removePath)
				}
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

func setLogLevel(level int) int {
	setLevel := 0
	switch level {
	case 1, 2:
		setLevel = 0
	case 3:
		setLevel = EnableWarnLog
	case 4:
		setLevel = EnableWarnLog | EnableInfoLog
	case 5:
		setLevel = EnableWarnLog | EnableInfoLog | EnableTraceLog
	}

	return setLevel
}

func InitLog(logPath string, logLevel int, keepDays int, zipPath string, zipTime int) {
	w := &logWriter{
		basePath: logPath,
		zipPath:  zipPath,
	}

	// 로그 디렉토리 생성
	w.MakeLogDirectory()

	// 로그 삭제
	if keepDays > 0 {
		go w.logTrucate(keepDays)
	}

	// 로그 파일 압축
	go LogCompress(logPath, zipPath, zipTime)

	lastCheckedLogTruncate = time.Now().Add(-24 * time.Hour)
	// fmt.Println("lastCheckedLogTruncate : ", lastCheckedLogTruncate)

	// 레벨 별 로그 설정
	var enableLogs int = 0
	if logLevel != 0 {
		enableLogs = setLogLevel(logLevel)
	}
	fmt.Printf("enableLogs : %d\n", enableLogs)

	Trace = log.New(Option(enableLogs&EnableTraceLog == EnableTraceLog, w, ioutil.Discard).(io.Writer), "TRACE\t :", log.Ltime|log.Lshortfile)
	Info = log.New(Option(enableLogs&EnableInfoLog == EnableInfoLog, w, ioutil.Discard).(io.Writer), "INFO\t :", log.Ltime|log.Lshortfile)
	Warn = log.New(Option(enableLogs&EnableWarnLog == EnableWarnLog, w, ioutil.Discard).(io.Writer), "WARN\t :", log.Ltime|log.Lshortfile)
	Err = log.New(w, "ERROR\t :", log.Ltime|log.Lshortfile)
	Critical = log.New(w, "CRITICAL :", log.Ltime|log.Lshortfile)
}
