package log

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestLogFileDebugFlase(t *testing.T) {
	file, _ := logFile()
	fmt.Println(file.Name())
	SetConfig(Config{
		ShowFilePath:  true,
		ShortFilePath: true,
		GOPATHDeep:    2,
		Debug:         false,
	})
	Error("这是一个错误")
	Warring("这是一个警告")
	Info("这是一个信息")

	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		content = []byte(err.Error())
	}
	fmt.Println("File Content: \n", string(content))
}

func TestLogFileDebugTrue(t *testing.T) {
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
