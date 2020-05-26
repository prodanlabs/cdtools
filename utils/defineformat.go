package utils

import (
	"io"
	"log"
	"os"
	"time"
)
var (
	Info         *log.Logger
	Warning      *log.Logger
	Error        *log.Logger
)
//time类型格式化
func DefineFormatTime(timeInFormat time.Time) time.Time {
	formatted,err := time.Parse("2006-01-02 15:04:05",timeInFormat.Format("2006-01-02 15:04:05"))
	if err != nil {
		Error.Println(err)
	}
	return formatted
}

//日志格式化
func DefineLog()  {
	//日志输出到message.log文件
	file := "./" + time.Now().Format("2006-01-02-") + "message" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	/*	log.SetOutput(logFile) // 将文件设置为log输出的文件
		log.SetPrefix("[cdTools]")
		log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)*/
	//日志终端标准输出并写入文件
	Info = log.New(io.MultiWriter(os.Stdout, logFile), "Info:", log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stdout, logFile), "Warning:", log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, logFile), "Error:", log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
}