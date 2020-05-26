package utils

import (
	"archive/zip"
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}
func GetLocalFileLastWriteTime(path string) (createTime time.Time, err error) {
	osType := runtime.GOOS
	fileInfo, err := os.Stat(path)
	if err != nil {
		info, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
		return info, err
	}
	if osType == "windows" {
		// Sys()返回的是interface{}，所以需要类型断言，不同平台需要的类型不一样
		wFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
		//获取文件最后修改的时间
		//转换time.Time类型与OSS文件比较
		//tNanSeconds := wFileSys.CreationTime.Nanoseconds() / 1e9 /// 返回的是纳秒,/1e9秒
		tNanSeconds := wFileSys.LastWriteTime.Nanoseconds() / 1e9
		return time.Unix(tNanSeconds, 0), err
	}
	info, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	return info, err
}

func WinCmdDel(localFileName string) {
	cmd := exec.Command("cmd", "/c", "del", localFileName)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		Error.Fatal(err)
	}
	Info.Printf("%s\n", ConvertByte2String(stdoutStderr, GB18030))
}

func WinCmdRename(localFileName string) {
	t := time.Now().Format("20060102-1504")
	bakLocalFileName := localFileName + t
	err := os.Rename(localFileName, bakLocalFileName)
	if err != nil {
		Error.Fatal(err)
	} else {
		Info.Println(`备份文件成功`)
	}
}


func WinCmdServer(flags, servername string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("net", flags, servername)
	//cmd.Dir = workDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	Info.Println(stdout.String(), stderr.String(), err)
}

func IsProcessExist(appName string) (bool, string, int) {
	appAry := make(map[string]int)
	cmd := exec.Command("cmd", "/C", "tasklist")
	output, _ := cmd.Output()
	n := strings.Index(string(output), "System")
	if n == -1 {
		Info.Println("no find")
		os.Exit(1)
	}
	data := string(output)[n:]
	fields := strings.Fields(data)
	for k, v := range fields {
		if v == appName {
			appAry[appName], _ = strconv.Atoi(fields[k+1])
			return true, appName, appAry[appName]
		}
	}
	return false, appName, -1
}

func Health(webHook, localServerName string) bool {
	app := localServerName + ".exe"
	b, s, i := IsProcessExist(app)
	Info.Println(b, s, i)
	if b == true {
		successMsg := app + "更新成功"
		SendDingMsg(webHook, successMsg)
	} else {
		Error.Printf("更新失败！%s服务状态为%d，进程找不到\n", s, i)
		errMsg := app + "进程不存在！更新失败!"
		SendDingMsg(webHook, errMsg)
	}
	return b
}

//解压
func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}