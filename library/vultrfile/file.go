package vultrfile

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"case-study-kredit-plus/configs"
)

func UploadFile(file *bytes.Reader, dst string) error {
	atomConfig, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalln("failed to get configuration: ", errConfig)
	}

	region := "us-east-1"

	credentials := credentials.NewStaticCredentialsProvider(
		atomConfig.VultrAccessKey,
		atomConfig.VultrSecretKey,
		"",
	)

	endpoint := aws.Endpoint{
		PartitionID:   "aws",
		URL:           atomConfig.VultrHostname,
		SigningRegion: region,
	}

	endpointResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) { return endpoint, nil })

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials),
		config.WithRegion(region),
		config.WithEndpointResolver(endpointResolver),
	)

	if err != nil {
		return err
	}

	uploadObj := &s3.PutObjectInput{
		Bucket: aws.String(atomConfig.VultrBucket),
		Key:    aws.String(dst),
		Body:   file,
		ACL:    "public-read",
	}

	config := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })
	_, err = config.PutObject(context.TODO(), uploadObj)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeleteFile(dst string) error {
	atomConfig, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalln("failed to get configuration: ", errConfig)
	}

	credentials := credentials.NewStaticCredentialsProvider(
		atomConfig.VultrAccessKey,
		atomConfig.VultrSecretKey,
		"",
	)

	endpoint := aws.Endpoint{
		PartitionID:   "aws",
		URL:           atomConfig.VultrHostname,
		SigningRegion: atomConfig.VultrRegion,
	}

	endpointResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) { return endpoint, nil })

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials),
		config.WithRegion(atomConfig.VultrRegion),
		config.WithEndpointResolver(endpointResolver),
	)

	if err != nil {
		return err
	}

	uploadObj := &s3.DeleteObjectInput{
		Bucket: aws.String(atomConfig.VultrBucket),
		Key:    aws.String(dst),
	}

	config := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })
	_, err = config.DeleteObject(context.TODO(), uploadObj)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetMainURL() string {
	atomConfig, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalln("failed to get configuration: ", errConfig)
	}

	url := fmt.Sprintf("%s/%s", atomConfig.VultrHostname, atomConfig.VultrBucket)

	return url
}
