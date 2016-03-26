# 日志记录程序

可对日志进行分级

## 日志存储位置

日志文件存放在程序的运行目录下的`log`目录中

## debug 模式

debug模式下可显示在控制器上,你所打开出来的日志

> log.SetDebug(true)

## 写入日志

写入错误类型

>func Error(msg ...interface{})

写入警告类型

>func Warring(msg ...interface{})

写入信息类型

>func Info(msg ...interface{})

写入panic类型

>func Panic(msg ...interface{})