package client

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"

	"github.com/cortexproject/cortex/pkg/util/flagext"
	"github.com/pkg/errors"
)

type S3Config struct {
	S3ForcePathStyle bool
	Insecure         bool
	Level            aws.LogLevelType

	Bucket   string
	Region   string
	Endpoint string

	AccessKeyID     string
	SecretAccessKey flagext.Secret
}

type S3ObjectClient struct {
	cfg          S3Config
	S3           s3iface.S3API
	S3Uploader   s3manageriface.UploaderAPI
	S3Downloader s3manageriface.DownloaderAPI
}

func NewS3ObjectClient(cfg S3Config) (*S3ObjectClient, error) {
	config, err := generateAwsConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate aws config")
	}

	s3Client, err := buildS3Client(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build s3 client")
	}

	s3Uploader, err := buildS3Uploader(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build s3 uploader")
	}

	s3Downloader, err := buildS3Downloader(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build s3 downloader")
	}

	client := S3ObjectClient{
		cfg:          cfg,
		S3:           s3Client,
		S3Uploader:   s3Uploader,
		S3Downloader: s3Downloader,
	}
	return &client, nil
}

func buildS3Client(config *aws.Config) (*s3.S3, error) {
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new s3 session")
	}

	s3Client := s3.New(sess)
	return s3Client, nil
}

func buildS3Uploader(config *aws.Config) (*s3manager.Uploader, error) {
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new s3 session")
	}

	s3Uploader := s3manager.NewUploader(sess)
	return s3Uploader, nil
}

func buildS3Downloader(config *aws.Config) (*s3manager.Downloader, error) {
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new s3 session")
	}

	s3Downloader := s3manager.NewDownloader(sess)
	return s3Downloader, nil
}

func generateAwsConfig(cfg S3Config) (*aws.Config, error) {
	config := aws.NewConfig()

	config = config.WithMaxRetries(3)
	config = config.WithS3ForcePathStyle(cfg.S3ForcePathStyle)

	if cfg.Insecure {
		config = config.WithDisableSSL(true)
	}

	if cfg.Level != 0 {
		config = config.WithLogLevel(cfg.Level)
	}

	if cfg.Endpoint != "" {
		config = config.WithEndpoint(cfg.Endpoint)
	}

	if cfg.Region != "" {
		config = config.WithRegion(cfg.Region)
	}

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey.String() == "" ||
		cfg.AccessKeyID == "" && cfg.SecretAccessKey.String() != "" {
		return nil, errors.New("must supply both an Access Key ID and Secret Access Key or neither")
	}

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey.String() != "" {
		creds := credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey.String(), "")
		config = config.WithCredentials(creds)
	}

	transport := http.RoundTripper(&http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   200,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	httpClient := &http.Client{
		Transport: transport,
	}

	config = config.WithHTTPClient(httpClient)
	return config, nil
}
