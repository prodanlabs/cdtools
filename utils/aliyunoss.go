package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"time"
)

func HandleError(err error) {
	Error.Println("Error:", err)
	os.Exit(-1)
}

// 定义进度条监听器。
type OssProgressListener struct {
}

// 定义进度变更事件处理函数。
func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		Info.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		Info.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		Info.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		Info.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

//通过创建时间判断最新的文件，并返回最新的文件名与创建时间
func GetOssFileInfo(ossBucket *oss.Bucket, ossFile string) (ossFileName string, fileTime time.Time) {
	bucketName := ossBucket
	// 列举文件。
	prefix := ossFile
	for {
		lsRes, err := bucketName.ListObjects(oss.Prefix(prefix), oss.MaxKeys(200))
		if err != nil {
			HandleError(err)
		}
		/*		for index, object := range lsRes.Objects {
				fmt.Println("Bucket: ", index, object.Key)
				//fmt.Println("tNow(string format):", object.LastModified.Format("2006-01-02 15:04:05"))
				v := lsRes.Objects[index].LastModified.Format("2006-01-02 15:04:05")
				fmt.Println(v)
			}*/
		//遍历oss.ObjectProperties切片，找出最新的文件
		newFileTime := DefineFormatTime(lsRes.Objects[0].LastModified)
		newFileTimeIndex := 0
		for i := 0; i < len(lsRes.Objects); i++ {
			//从第二个元素开始循环比较，如果发现有新的数，则交换
			comPareFile := DefineFormatTime(lsRes.Objects[i].LastModified)
			if newFileTime.Before(comPareFile) {
				newFileTime = DefineFormatTime(lsRes.Objects[i].LastModified)
				newFileTimeIndex = i
			} else {
				Info.Printf("oss文件: %v 是旧版本\n",lsRes.Objects[i].Key)
			}
		}
		//阿里云oss创建时间和time.Parse方法用的是UTC时区，newFileTime需要转换时区，+8小时
		return lsRes.Objects[newFileTimeIndex].Key, time.Unix(newFileTime.Unix(), 0)
	}
}

func GetOssFile(ossBucket *oss.Bucket, fileName, localPath string) {
	isExist, err := ossBucket.IsObjectExist(localPath)
	if err != nil {
		Error.Println("Error:", err)
		os.Exit(-1)
	}
	Info.Println("文件是否存在:", isExist)
	// 分片下载。3个协程并发下载分片，开启断点续传下载。
	// 带进度条的下载。
	// "LocalFile"为filePath，100*1024为partSize。
	err = ossBucket.DownloadFile(fileName, localPath, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""), oss.Progress(&OssProgressListener{}))
	if err != nil {
		Error.Println("Error:", err)
		os.Exit(-1)
	}
}
