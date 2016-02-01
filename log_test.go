package log
import (
	"testing"
	"fmt"
	"time"
)

func TestLogFile(t *testing.T) {
	file, _ := logFile()
	fmt.Println(file.Name())
	SetDebug(true)
	Error("这是一个错误")
	time.Sleep(100*time.Second)
}
