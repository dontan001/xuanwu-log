package s3

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/cortexproject/cortex/pkg/util/flagext"
	"github.com/kyligence/xuanwu-log/pkg/storage/s3/client"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

var (
	bucket = "donggetest"

	cfg = client.S3Config{
		Level:           aws.LogDebugWithHTTPBody,
		Insecure:        true,
		Region:          endpoints.UsWest2RegionID,
		AccessKeyID:     AccessKeyID,
		SecretAccessKey: flagext.Secret{Value: SecretAccessKey},
	}

	objectClient, _ = client.NewS3ObjectClient(cfg)
)

func GetBuckets() error {
	ctx := context.Background()

	result, err := objectClient.S3.ListBucketsWithContext(ctx, nil)
	if err != nil {
		log.Printf("Unable to list buckets, %v", err)
	}

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	return nil
}

func GetObjects() error {
	ctx := context.Background()

	objects := []string{}
	err := objectClient.S3.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		for _, o := range p.Contents {
			objects = append(objects, aws.StringValue(o.Key))
		}
		return true // continue paging
	})
	if err != nil {
		return fmt.Errorf("failed to list objects for bucket, %s, %v", bucket, err)
	}

	for idx, obj := range objects {
		fmt.Printf("Objects in bucket: %d, %s \n", idx, obj)
	}
	return nil
}

func HeadObject(remotePath string) error {
	result, err := objectClient.S3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remotePath),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			// Specific error code handling
		}
		return err
	}

	fmt.Println(result)
	return nil
}

func GetObject(remotePath, destFile string) error {
	file, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("unable to open file %s, %v", destFile, err)
	}

	defer file.Close()
	numBytes, err := objectClient.S3Downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(remotePath),
		})
	if err != nil {
		return fmt.Errorf("unable to download file %q, %v", remotePath, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return nil
}

func PutObject(remotePath, srcFile string) error {
	file, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("unable to open file %v", err)
	}

	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get size of file %s: %v", srcFile, err)
	}

	defer util.TimeMeasureRate("PutObject", info.Size())()
	_, err = objectClient.S3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remotePath),
		Body:   file,
	}, func(uploader *s3manager.Uploader) {
		// set the part size as low as possible to avoid timeouts and aborts
		var partSize int64 = s3manager.MinUploadPartSize
		maxParts := math.Ceil(float64(info.Size() / partSize))

		// 10000 parts is the limit for AWS S3. If the resulting number of parts would exceed that limit, increase the
		// part size as much as needed but as little possible
		if maxParts > s3manager.MaxUploadParts {
			partSize = int64(math.Ceil(float64(info.Size()) / s3manager.MaxUploadParts))
		}

		uploader.Concurrency = s3manager.DefaultUploadConcurrency
		uploader.LeavePartsOnError = false
		uploader.PartSize = partSize
	})
	if err != nil {
		return fmt.Errorf("unable to upload %s to %s, %v", srcFile, bucket, err)
	}

	log.Printf("Successfully uploaded %q to %q\n", srcFile, bucket)
	return nil
}

func DelObject(remotePath string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remotePath),
	}

	_, err := objectClient.S3.DeleteObject(input)
	return err
}
