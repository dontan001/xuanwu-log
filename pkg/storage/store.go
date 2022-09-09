package storage

import (
	"github.com/kyligence/xuanwu-log/pkg/storage/s3"
)

func Upload(remotePath, srcFile string) error {
	return s3.PutObject(remotePath, srcFile)
}

func Exist(remotePath string) (bool, error) {
	obj, err := s3.HeadObject(remotePath)
	if err != nil {
		return false, err
	}

	if obj == nil {
		return false, nil
	}

	return true, nil
}
