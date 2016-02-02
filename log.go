package log

import (
	"fmt"
	sysLog "log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"runtime"
	"strings"
	"go/build"
)

const (
	ERROR   = "error"
	INFO    = "info"
	WARRING = "warring"
	PANIC = "panic"
)

// 日志数据处理
type data struct {
	filePath *os.File
	debug    bool
}

// 日志对象
var log *sysLog.Logger

// 日志数据对象
var logData data

// 程序的运行目录
var runDir string

func init() {
	runDir = RunDir();

	logFilePath, err := logFile()
	if err != nil {
		panic(err)
	}

	logData = data{
		filePath: logFilePath,
		debug:    false,
	}

	log = sysLog.New(&logData, "", 0)

	// 日志文件计算
	go func() {
		timeLocation, err := time.LoadLocation("Asia/Chongqing");
		if err != nil {
			log.Println(err)
		}
		now := time.Now()
		tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, timeLocation);

		time.Sleep(time.Duration(tomorrow.Unix()-now.Unix()) * time.Second)
		for {
			timeName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
			if timeName != logData.FileName() {
				if logFilePath, err := logFile(); err == nil {
					logData.CloseFile()
					logData.filePath = logFilePath
				}
			}
			// 开始处理休眠到一次应该的启动的时候
			time.Sleep(86400 * time.Second)
		}
	}()
}

// 写入数据到文件中
func (this *data) Write(p []byte) (int, error) {
	if this.debug {
		pStr := string(p)
		pStr = strings.Replace(pStr, fmt.Sprintf("[%s]", ERROR), fmt.Sprintf("[\033[31m%s\033[0m]", ERROR), 1)
		pStr = strings.Replace(pStr, fmt.Sprintf("[%s]", WARRING), fmt.Sprintf("[\033[33m%s\033[0m]", WARRING), 1)
		pStr = strings.Replace(pStr, fmt.Sprintf("[%s]", INFO), fmt.Sprintf("[\033[32m%s\033[0m]", INFO), 1)
		fmt.Print(pStr)
	}
	return this.filePath.Write(p)
}

// 获取当前正在写入的日志文件名
func (this *data) FileName() string {
	return filepath.Base(this.filePath.Name())
}

// 关闭文件指针
func (this *data) CloseFile() {
	this.filePath.Close()
}

// 设置是否为 debug 模式
func SetDebug(debugs bool) {
	logData.debug = debugs
}

// 写入错误
func Error(msg interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Printf("[%s] %s %s:%d %s", ERROR, time.Now().Format("2006-01-02 15:04:05"), strings.Replace(filePath, build.Default.GOPATH, "", 1), line, msg)
}

// 写入信息
func Info(msg interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Printf("[%s] %s %s:%d %s", INFO, time.Now().Format("2006-01-02 15:04:05"), strings.Replace(filePath, build.Default.GOPATH, "", 1), line, msg)
}

// 写入警告
func Warring(msg interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Printf("[%s] %s %s:%d %s", WARRING, time.Now().Format("2006-01-02 15:04:05"), strings.Replace(filePath, build.Default.GOPATH, "", 1), line, msg)
}

// 显示 panic
func Panic(msg interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Panicf("[%s] %s %s:%d %s", PANIC, time.Now().Format("2006-01-02 15:04:05"), strings.Replace(filePath, build.Default.GOPATH, "", 1), line, msg)
}

// 获取日志文件
func logFile() (*os.File, error) {
	rootDir := fmt.Sprintf("%s/log/", runDir)
	logFilePath := fmt.Sprintf("%s%s.log", rootDir, time.Now().Format("2006-01-02"))

	fileinfo, err := os.Stat(rootDir)
	// 当文件夹不存在时
	if err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(rootDir, 0700); err != nil {
			panic(err)
		}
	} else if err != nil { // 文件夹存在,不过出现了错误
		panic(err)
	} else if err == nil { // 文件夹存在,无错误
		// 当日志目录不是目录时
		if !fileinfo.IsDir() {
			panic(fmt.Sprintf("path %s not a dir", rootDir))
		}
	}

	return os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// 获取程序的运行目录
func RunDir() string {
	rootDir, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}

	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		panic(err)
	}

	return filepath.Dir(rootDir)
}
