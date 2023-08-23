package ossexample

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/exp/slices"
)

var oss_client *s3.S3
var uploader *s3manager.Uploader
var downloader *s3manager.Downloader

func InitS3() {
	var aws_access_key_id, aws_secret_access_key string = "9ITW6fyeQAy4zfapVlej", "BjdwwduLWAbqTHrMxDMwyoSGe2xsRVrAkLJ3YQ64"
	sess, err := session.NewSessionWithOptions(session.Options{Config: aws.Config{Region: aws.String("us-west-weak-2"), Endpoint: aws.String("http://192.168.214.133:30009"), Credentials: credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, "")}})

	if err != nil {
		log.Fatal("unable to create aws session")
	}
	oss_client = s3.New(sess)

	uploader = s3manager.NewUploader(sess)
	downloader = s3manager.NewDownloader(sess)
}

// 注意bucketname需要以/开头
func ListBucket(bucketname string) []string {
	var values []string
	resp, _ := oss_client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String("/ellis")})
	for _, key := range resp.Contents {
		values = append(values, *key.Key)
	}
	return values
}

func CheckBucketExist(buckname string) bool {
	var contains []string
	lbo, err := oss_client.ListBuckets(nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		for _, v := range lbo.Buckets {
			contains = append(contains, *v.Name)
		}
	}

	result := slices.Contains(contains, buckname)
	return result
}

// 注意bucket需要以/开头
func CreateBucket(bucketname string) {
	oss_client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucketname)})
}

// 上传文件注意bucket以/开头
func UploadFile(filename string, bucketname string, destination string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(destination),
		Body:   f,
	}
	uploader.Upload(upParams)
}

// 文件copy其中sourcebucketandpath是/bucketname+key
func CopyFile(sourcebucketandpath string, destinationbucket string, key string) {
	oss_client.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(destinationbucket), CopySource: aws.String(sourcebucketandpath), Key: aws.String(key)})
}

// 删除文件
func DeleteObject(bucketname string, key string) {
	doo, err := oss_client.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketname), Key: aws.String(key)})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("doo.String(): %v\n", doo.String())
	}
}

// 批量删除文件
func DeleteObjects(bucketname string, keys []string) {
	delete := []*s3.ObjectIdentifier{}

	for i := 0; i < len(keys); i++ {
		delete = append(delete, &s3.ObjectIdentifier{Key: aws.String(keys[i])})
	}
	oss_client.DeleteObjects(&s3.DeleteObjectsInput{Bucket: aws.String(bucketname), Delete: &s3.Delete{Objects: delete}})
}

// 生成下载链接
func GenerateDownloadLink(bucketname string, key string) {
	req, _ := oss_client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("urlStr: %v\n", urlStr)
}

// 下载文件
func DownloadFile(filename string, bucketname string, key string) {
	f, _ := os.Create(filename)
	defer f.Close()
	downloader.Download(f, &s3.GetObjectInput{Bucket: aws.String(bucketname), Key: aws.String(key)})
}

// func main() {

// 	// _, err2 := oss_client.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String("/ellis")})
// 	// if err2 != nil {
// 	// 	fmt.Printf("err2: %v\n", err2)
// 	// }

// 	// resp, _ := oss_client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String("/ellis")})
// 	// for _, key := range resp.Contents {
// 	// 	fmt.Println(*key.Key)
// 	// }
// 	InitS3()
// 	// result := CheckBucketExist("haha")
// 	// if !result {
// 	// 	CreateBucket("/haha")
// 	// }

// 	// UploadFile("README.md", "/haha", "/test/test.md")

// 	// CopyFile("/haha//test/test.md", "/ellis", "/test.md")

// 	// DeleteObject("/ellis", "test.md")

// 	// DeleteObjects("/ellis", []string{"credentials.json", "confluentinc-kafka-connect-elasticsearch-14.0.4.zip"})

// 	// GenerateDownloadLink("/ellis", "ellis.py")

// 	DownloadFile("haha.py", "/ellis", "ellis.py")
// }
