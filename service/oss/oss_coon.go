package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	cfg "ntfileupload/config"
)

var ossCli *oss.Client

// Client 创建oss client 对象
func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}

	ossCli, err := oss.New(cfg.OSSEndpoint, cfg.OSSAccessKeyID, cfg.OSSAccessKeySecret)
	if nil != err {
		fmt.Println("oss连接错误: ", err.Error())
		return nil
	}
	return ossCli
}

// Bucket 获取bucket储存空间
func Bucket() *oss.Bucket {
	cli := Client()
	if nil != cli {
		bucket, err := cli.Bucket(cfg.OSSBucket)
		if nil != err {
			fmt.Println("获取Bucket错误: ", err.Error())
			return nil
		}
		return bucket
	}
	return nil
}
