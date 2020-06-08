package comm

import (
	"cdtoos/utils"
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"os"
	"time"
)

func GetConfig() {
	var (
		config string
		path   string
		suffix string
	)
	flag.StringVar(&config, "c", "config", "配置文件名")
	flag.StringVar(&path, "d", ".", "配置文件路径，默认当前目录")
	flag.StringVar(&suffix, "t", "json", "配置文件后缀")
	v := flag.Bool("v", false, "version")
	flag.Parse()
	if *v {
		_, _ = fmt.Fprint(os.Stderr, StringifyMultiLine())
		os.Exit(1)
	}

	utils.Info.Printf("cdtoos的配置文件:%s", path+config+"."+suffix)
	viper.SetConfigName(config) //设置配置文件
	viper.AddConfigPath(path)   //添加配置文件所在的路径
	viper.SetConfigType(suffix) //设置配置文件类型
	confErr := viper.ReadInConfig()
	if confErr != nil {
		utils.Error.Printf("config file error: %s\n", confErr)
		os.Exit(1)
	}
}

func CmdComm(ossBucket *oss.Bucket, localPath, localServerName, localServerFile, webHook string) {

	utils.Info.Printf("Start: 服务%s开始检查版本更新....\n", localServerName)

	//获取oss最新的文件名、文件创建时间
	ossFileName, ossFileCreateTime := utils.GetOssFileInfo(ossBucket, localServerName)
	utils.Info.Printf("oss最新的文件:%s 创建时间:%v\n", ossFileName, ossFileCreateTime)

	//获取本地文件的创建时间
	localFileCreateTimeInfo, err := utils.GetLocalFileLastWriteTime(localServerFile)
	if err != nil {
		defer func() {
			if msg := recover(); msg != nil {
				utils.Warning.Println(msg)
			}
		}()
		utils.Error.Println(err)
	}

	//格式化时间，oss与本地文件比较创建时间
	if utils.DefineFormatTime(ossFileCreateTime).Before(utils.DefineFormatTime(localFileCreateTimeInfo)) {
		//2020-05-10 22:02:27 +0000 UTC  > 2020-05-10 05:39:55 +0000 UTC
		utils.Info.Printf("本地已是最新版本:local %v => oss %v", utils.DefineFormatTime(localFileCreateTimeInfo), utils.DefineFormatTime(ossFileCreateTime))
	} else {
		utils.Info.Printf("发现新版本:local %v < oss %v", utils.DefineFormatTime(localFileCreateTimeInfo), utils.DefineFormatTime(ossFileCreateTime))
		//下载最新文件到本地
		newLocalFileName := localPath + "\\" + ossFileName
		utils.GetOssFile(ossBucket, ossFileName, newLocalFileName)

		////停止服务
		utils.Info.Printf("正在停止%s服务...", localServerName)
		utils.WinCmdServer("stop", localServerName)

		//备份旧版本文件
		utils.WinCmdRename(localServerFile)

		//解压
		unzipLocalFilePath := localPath + localServerName
		utils.Info.Printf("解压路径: %s",unzipLocalFilePath)
		err := utils.UnZip(newLocalFileName, unzipLocalFilePath)
		if err != nil {
			utils.Error.Printf("解压失败: %v", err)
		} else {
			utils.Info.Printf("解压完成.")
		}

		//删除oss下载的压缩文件
		delErr := os.Remove(newLocalFileName)
		if delErr != nil {
			utils.Error.Printf("删除失败: %v", delErr)
		} else {
			utils.Info.Printf("已删除压缩文件.")
		}

		//启动服务
		utils.Info.Printf("正在启动%s服务...", localServerName)
		utils.WinCmdServer("start", localServerName)
		time.Sleep(30 * time.Second)

		//进程检查
		state := utils.Health(webHook, localServerName)
		if state == false {
			utils.Error.Println("第一次启动失败！")
			utils.Warning.Printf("正在尝试启动%s服务...", localServerName)
			utils.WinCmdServer("start", localServerName)
			utils.Health(webHook, localServerName)
			time.Sleep(30 * time.Second)
		} else {
			utils.Info.Printf("健康检查: 服务%s进程存在.\n", localServerName)
		}
	}
	utils.Info.Printf("END: 服务%s检查更新操作退出.\n", localServerName)
}

func CmdCommWeb(ossBucket *oss.Bucket, localPath, localServerName, localServerFile string) {
	utils.Info.Printf("Start: 服务%s开始检查版本更新....\n", localServerName)

	//获取oss最新的文件名、文件创建时间
	ossFileName, ossFileCreateTime := utils.GetOssFileInfo(ossBucket, localServerName)
	utils.Info.Printf("oss最新的文件:%s 创建时间:%v\n", ossFileName, ossFileCreateTime)

	//获取本地文件的创建时间
	localFileCreateTimeInfo, err := utils.GetLocalFileLastWriteTime(localServerFile)
	if err != nil {
		defer func() {
			if msg := recover(); msg != nil {
				utils.Warning.Println(msg)
			}
		}()
		utils.Error.Println(err)
	}

	//格式化时间，oss与本地文件比较创建时间
	if utils.DefineFormatTime(ossFileCreateTime).Before(utils.DefineFormatTime(localFileCreateTimeInfo)) {
		//2020-05-10 22:02:27 +0000 UTC  > 2020-05-10 05:39:55 +0000 UTC
		utils.Info.Printf("本地已是最新版本:local %v => oss %v", utils.DefineFormatTime(localFileCreateTimeInfo), utils.DefineFormatTime(ossFileCreateTime))
	} else {
		utils.Info.Printf("发现新版本:local %v < oss %v", utils.DefineFormatTime(localFileCreateTimeInfo), utils.DefineFormatTime(ossFileCreateTime))
		//下载最新文件到本地
		newLocalFileName := localPath + "\\" + ossFileName
		utils.GetOssFile(ossBucket, ossFileName, newLocalFileName)

		//备份旧版本文件
		utils.WinCmdRename(localServerFile)

		//解压
		unzipLocalFilePath := localPath + localServerName
		utils.Info.Printf("解压路径: %s",unzipLocalFilePath)
		err := utils.UnZip(newLocalFileName, unzipLocalFilePath)
		if err != nil {
			utils.Error.Printf("解压失败: %v", err)
		} else {
			utils.Info.Printf("解压完成.")
		}

		//删除oss下载的压缩文件
		delErr := os.Remove(newLocalFileName)
		if delErr != nil {
			utils.Error.Printf("删除失败: %v", delErr)
		} else {
			utils.Info.Printf("已删除压缩文件.")
		}
	}
}