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

const (
	bucket = "donggetest"
	region = endpoints.UsWest2RegionID

	logLevel            = aws.LogDebug
	uploadConcurrency   = s3manager.DefaultUploadConcurrency
	downloadConcurrency = s3manager.DefaultDownloadConcurrency
)

/*var (
	cfg = client.S3Config{
		Level:           logLevel,
		Insecure:        true,
		Region:          region,
		AccessKeyID:     AccessKeyID,
		SecretAccessKey: flagext.Secret{Value: SecretAccessKey},
	}

	objectClient, _ = client.NewS3ObjectClient(cfg)
)*/

type S3Config struct {
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}

type S3Store struct {
	Config *S3Config
	client *client.S3ObjectClient
}

func (s *S3Store) Setup() {
	if s.Config != nil {
		cfg := client.S3Config{
			Level:           logLevel,
			Insecure:        true,
			Bucket:          s.Config.Bucket,
			Region:          s.Config.Region,
			AccessKeyID:     AccessKeyID,
			SecretAccessKey: flagext.Secret{Value: SecretAccessKey},
		}

		objectClient, _ := client.NewS3ObjectClient(cfg)
		s.client = objectClient
	}
}

func (s *S3Store) GetBuckets() error {
	ctx := context.Background()

	result, err := s.client.S3.ListBucketsWithContext(ctx, nil)
	if err != nil {
		log.Printf("Unable to list buckets, %v", err)
	}

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	return nil
}

func (s *S3Store) GetObjects() error {
	ctx := context.Background()

	objects := []string{}
	err := s.client.S3.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(s.Config.Bucket),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		for _, o := range p.Contents {
			objects = append(objects, aws.StringValue(o.Key))
		}
		return true // continue paging
	})
	if err != nil {
		return fmt.Errorf("failed to list objects for bucket, %s, %v", s.Config.Bucket, err)
	}

	for idx, obj := range objects {
		fmt.Printf("Objects in bucket: %d, %s \n", idx, obj)
	}
	return nil
}

func (s *S3Store) HeadObject(remotePath string) (*s3.HeadObjectOutput, error) {
	result, err := s.client.S3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(remotePath),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == "NotFound" {
			return nil, nil
		}

		return nil, err
	}

	log.Printf("Object head: %v", result)
	return result, nil
}

func (s *S3Store) GetObject(remotePath, destFile string) error {
	file, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("unable to open file %s, %v", destFile, err)
	}

	defer file.Close()

	objectMeta, err := s.client.S3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(remotePath),
	})
	if err != nil {
		return fmt.Errorf("unable to get object info %q, %v", remotePath, err)
	}
	size := objectMeta.ContentLength

	defer util.TimeMeasureRate("GetObject", *size)()
	numBytes, err := s.client.S3Downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.Config.Bucket),
			Key:    aws.String(remotePath),
		}, func(downloader *s3manager.Downloader) {
			downloader.PartSize = s3manager.DefaultDownloadPartSize * 10
			/*downloader.BufferProvider = s3manager.NewPooledBufferedWriterReadFromProvider(
			int(s3manager.DefaultDownloadPartSize * 5))*/

			downloader.Concurrency = downloadConcurrency
		})
	if err != nil {
		return fmt.Errorf("unable to download file %q, %v", remotePath, err)
	}

	log.Printf("Successfully downloaded %q from bucket %q with %d bytes \n", file.Name(), s.Config.Bucket, numBytes)
	return nil
}

func (s *S3Store) PutObject(remotePath, srcFile string) error {
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
	_, err = s.client.S3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(remotePath),
		Body:   file,
	}, func(uploader *s3manager.Uploader) {
		// set the part size as low as possible to avoid timeouts and aborts
		var partSize int64 = s3manager.DefaultUploadPartSize
		maxParts := math.Ceil(float64(info.Size() / partSize))

		// 10000 parts is the limit for AWS S3. If the resulting number of parts would exceed that limit, increase the
		// part size as much as needed but as little possible
		if maxParts > s3manager.MaxUploadParts {
			partSize = int64(math.Ceil(float64(info.Size()) / s3manager.MaxUploadParts))
		}

		uploader.Concurrency = uploadConcurrency
		uploader.LeavePartsOnError = false
		uploader.PartSize = partSize
	})
	if err != nil {
		return fmt.Errorf("unable to upload %s to %s, %v", srcFile, s.Config.Bucket, err)
	}

	log.Printf("Successfully uploaded %q to %q\n", srcFile, s.Config.Bucket)
	return nil
}

func (s *S3Store) DelObject(remotePath string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(remotePath),
	}

	_, err := s.client.S3.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "NotFound" {
				return nil
			}
		}
		return err
	}

	return nil
}
