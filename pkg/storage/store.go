package storage

import (
	"github.com/kyligence/xuanwu-log/pkg/storage/s3"
)

func Upload(remotePath, srcFile string) error {
	return s3.PutObject(remotePath, srcFile)
}

func Exist(remotePath string) error {
	return nil
}
