package main

import (
	"cdtoos/comm"
	"cdtoos/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func init() {
	//日志格式化
	utils.DefineLog()
	return
}

func main() {
	//读取配置文件
	comm.GetConfig()

	//创建OSSClient实例。
	client, err := oss.New(viper.Get("oss.url").(string), viper.Get("oss.id").(string), viper.Get("oss.key").(string))
	if err != nil {
		utils.HandleError(err)
	}

	//获取Bucket存储空间
	ossBucket, err := client.Bucket(viper.Get("oss.bucket").(string))
	if err != nil {
		utils.HandleError(err)
	}
	//读取配置
	filePath := viper.Get("local.Path").(string)
	utils.Info.Printf("读取配置文件中的本地文件路径:%s", filePath)
	webFilePath := viper.Get("local.WebPath").(string)
	utils.Info.Printf("读取配置文件中的WEB静态文件本地文件路径:%s", webFilePath)
	keyWordsOfWeb := viper.Get("local.KeywordsOfWeb").(string)
	utils.Info.Printf("区分前后端服务:配置文件local.Server中包含%s关键字的为前端服务", keyWordsOfWeb)
	serverGroup := viper.Get("local.Server").(string)
	DingDing := viper.Get("webHook").(string)
	utils.Info.Printf("webHook地址:%s", DingDing)
	//以字符串","切割服务组,转换为切片
	sep := ","
	arr := strings.Split(serverGroup, sep)
	for {
		for _, serverName := range arr {
			web := strings.Contains(serverName, keyWordsOfWeb)
			if web == true {
				// 更新前端静态文件
				utils.Info.Printf("读取配置文件中的服务名:%s", serverName)
				webServerFile := webFilePath  + serverName + "\\" + "dist"
				utils.Info.Printf("拼接服务本地文件路径:%s", webServerFile)
				comm.CmdCommWeb(ossBucket,webFilePath,serverName,webServerFile)

			} else {
				// 更新后端服务
				utils.Info.Printf("读取配置文件中的服务名:%s", serverName)
				serverFile := filePath  + serverName + "\\" + serverName + ".jar"
				utils.Info.Printf("拼接服务本地文件路径:%s", serverFile)
				comm.CmdComm(ossBucket, filePath, serverName, serverFile, DingDing)
			}
		}
		time.Sleep(120 * time.Second)
	}
}