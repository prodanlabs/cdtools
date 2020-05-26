# cdtools

A CD tool for windows to actively synchronize OSS files

```PowerShell
> ./cdtoos.exe -v
Version=v0.0.2
GitCommitLog=6861f45995c3348abf7db0ee29a6434bb7dccc97 update version 0.0.2
BuildTime=2020-05-18 11:32:39
GoVersion=go version go1.14.1 windows/amd64
runtime=windows/amd64


> ./cdtoos.exe -h
Usage of G:\GoDevWorker\cdtools\cdtoos.exe:
  -c string
        配置文件名 (default "config")
  -d string
        配置文件路径，默认当前目录 (default ".")
  -t string
        配置文件后缀 (default "json")
  -v    version
```

#### 先决条件

* 阿里云oss
* 服务需注册到window服务
* 因为调用net命令，需管理员权限运行，否则提示没有权限

#### 配置文件

```josn
{
  "oss": {
    "url": "http://xxx.aliyuncs.com",
    "id": "xxx",
    "key": "xxx",
    "bucket": "xxx"
  },
  "local": {
    "Path": "E:\\APP\\",
    "Server": "demo-server,demo1-server"
  },
  "webHook": "https://oapi.dingtalk.com/robot/send?access_token=xxx"
}
```

1. 支持多服务，在`local.Server`增加, 以`,` 号隔开;
2. webHook是钉钉机器人的钩子，其他的钩子未测试;
3. `local.Path`服务文件路径。如 `E:\\APP\\demo-server\demo-server.jar`  oss同理 `${bucket}/demo-server/demo-server*.zip` ;
4. 目前只支持java的`.jar` ;

## 示例

Jenkins上传编译好的jar包到oss

```
mv ${projectName}*.jar ${projectName}.jar
zip -q -r ${projectName}-${BUILD_NUMBER}.zip ${projectName}.jar
sh /root/sh/oss-test.sh put ${projectName}-${BUILD_NUMBER}.zip ${projectName}/${projectName}-${BUILD_NUMBER}.zip
rm -f ${projectName}-${BUILD_NUMBER}.zip ${projectName}.jar
```

```

```



<img src="https://github.com/ProdanLabs/Golang-practice-project/blob/master/image/qrcode_for_gh.jpg" width="120">
