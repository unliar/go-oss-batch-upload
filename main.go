package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// 定义上传文件的并发数量
const defaultConcurrency = 10

func main() {
	// 定义命令行参数
	var accessKeyId string
	var accessKeySecret string
	var endpoint string
	var resourceDir string
	var pathPrefix string
	var bucketName string
	var concurrency int

	flag.StringVar(&accessKeyId, "access-key-id", "", "OSS 账号的 AccessKeyId")
	flag.StringVar(&accessKeySecret, "access-key-secret", "", "OSS 账号的 AccessKeySecret")
	flag.StringVar(&endpoint, "endpoint", "", "OSS 的 Endpoint 地址")
	flag.StringVar(&resourceDir, "resource-dir", "", "资源文件夹路径")
	flag.StringVar(&pathPrefix, "path-prefix", "", "上传到 OSS 的子目录路径")
	flag.StringVar(&bucketName, "bucket-name", "", "上传文件的 Bucket 名称")
	flag.IntVar(&concurrency, "concurrency", defaultConcurrency, "上传文件的并发数量")

	flag.Parse()

	// 验证参数
	if accessKeyId == "" || accessKeySecret == "" || endpoint == "" || resourceDir == "" || bucketName == "" {
		fmt.Println("Error: Missing required parameters")
		flag.Usage()
		return
	}

	// 获取 OSS 客户端
	client, err := oss.New(
		accessKeyId,
		accessKeySecret,
		endpoint,
	)
	if err != nil {
		fmt.Println("Error creating OSS client:", err)
		return
	}

	// 创建一个通道，用于存储待上传的文件
	files := make(chan string, concurrency)

	// 遍历资源文件夹
	err = filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 判断是否是文件
		if !info.IsDir() {
			// 将文件路径发送到通道
			files <- path
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking directory:", err)
		return
	}

	// 创建WaitGroup，用于等待所有文件上传完成
	wg := new(sync.WaitGroup)
	wg.Add(concurrency)

	// 启动并发上传
	for i := 0; i < concurrency; i++ {
		go func() {
			for path := range files {
				// 上传文件
				err := uploadFile(client, bucketName, path, pathPrefix)
				if err != nil {
					fmt.Println("Error uploading file:", err)
				}
			}

			wg.Done()
		}()
	}

	// 等待所有文件上传完成
	wg.Wait()

	fmt.Println("All files uploaded successfully!")
}

func uploadFile(client *oss.Client, bucketName string, path string, pathPrefix string) error {
	// 获取文件名称
	filename := filepath.Base(path)

	// 拼接 OSS 对象键
	objectKey := filepath.Join(pathPrefix, filename)

	// 创建 Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 上传文件
	err = bucket.PutObjectFromFile(objectKey, path)
	if err != nil {
		return err
	}

	return nil
}
