package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

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
	flag.IntVar(&c.Concurrency, "concurrency", 10, "上传文件的并发数量")
	flag.Parse()
}

func uploadFile(client *oss.Client, path string, resourceDir string, pathPrefix string, bucketName string) error {
	filename := filepath.Base(path)
	objectKey := filepath.Join(pathPrefix, filepath.Dir(strings.TrimPrefix(path, resourceDir)), filename)
	fmt.Println("File upload Ready", objectKey)
	bk, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error creating OSS bucket:", err)
		return err
	}
	err = bk.PutObjectFromFile(objectKey, path)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return err
	} else {
		fmt.Println("File uploaded successfully:", objectKey)
		return nil
	}
}

func createOSSClient(c *Config) *oss.Client {
	client, err := oss.New(
		c.Endpoint,
		c.AccessKeyID,
		c.AccessKeySecret,
	)
	if err != nil {
		fmt.Println("Error creating OSS client:", err)
		os.Exit(1)
	}
	return client
}

func uploadFiles(c *Config, client *oss.Client) {
	files := make(chan string, c.Concurrency)
	var wg sync.WaitGroup
	wg.Add(c.Concurrency)
	for i := 0; i < c.Concurrency; i++ {
		go func() {
			for path := range files {
				err := uploadFile(client, path, c.ResourceDir, c.PathPrefix, c.BucketName)
				if err != nil {
					fmt.Println("Error uploading file:", err)
					os.Exit(1)
				}
			}
			wg.Done()
		}()
	}
	err := filepath.Walk(c.ResourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
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
	wg.Wait()
	fmt.Println("Uploaded all files successfully!")
}

func main() {
	c := Config{}
	c.Init()

	client := createOSSClient(&c)
	uploadFiles(&c, client)
}
