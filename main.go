package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const ConcurrencyDefault = 10

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
	ResourceDir     string
	PathPrefix      string
	BucketName      string
	Concurrency     int
}

func (c *Config) Init() {
	flag.StringVar(&c.AccessKeyID, "access-key-id", "", "OSS 账号的 AccessKeyId")
	flag.StringVar(&c.AccessKeySecret, "access-key-secret", "", "OSS 账号的 AccessKeySecret")
	flag.StringVar(&c.Endpoint, "endpoint", "", "OSS 的 Endpoint 地址")
	flag.StringVar(&c.ResourceDir, "resource-dir", "", "资源文件夹路径")
	flag.StringVar(&c.PathPrefix, "path-prefix", "", "上传到 OSS 的子目录路径")
	flag.StringVar(&c.BucketName, "bucket-name", "", "上传文件的 Bucket 名称")
	flag.IntVar(&c.Concurrency, "concurrency", ConcurrencyDefault, "上传文件的并发数量")
	flag.Parse()
	flag.Usage()
	if c.AccessKeyID == "" {
		handleError(errors.New("-access-key-id is required"), "OSS 账号的 AccessKeyId 不能为空")
	}
	if c.AccessKeySecret == "" {
		handleError(errors.New("-access-key-secret is required"), "OSS 账号的 AccessKeySecret 不能为空")
	}
	if c.Endpoint == "" {
		handleError(errors.New("-endpoint is required"), "OSS 的 Endpoint 地址不能为空")
	}
	if c.ResourceDir == "" {
		handleError(errors.New("-resource-dir is required"), "资源文件夹路径不能为空")
	}
	if c.PathPrefix == "" {
		handleError(errors.New("-path-prefix is required"), "上传到 OSS 的子目录路径不能为空")
	}
	if c.BucketName == "" {
		handleError(errors.New("-bucket-name is required"), "上传文件的 Bucket 名称不能为空")
	}
	if c.Concurrency <= 0 {
		handleError(errors.New("-concurrency must be more than 0"), "上传文件的并发数量不能小于0")
	}
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Println(message, err)
		os.Exit(1)
	}
}

func uploadFile(client *oss.Client, path string, resourceDir string, pathPrefix string, bucketName string) error {
	filename := filepath.Base(path)
	objectKey := filepath.Join(pathPrefix, filepath.Dir(strings.TrimPrefix(path, resourceDir)), filename)
	fmt.Println("File upload Ready", objectKey)
	bk, err := client.Bucket(bucketName)
	handleError(err, "Error creating OSS bucket:")

	err = bk.PutObjectFromFile(objectKey, path)
	if err != nil {
		handleError(err, "Error uploading file:"+err.Error())
	}
	fmt.Println("File uploaded successfully:", objectKey)
	return nil
}

func createOSSClient(c *Config) *oss.Client {
	client, err := oss.New(
		c.Endpoint,
		c.AccessKeyID,
		c.AccessKeySecret,
	)
	handleError(err, "Error creating OSS client:")
	return client
}

func uploadFiles(c *Config, client *oss.Client) {
	files := make(chan string, c.Concurrency)
	var wg sync.WaitGroup
	defer close(files) // 在发送完所有文件路径后关闭通道
	// 启动 goroutine 来上传文件
	for i := 0; i < c.Concurrency; i++ {
		wg.Add(1) // 为每个 goroutine 添加到等待组
		go func() {
			defer wg.Done() // 确保在 goroutine 退出时调用 Done
			for path := range files {
				if err := uploadFile(client, path, c.ResourceDir, c.PathPrefix, c.BucketName); err != nil {
					fmt.Println("上传文件时出错:", err)
				}
			}
		}()
	}

	// 遍历目录并将文件路径发送到通道
	err := filepath.Walk(c.ResourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 处理遍历目录时的错误
		}
		if !info.IsDir() {
			files <- path // 将文件路径发送到通道
		}
		return nil
	})
	handleError(err, "遍历目录时出错:")

	wg.Wait() // 等待所有 goroutine 完成
	fmt.Println("所有文件上传成功!")
}

func main() {
	c := Config{}
	c.Init()
	client := createOSSClient(&c)
	uploadFiles(&c, client)
}
