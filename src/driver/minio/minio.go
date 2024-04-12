package minio

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kozmoai/kozmo-supervisor-backend/src/utils/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const MINIO_DEFAULT_SERVE_ADDRESS = "http://127.0.0.1:9000/"
const DEFAULT_PUBLIC_POLICY = `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"],"Sid": ""}]}`
const MINIO_CONNECT_RETRY_TIMES = 6

type MINIOConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
	BucketName      string
	SSLEnabled      bool
	UploadTimeout   time.Duration
}

func NewSystemMINIOConfigByGlobalConfig(config *config.Config) *MINIOConfig {
	return &MINIOConfig{
		AccessKeyID:     config.GetMINIOAccessKeyID(),
		AccessKeySecret: config.GetMINIOAccessKeySecret(),
		Endpoint:        config.GetMINIOEndpoint(),
		BucketName:      config.GetMINIOSystemBucketName(),
		UploadTimeout:   config.GetMINIOTimeout(),
	}
}

func NewTeamMINIOConfigByGlobalConfig(config *config.Config) *MINIOConfig {
	return &MINIOConfig{
		AccessKeyID:     config.GetMINIOAccessKeyID(),
		AccessKeySecret: config.GetMINIOAccessKeySecret(),
		Endpoint:        config.GetMINIOEndpoint(),
		BucketName:      config.GetMINIOTeamBucketName(),
		UploadTimeout:   config.GetMINIOTimeout(),
	}
}

func CreateMINIOInstance(minioConfig *MINIOConfig) *minio.Client {
	var err error

	minioIsntance, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.AccessKeySecret, ""),
		Secure: minioConfig.SSLEnabled,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return minioIsntance
}

type S3Drive struct {
	Instance *minio.Client
	Config   *MINIOConfig
}

func NewS3Drive(minioConfig *MINIOConfig) *S3Drive {
	var errInInit error
	retries := MINIO_CONNECT_RETRY_TIMES
	s3Drive := &S3Drive{
		Config: minioConfig,
	}
	// connect
	s3Drive.Instance = CreateMINIOInstance(minioConfig)
	errInInit = s3Drive.initDefaultBucket()
	for errInInit != nil {
		if retries > 1 {
			log.Printf("Can not connect minio, retrying in %d times\n", retries)
			retries--
			time.Sleep(10 * time.Second)
			errInInit = s3Drive.initDefaultBucket()
			continue
		}
		panic("connect minio failed, please check your minio server config. " + errInInit.Error())
	}
	return s3Drive
}

func (s3Drive *S3Drive) initDefaultBucket() error {
	ctx := context.Background()
	bucketName := s3Drive.Config.BucketName
	err := s3Drive.Instance.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s3Drive.Instance.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own bucket \"%s\"\n", bucketName)
		} else {
			// connection error
			return err
		}
	} else {
		log.Printf("Successfully created bucket \"%s\"\n", bucketName)
	}

	// set policy
	policy := fmt.Sprintf(DEFAULT_PUBLIC_POLICY, bucketName)
	fmt.Printf("policy: %v\n", policy)
	errInSetPloicy := s3Drive.Instance.SetBucketPolicy(context.Background(), bucketName, policy)
	if errInSetPloicy != nil {
		log.Printf("[error in set minio bucket policy] %v\n", errInSetPloicy)
		return errInSetPloicy
	}
	return nil
}

func formatPresignedURLForSelfHostEnv(rawPresignedURL string) string {
	return strings.Replace(rawPresignedURL, MINIO_DEFAULT_SERVE_ADDRESS, "/object-storage/", -1)
}

func (s3Drive *S3Drive) GetPreSignedPutURL(fileName string) (string, error) {
	ctx := context.Background()
	// get put request
	presignedURL, err := s3Drive.Instance.PresignedPutObject(ctx, s3Drive.Config.BucketName, fileName, s3Drive.Config.UploadTimeout)
	if err != nil {
		return "", err
	}
	return formatPresignedURLForSelfHostEnv(presignedURL.String()), nil
}
