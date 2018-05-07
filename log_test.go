package log

import (
	"fmt"
	"testing"
)

func TestLogFile(t *testing.T) {
	file, _ := logFile()
	fmt.Println(file.Name())
	SetConfig(Config{
		ShowFilePath:  true,
		ShortFilePath: true,
		GOPATHDeep:    2,
		Debug:         true,
	})
	Error("这是一个错误")
	Warring("这是一个警告")
	Info("这是一个信息")
}
