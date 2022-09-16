package storage

import (
	"log"

	"github.com/kyligence/xuanwu-log/pkg/storage/s3"
)

type Store struct {
	Config *s3.S3Config
	S3     *s3.S3Store
}

func (store *Store) Setup() {
	if store.Config != nil {
		store.S3 = &s3.S3Store{
			Config: &s3.S3Config{
				Bucket: store.Config.Bucket,
				Region: store.Config.Region,
			},
		}

		store.S3.Setup()
		log.Printf("Store setup finish.")
	}
}

func (store *Store) Upload(remotePath, srcFile string) error {
	return store.S3.PutObject(remotePath, srcFile)
}

func (store *Store) Exist(remotePath string) (bool, error) {
	obj, err := store.S3.HeadObject(remotePath)
	if err != nil {
		return false, err
	}

	if obj == nil {
		return false, nil
	}

	return true, nil
}
