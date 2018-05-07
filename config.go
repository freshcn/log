package log

import (
	"time"
)

// Config 存储路径
type Config struct {
	// Timezone 时区
	Timezone string
	// TimeFormat 时间配置
	TimeFormat string
	// ShowFilePath 是否显示目录
	ShowFilePath bool
	// ShortFilePath 短文件名
	ShortFilePath bool
	// GOPATHDeep 除了gopath以外，项目目录的深度
	GOPATHDeep int
	// FilePath 日志文件存储目录
	FilePath string
	// Debug 是否为测试环境
	Debug bool
}

var (
	// DefaultConfig 默认配置
	DefaultConfig = Config{
		Timezone:      "Asia/Chongqing",
		TimeFormat:    "2006-01-02 15:04:05",
		ShowFilePath:  true,
		ShortFilePath: true,
		GOPATHDeep:    2,
		FilePath:      RunDir(),
		Debug:         false,
	}

	config = DefaultConfig

	timezone = setTimeZone()
)

// SetConfig 设置配置信息
func SetConfig(conf Config) {
	config = conf

	if config.Timezone == "" {
		config.Timezone = DefaultConfig.Timezone
	}
	if config.TimeFormat == "" {
		config.TimeFormat = DefaultConfig.TimeFormat
	}
	if config.FilePath == "" {
		config.FilePath = RunDir()
	}
	SetDebug(config.Debug)

	setTimeZone()
}

// setTime
func setTimeZone() *time.Location {
	timeLocation, err := time.LoadLocation(config.Timezone)
	if err != nil {
		log.Println(err)
		return nil
	}
	return timeLocation
}
