package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		endpoint,
		accessKeyId,
		accessKeySecret,
	)
	if err != nil {
		fmt.Println("Error creating OSS client:", err)
		return
	}

	// 创建一个通道，用于存储待上传的文件
	files := make(chan string, concurrency)

	// 并发上传
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			for path := range files {
				// 获取文件名称
				filename := filepath.Base(path)

				// 拼接 OSS 对象键
				objectKey := filepath.Join(pathPrefix, filepath.Dir(strings.TrimPrefix(path, resourceDir)), filename)
				fmt.Println("File upload Ready", objectKey)
				// 上传文件
				bk, err := client.Bucket(bucketName)
				if err != nil {
					fmt.Println("Error uploading Bucket:", err)
				} else {
					err = bk.PutObjectFromFile(objectKey, path)
					if err != nil {
						fmt.Println("Error uploading file:", err)
					} else {
						fmt.Println("File uploaded successfully:", objectKey)
					}

				}
			}

			wg.Done()
		}()
	}

	// 遍历资源文件夹，将文件路径发送到通道
	err = filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 判断是否是文件
		if !info.IsDir() {
			files <- path
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking directory:", err)
		return
	}

	close(files)

	// 等待所有文件上传完成
	wg.Wait()
	fmt.Println("All files uploaded successfully!")
}
