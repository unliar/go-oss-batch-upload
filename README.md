## OSS 批量上传工具

### 中文

**简介**

OSS 批量上传工具是一个 Golang 程序，用于并发上传文件到 OSS。该程序可以从命令行或环境变量中读取配置，并支持上传到 OSS 的子目录下。

**安装**

```
go install github.com/unliar/go-oss-batch-upload@0.0.3
```

**使用**

```dockerfile
FROM registry.cn-shenzhen.aliyuncs.com/happysooner/go-oss-batch-upload:$TAG

CMD /app/main -access-key-id your-access-key-id -access-key-secret your-access-key-secret -endpoint your-endpoint -resource-dir /path/to/resource/dir -path-prefix subdir -bucket-name my-bucket
```

```
go-oss-batch-upload -access-key-id your-access-key-id -access-key-secret your-access-key-secret -endpoint your-endpoint -resource-dir /path/to/resource/dir -path-prefix subdir -bucket-name my-bucket
```

**参数**

* `-access-key-id`: OSS 账号的 AccessKeyId
* `-access-key-secret`: OSS 账号的 AccessKeySecret
* `-endpoint`: OSS 的 Endpoint 地址
* `-resource-dir`: 资源文件夹路径
* `-path-prefix`: 上传到 OSS 的子目录路径
* `-bucket-name`: 上传文件的 Bucket 名称
* `-concurrency`: 上传文件的并发数量 (可选，默认值: 10)

**示例**

```
go-oss-batch-upload -access-key-id AKID... -access-key-secret YOUR_ACCESS_KEY_SECRET -endpoint oss-cn-hangzhou.aliyuncs.com -resource-dir /path/to/resource/dir -path-prefix subdir -bucket-name my-bucket
```

**注意**

* 请确保您已安装 Golang 并将 `go` 命令添加到系统路径中。
* 请替换命令中的参数值以匹配您的实际情况。
* -path-prefix 不可以有前缀 / 

### English

**Introduction**

OSS Batch Uploader is a Golang program that uploads files to OSS concurrently. It can read configuration from the command line or environment variables, and supports uploading to subdirectories under OSS.

**Installation**

```
go install github.com/unliar/go-oss-batch-upload@0.0.3
```

**Usage**

```
oss-batch-uploader -access-key-id your-access-key-id -access-key-secret your-access-key-secret -endpoint your-endpoint -resource-dir /path/to/resource/dir -path-prefix subdir -bucket-name my-bucket
```

**Parameters**

* `-access-key-id`: The AccessKeyId of your OSS account
* `-access-key-secret`: The AccessKeySecret of your OSS account
* `-endpoint`: The endpoint of OSS
* `-resource-dir`: The path of the resource directory
* `-path-prefix`: The subdirectory path to upload to OSS
* `-bucket-name`: The name of the bucket to upload files to
* `-concurrency`: The number of concurrent uploads (optional, default: 10)

**Example**

```dockerfile
FROM registry.cn-shenzhen.aliyuncs.com/happysooner/go-oss-batch-upload:$TAG

CMD /app/main -access-key-id your-access-key-id -access-key-secret your-access-key-secret -endpoint your-endpoint -resource-dir /path/to/resource/dir -path-prefix subdir -bucket-name my-bucket
```

```
go-oss-batch-upload -access-key-id AKID... -access-key-secret YOUR_ACCESS_KEY_SECRET -endpoint oss-cn-hangzhou.aliyuncs.com -resource-dir /path/to/resource/dir -path-prefix subdir -bucket-name my-bucket
```

**Notes**

* Make sure you have installed Golang and added the `go` command to the system path.
* Please replace the parameter values in the command to match your actual situation.

**Additional Information**

* Golang Documentation: [https://golang.org/](https://golang.org/)