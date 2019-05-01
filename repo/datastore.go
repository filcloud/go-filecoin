package repo

import (
	"errors"
	"fmt"
	"os"

	badgerds "github.com/ipfs/go-ds-badger"
	s3ds "github.com/ipfs/go-ds-s3"
)

const filS3AccessKeyVar = "FIL_S3_ACCESS_KEY"
const filS3SecretKeyVar = "FIL_S3_SECRET_KEY"
const filS3EndpointVar = "FIL_S3_ENDPOINT"
const filS3RegionVar = "FIL_S3_REGION"
const filS3BucketVar = "FIL_S3_BUCKET"

func newDatastore(dsType, path string) (Datastore, error) {
	switch dsType {
	case "badgerds":
		ds, err := badgerds.NewDatastore(path, badgerOptions())
		if err != nil {
			return nil, err
		}
		return ds, nil
	case "s3ds":
		cfg, err := s3dsConfig()
		if err != nil {
			return nil, err
		}
		ds, err := s3ds.NewS3Datastore(*cfg)
		if err != nil {
			return nil, err
		}
		return ds, nil
	default:
		return nil, fmt.Errorf("unknown datastore type: %s", dsType)
	}
}

func s3dsConfig() (*s3ds.Config, error) {
	accessKey := os.Getenv(filS3AccessKeyVar)
	secretKey := os.Getenv(filS3SecretKeyVar)
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("access key and secret key must be specified")
	}
	endpoint := os.Getenv(filS3EndpointVar)
	if endpoint == "" {
		endpoint = "http://127.0.0.1:9000"
	}
	region := os.Getenv(filS3RegionVar)
	if region == "" {
		region = "us-east-1"
	}
	bucket := os.Getenv(filS3BucketVar)
	if bucket == "" {
		bucket = "fil"
	}
	return &s3ds.Config{
		AccessKey:      accessKey,
		SecretKey:      secretKey,
		RegionEndpoint: endpoint,
		Region:         region,
		Bucket:         bucket,
	}, nil
}
