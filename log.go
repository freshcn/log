package log

import (
	"fmt"
	"go/build"
	"io"
	sysLog "log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// ERROR 错误类型
	ERROR = "error"
	// INFO 普通信息
	INFO = "info"
	// WARRING 警告信息
	WARRING = "warring"
	// PANIC 恐慌信息
	PANIC = "panic"
)

// 日志数据处理
type data struct {
	filePath *os.File
	debug    bool
}

var (
	// 日志数据对象
	logData data
	// 日志对象
	log *sysLog.Logger
	// pathSeparator 系统的目录间隔符
	pathSeparator = string(os.PathSeparator)
)

func init() {
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
		now := time.Now().In(timezone)
		tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, timezone)

		time.Sleep(time.Duration(tomorrow.Unix()-now.Unix()) * time.Second)
		for {
			if !logData.debug {
				timeName := fmt.Sprintf("%s.log", time.Now().In(timezone).Format("2006-01-02"))
				if timeName != logData.FileName() {
					if logFilePath, err := logFile(); err == nil {
						logData.CloseFile()
						logData.filePath = logFilePath
					}
				}
			}
			// 开始处理休眠到一次应该的启动的时候
			time.Sleep(86400 * time.Second)
		}
	}()
}

// Write 写入数据到文件中
func (d *data) Write(p []byte) (int, error) {
	if d.debug {
		pStr := string(p)
		pStr = strings.Replace(pStr, fmt.Sprintf("[%s]", ERROR), fmt.Sprintf("[\033[31m%s\033[0m]", ERROR), 1)
		pStr = strings.Replace(pStr, fmt.Sprintf("[%s]", PANIC), fmt.Sprintf("[\033[31m%s\033[0m]", PANIC), 1)
		pStr = strings.Replace(pStr, fmt.Sprintf("[%s]", WARRING), fmt.Sprintf("[\033[33m%s\033[0m]", WARRING), 1)
		pStr = strings.Replace(pStr, fmt.Sprintf("[%s]", INFO), fmt.Sprintf("[\033[32m%s\033[0m]", INFO), 1)
		fmt.Print(pStr)
		return 0, nil
	}
	return d.filePath.Write(p)
}

// FileName 获取当前正在写入的日志文件名
func (d *data) FileName() string {
	return filepath.Base(d.filePath.Name())
}

// CloseFile 关闭文件指针
func (d *data) CloseFile() {
	d.filePath.Close()
}

// SetDebug 设置是否为 debug 模式
func SetDebug(debugs bool) {
	logData.debug = debugs
}

// Error 写入错误
func Error(msg ...interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Printf(formatLayout(), ERROR, formatTime(), formatFilePath(filePath, line), fmt.Sprint(msg...))
}

// Errorf 格式化错误信息输出
func Errorf(format string, msg ...interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Printf(formatLayout(), ERROR, formatTime(), formatFilePath(filePath, line), fmt.Sprintf(format, msg...))
}

// Info 写入信息
func Info(msg ...interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Printf(formatLayout(), INFO, formatTime(), formatFilePath(filePath, line), fmt.Sprint(msg...))
}

// Infof 格式化信息输出
func Infof(format string, msg ...interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Printf(formatLayout(), INFO, formatTime(), formatFilePath(filePath, line), fmt.Sprintf(format, msg...))
}

// Warring 写入警告
func Warring(msg ...interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Printf(formatLayout(), WARRING, formatTime(), formatFilePath(filePath, line), fmt.Sprint(msg...))
}

// Warringf 格式化错误信息输出
func Warringf(format string, msg ...interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Printf(formatLayout(), WARRING, formatTime(), formatFilePath(filePath, line), fmt.Sprintf(format, msg...))
}

// Panic 显示 panic
func Panic(msg ...interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Panicf(formatLayout(), PANIC, formatTime(), formatFilePath(filePath, line), fmt.Sprint(msg...))
}

// Panicf 格式化panic信息输出
func Panicf(format string, msg ...interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	log.Panicf(formatLayout(), PANIC, formatTime(), formatFilePath(filePath, line), fmt.Sprintf(format, msg...))
}

// formatLayout 初始化显示格式
func formatLayout() (l string) {
	if config.ShowFilePath {
		l = "[%s] %s %s %s"
	} else {
		l = "[%s] %s%s%s"
	}
	return
}

// formatTime 格式化时间
func formatTime() (t string) {
	return time.Now().In(timezone).Format(config.TimeFormat)
}

// formatFilePath 处理文件地址
func formatFilePath(filepath string, line int) string {
	// 是否显示
	if !config.ShowFilePath {
		return " "
	}

	// 是否显示为长路径
	if !config.ShortFilePath {
		return fmt.Sprintf("%s:%d", filepath, line)
	}

	filepath = strings.TrimPrefix(filepath, fmt.Sprintf("%s%ssrc%s", build.Default.GOPATH, pathSeparator, pathSeparator))
	if config.GOPATHDeep > 0 {
		if tmp := strings.Split(filepath, pathSeparator); len(tmp) >= config.GOPATHDeep {
			filepath = strings.Join(tmp[config.GOPATHDeep:], pathSeparator)
		}
	}
	return fmt.Sprintf("%s:%d", filepath, line)
}

// logFile 获取日志文件
func logFile() (*os.File, error) {
	rootDir := fmt.Sprintf("%s%slog%s", config.FilePath, pathSeparator, pathSeparator)
	logFilePath := fmt.Sprintf("%s%s.log", rootDir, time.Now().Format("2006-01-02"))

	fileinfo, err := os.Stat(rootDir)
	// 当文件夹不存在时
	if err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(rootDir, 0766); err != nil {
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

// RunDir 获取程序的运行目录
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

// LogWriter 返回
func LogWriter() io.Writer {
	return &logData
}
